# 架构设计文档

本文档描述 BookImporter 项目的架构设计和实现细节。

## 目录

- [项目概述](#项目概述)
- [整体架构](#整体架构)
- [核心模块](#核心模块)
- [命令流程](#命令流程)
- [数据流](#数据流)
- [设计决策](#设计决策)
- [扩展性](#扩展性)

## 项目概述

BookImporter 是一个命令行工具，用于管理和整理电子书库。采用 Go 语言开发，使用 Cobra 框架构建 CLI。

### 设计目标

1. **简单易用**: 提供直观的命令行界面
2. **性能优良**: 高效处理大量文件
3. **可靠稳定**: 完善的错误处理
4. **易于扩展**: 模块化设计，便于添加新功能

## 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                      用户界面 (CLI)                      │
│                     main.go → cmd/                       │
└─────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        ↓                   ↓                   ↓
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│ clname 命令  │   │ rename 命令  │   │ version 命令 │
│ cmd/clname.go│   │ cmd/rename.go│   │cmd/version.go│
└──────────────┘   └──────────────┘   └──────────────┘
        │                   │
        └───────────────────┼───────────────────┘
                            ↓
        ┌──────────────────────────────────────┐
        │          工具库 (pkg/util/)          │
        │                                      │
        │  ┌────────────┐   ┌──────────────┐  │
        │  │ cleanname  │   │  filetool    │  │
        │  └────────────┘   └──────────────┘  │
        └──────────────────────────────────────┘
                            ↓
        ┌──────────────────────────────────────┐
        │          外部依赖                     │
        │  • Cobra (CLI 框架)                  │
        │  • epub (EPUB 解析)                  │
        │  • ebook-meta (Calibre)              │
        └──────────────────────────────────────┘
```

### 分层架构

1. **表现层 (Presentation)**: CLI 命令和参数处理
2. **业务逻辑层 (Business Logic)**: 核心功能实现
3. **工具层 (Utility)**: 通用工具函数
4. **数据层 (Data)**: 文件系统操作

## 核心模块

### 1. 命令模块 (cmd/)

#### cmd/root.go

根命令，作为所有子命令的容器。

```go
var rootCmd = &cobra.Command{
    Use:   "bookimporter",
    Short: "Import books into your library",
}

func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}

func init() {
    // 注册子命令
    rootCmd.AddCommand(clnameCmd)
    rootCmd.AddCommand(renameCmd)
    rootCmd.AddCommand(versionCmd)
}
```

**职责:**
- 初始化 Cobra 应用
- 注册子命令
- 提供统一的错误处理

#### cmd/clname.go

清理书籍标题命令。

**核心结构:**

```go
type ClnameConfig struct {
    Path  string  // 文件或目录路径
    DoTry bool    // 尝试模式
    Debug bool    // 调试模式
    Skip  bool    // 跳过错误
}
```

**主要流程:**

1. 参数验证 (`ValidateConfig`)
2. 文件发现（单个文件或目录）
3. EPUB 解析 (`ParseEpub`)
4. 标题清理 (`util.TryCleanTitle`)
5. 元数据更新（调用 `ebook-meta`）

**关键函数:**

- `ValidateConfig()`: 验证配置参数
- `ParseEpub()`: 解析 EPUB 文件并处理
- `exec.Command()`: 调用外部 ebook-meta 工具

#### cmd/rename.go

批量重命名命令。

**核心结构:**

```go
type RenameConfig struct {
    Debug      bool
    DoTry      bool
    Formats    []string  // 文件格式过滤
    Recursive  bool      // 递归搜索
    SourceDir  string    // 源目录
    OutputDir  string    // 输出目录
    Template   string    // 文件名模板
    StartIndex int       // 起始序号
}
```

**主要流程:**

1. 参数解析和验证
2. 文件发现 (`findFiles`)
3. 文件名构建 (`buildNewName`)
4. 文件重命名/移动 (`os.Rename`)

**关键函数:**

- `validateConfig()`: 验证模板格式
- `findFiles()`: 递归查找匹配的文件
- `buildNewName()`: 根据模板生成新文件名

### 2. 工具模块 (pkg/util/)

#### pkg/util/cleanname.go

标题清理逻辑。

**核心功能:**

```go
func TryCleanTitle(title string) string {
    // 移除中文括号
    title = removePattern(title, `【.*?】`)
    title = removePattern(title, `（.*?）`)
    
    // 移除英文括号
    title = removePattern(title, `\[.*?\]`)
    title = removePattern(title, `\(.*?\)`)
    
    return strings.TrimSpace(title)
}
```

**设计考虑:**

- 使用正则表达式进行模式匹配
- 支持嵌套括号
- 保留其他格式的文本

#### pkg/util/filetool.go

文件操作工具。

**主要函数:**

```go
func Exists(path string) bool
func IsFile(path string) bool
func IsDir(path string) bool
```

提供统一的文件系统操作接口。

## 命令流程

### clname 命令流程

```
用户输入命令
    ↓
参数解析 (Cobra)
    ↓
验证配置 (ValidateConfig)
    ↓
判断路径类型
    ├─→ 单个文件 → ParseEpub
    └─→ 目录 → filepath.Glob → 遍历 → ParseEpub
              ↓
        解析 EPUB (epub.Open)
              ↓
        获取标题 (book.Opf.Metadata.Title)
              ↓
        清理标题 (util.TryCleanTitle)
              ↓
        比较新旧标题
              ├─→ 相同 → 跳过
              └─→ 不同 → 调用 ebook-meta 更新
                        ↓
                  显示结果
```

### rename 命令流程

```
用户输入命令
    ↓
参数解析 (Cobra)
    ↓
验证配置 (validateConfig)
    ↓
查找文件 (findFiles)
    ├─→ 递归模式 → 遍历子目录
    └─→ 非递归 → 当前目录
              ↓
        格式过滤
              ↓
文件列表 → 遍历
              ↓
        生成新文件名 (buildNewName)
              ↓
        构建目标路径
              ├─→ 有输出目录 → 移动到输出目录
              └─→ 无输出目录 → 原地重命名
                        ↓
                  执行重命名 (os.Rename)
                        ↓
                  记录结果
                        ↓
                  显示统计
```

## 数据流

### clname 数据流

```
EPUB 文件
    ↓
[读取] epub.Open()
    ↓
Metadata 对象
    ↓
[提取] Title 字段
    ↓
原始标题字符串
    ↓
[清理] TryCleanTitle()
    ↓
清理后标题字符串
    ↓
[比较] 新旧标题
    ↓
[更新] ebook-meta 命令
    ↓
更新的 EPUB 文件
```

### rename 数据流

```
源目录路径
    ↓
[扫描] filepath.Glob()
    ↓
文件路径列表
    ↓
[过滤] 格式匹配
    ↓
待处理文件列表
    ↓
[转换] buildNewName()
    ↓
(原路径, 新路径) 对
    ↓
[执行] os.Rename()
    ↓
重命名后的文件
```

## 设计决策

### 1. 为什么使用 Cobra 框架？

**优势:**
- 自动生成帮助信息
- 强大的参数解析
- 子命令管理
- 广泛使用，社区活跃

**替代方案:**
- flag 标准库：功能较弱
- urfave/cli：功能类似，但 Cobra 更流行

### 2. 为什么调用外部 ebook-meta？

**原因:**
- Calibre 是成熟的电子书管理工具
- ebook-meta 功能完善，支持多种格式
- 避免重新实现复杂的 EPUB 元数据编写逻辑

**缺点:**
- 需要额外安装 Calibre
- 性能略低（进程调用开销）

**未来改进:**
可以考虑直接操作 EPUB 文件，减少外部依赖。

### 3. 为什么使用 kapmahc/epub 库？

**原因:**
- 轻量级，只用于读取元数据
- API 简单
- 满足当前需求

**替代方案:**
- taylorskalyo/goreader/epub：功能更强大
- 自己实现：过于复杂

### 4. 错误处理策略

**原则:**
- 用户操作错误：显示友好提示，退出
- 可恢复错误：记录并继续（使用 `-j` 选项）
- 严重错误：panic（开发阶段）→ 改为返回错误（生产）

**示例:**

```go
// 不可恢复错误
if !util.Exists(c.Path) {
    fmt.Println("文件路径不存在，请检查")
    os.Exit(1)
}

// 可恢复错误
if err != nil && c.Skip {
    fmt.Printf("跳过: %v\n", err)
    continue
}
```

## 扩展性

### 添加新命令

1. 在 `cmd/` 创建新文件，如 `cmd/export.go`
2. 定义 Cobra 命令：

```go
var exportCmd = &cobra.Command{
    Use:   "export",
    Short: "Export book list",
    Run: func(cmd *cobra.Command, args []string) {
        // 实现逻辑
    },
}

func init() {
    // 添加参数
    exportCmd.Flags().StringP("output", "o", "books.csv", "Output file")
}
```

3. 在 `cmd/root.go` 注册：

```go
func init() {
    rootCmd.AddCommand(exportCmd)
}
```

### 添加新的清理规则

修改 `pkg/util/cleanname.go`：

```go
func TryCleanTitle(title string) string {
    // 现有规则
    title = removePattern(title, `【.*?】`)
    
    // 添加新规则
    title = removePattern(title, `《.*?》`)
    
    return strings.TrimSpace(title)
}
```

### 支持新的文件格式

以 PDF 为例：

1. 添加 PDF 解析库

```go
import "github.com/pdfcpu/pdfcpu/pkg/api"
```

2. 实现 PDF 处理函数

```go
func ParsePDF(file string, config *Config) error {
    // 读取 PDF 元数据
    // 清理标题
    // 写回元数据
}
```

3. 在 clname 命令中添加支持

```go
if strings.HasSuffix(c.Path, ".pdf") {
    ParsePDF(c.Path, c)
} else if strings.HasSuffix(c.Path, ".epub") {
    ParseEpub(c.Path, c)
}
```

### 插件化设计（未来）

考虑引入插件系统：

```
plugins/
├── epub/
│   └── cleaner.go
├── pdf/
│   └── cleaner.go
└── interface.go

// interface.go
type Cleaner interface {
    SupportedFormats() []string
    Clean(file string) error
}
```

## 性能考虑

### 当前性能

- **clname**: 主要瓶颈是 ebook-meta 调用（每个文件 ~100-500ms）
- **rename**: 文件系统操作，性能较好（每个文件 ~1-10ms）

### 优化方向

1. **并发处理**

```go
// 使用 goroutine 池
pool := workerpool.New(runtime.NumCPU())
for _, file := range files {
    file := file
    pool.Submit(func() {
        ProcessFile(file)
    })
}
pool.StopWait()
```

2. **批量操作**

对于 rename 命令，可以批量调用 `os.Rename`。

3. **缓存**

缓存已处理文件的结果，避免重复处理。

4. **直接操作 EPUB**

实现直接修改 EPUB 元数据，避免调用外部工具。

## 测试策略

### 单元测试

- 每个公共函数都有对应测试
- 使用表驱动测试
- 模拟文件系统操作

### 集成测试

- 测试完整的命令执行流程
- 使用临时目录
- 验证最终结果

### 测试覆盖率

当前目标：> 80%

查看覆盖率：

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 安全考虑

1. **路径遍历**: 验证所有路径参数
2. **命令注入**: 转义传递给 ebook-meta 的参数
3. **文件权限**: 检查文件访问权限
4. **资源限制**: 限制并发数量，防止资源耗尽

## 未来规划

### 短期（v1.0）

- [ ] 完善错误处理
- [ ] 添加更多单元测试
- [ ] 改进日志输出
- [ ] 添加配置文件支持

### 中期（v2.0）

- [ ] 支持更多电子书格式（PDF, MOBI）
- [ ] 并发处理
- [ ] 进度条显示
- [ ] 导出书籍列表功能

### 长期（v3.0）

- [ ] GUI 界面
- [ ] 插件系统
- [ ] 云端同步
- [ ] 书籍元数据自动补全

## 贡献

欢迎贡献！请参阅 [CONTRIBUTING.md](CONTRIBUTING.md)

---

最后更新: 2025-11-28

