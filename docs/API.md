# API 文档

本文档描述 BookImporter 的代码 API，供开发者集成和扩展使用。

## 目录

- [包概览](#包概览)
- [cmd 包](#cmd-包)
- [pkg/util 包](#pkgutil-包)
- [使用示例](#使用示例)

## 包概览

```
github.com/jianyun8023/bookimporter/
├── cmd/                    # 命令行命令
│   ├── clname.go          # 清理标题命令
│   ├── rename.go          # 重命名命令
│   ├── root.go            # 根命令
│   └── version.go         # 版本命令
└── pkg/
    └── util/              # 工具函数
        ├── cleanname.go   # 标题清理
        └── filetool.go    # 文件操作
```

## cmd 包

### cmd/root.go

根命令和入口点。

#### Execute()

```go
func Execute()
```

执行根命令，启动 CLI 应用。这是程序的主入口点。

**用法:**

```go
import "github.com/jianyun8023/bookimporter/cmd"

func main() {
    cmd.Execute()
}
```

## pkg/util 包

### pkg/util/cleanname.go

标题清理相关函数。

#### TryCleanTitle()

```go
func TryCleanTitle(title string) string
```

清理书籍标题，移除括号和其中的内容。

**参数:**
- `title`: 原始标题字符串

**返回:**
- 清理后的标题字符串

**清理规则:**
- 移除 `【...】` 中文方括号及内容
- 移除 `（...）` 中文圆括号及内容
- 移除 `[...]` 英文方括号及内容
- 移除 `(...)` 英文圆括号及内容
- 去除首尾空白

**示例:**

```go
import "github.com/jianyun8023/bookimporter/pkg/util"

title := "三体【刘慈欣】"
cleaned := util.TryCleanTitle(title)
// cleaned = "三体"

title2 := "Python编程(第3版)[入门到实践]"
cleaned2 := util.TryCleanTitle(title2)
// cleaned2 = "Python编程"
```

**注意事项:**
- 支持嵌套括号
- 如果标题完全由括号包围，可能返回空字符串
- 不会修改其他类型的标点符号

### pkg/util/filetool.go

文件系统操作工具函数。

#### Exists()

```go
func Exists(path string) bool
```

检查文件或目录是否存在。

**参数:**
- `path`: 文件或目录路径

**返回:**
- `true`: 文件或目录存在
- `false`: 不存在

**示例:**

```go
if util.Exists("/path/to/file.epub") {
    // 文件存在，继续处理
}
```

#### IsFile()

```go
func IsFile(path string) bool
```

检查路径是否为文件。

**参数:**
- `path`: 文件路径

**返回:**
- `true`: 是文件
- `false`: 不是文件或不存在

**示例:**

```go
if util.IsFile("/path/to/file.epub") {
    // 是文件
}
```

#### IsDir()

```go
func IsDir(path string) bool
```

检查路径是否为目录。

**参数:**
- `path`: 目录路径

**返回:**
- `true`: 是目录
- `false`: 不是目录或不存在

**示例:**

```go
if util.IsDir("/path/to/books") {
    // 是目录
}
```

## 使用示例

### 作为库使用

虽然 BookImporter 主要是命令行工具，但也可以作为库在其他 Go 项目中使用。

#### 清理标题

```go
package main

import (
    "fmt"
    "github.com/jianyun8023/bookimporter/pkg/util"
)

func main() {
    titles := []string{
        "三体【刘慈欣】",
        "Python编程（第3版）",
        "算法导论[第三版]",
    }
    
    for _, title := range titles {
        cleaned := util.TryCleanTitle(title)
        fmt.Printf("原标题: %s\n", title)
        fmt.Printf("清理后: %s\n\n", cleaned)
    }
}
```

#### 文件操作

```go
package main

import (
    "fmt"
    "github.com/jianyun8023/bookimporter/pkg/util"
)

func main() {
    path := "/path/to/file.epub"
    
    if !util.Exists(path) {
        fmt.Println("文件不存在")
        return
    }
    
    if util.IsFile(path) {
        fmt.Println("是文件")
    } else if util.IsDir(path) {
        fmt.Println("是目录")
    }
}
```

### 扩展命令

添加自定义命令：

```go
package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/jianyun8023/bookimporter/pkg/util"
)

var myCmd = &cobra.Command{
    Use:   "mycmd",
    Short: "自定义命令",
    Run: func(cmd *cobra.Command, args []string) {
        // 使用 util 包的功能
        path := args[0]
        if util.Exists(path) {
            fmt.Println("处理文件:", path)
            // 自定义逻辑
        }
    },
}

func init() {
    rootCmd.AddCommand(myCmd)
}
```

### 批量处理

```go
package main

import (
    "fmt"
    "path/filepath"
    "github.com/jianyun8023/bookimporter/pkg/util"
    "github.com/kapmahc/epub"
)

func processBooks(dir string) error {
    files, err := filepath.Glob(filepath.Join(dir, "*.epub"))
    if err != nil {
        return err
    }
    
    for _, file := range files {
        if !util.IsFile(file) {
            continue
        }
        
        // 打开 EPUB
        book, err := epub.Open(file)
        if err != nil {
            fmt.Printf("无法打开 %s: %v\n", file, err)
            continue
        }
        
        // 获取标题
        if len(book.Opf.Metadata.Title) == 0 {
            continue
        }
        
        oldTitle := book.Opf.Metadata.Title[0]
        newTitle := util.TryCleanTitle(oldTitle)
        
        if oldTitle != newTitle {
            fmt.Printf("文件: %s\n", file)
            fmt.Printf("  旧标题: %s\n", oldTitle)
            fmt.Printf("  新标题: %s\n", newTitle)
            // 这里可以调用 ebook-meta 更新
        }
    }
    
    return nil
}

func main() {
    err := processBooks("/path/to/books")
    if err != nil {
        fmt.Println("错误:", err)
    }
}
```

### 自定义清理规则

```go
package main

import (
    "fmt"
    "regexp"
    "strings"
)

// CustomCleanTitle 自定义清理函数
func CustomCleanTitle(title string) string {
    // 使用内置函数
    // title = util.TryCleanTitle(title)
    
    // 添加自定义规则
    // 移除书名号
    re := regexp.MustCompile(`《.*?》`)
    title = re.ReplaceAllString(title, "")
    
    // 移除版本号
    re = regexp.MustCompile(`v\d+\.\d+`)
    title = re.ReplaceAllString(title, "")
    
    // 移除多余空格
    title = strings.Join(strings.Fields(title), " ")
    
    return strings.TrimSpace(title)
}

func main() {
    title := "《三体》【刘慈欣】v1.0"
    cleaned := CustomCleanTitle(title)
    fmt.Println(cleaned) // 输出: 三体
}
```

## 错误处理

### 错误类型

BookImporter 使用标准 Go 错误处理机制。主要错误类型：

1. **文件不存在**: `os.ErrNotExist`
2. **权限错误**: `os.ErrPermission`
3. **EPUB 解析错误**: epub 库返回的错误
4. **命令执行错误**: exec 包返回的错误

### 错误处理示例

```go
import (
    "errors"
    "os"
)

func processFile(path string) error {
    if !util.Exists(path) {
        return errors.New("文件不存在: " + path)
    }
    
    // 处理文件
    err := doSomething(path)
    if err != nil {
        if os.IsPermission(err) {
            return fmt.Errorf("权限不足: %w", err)
        }
        return fmt.Errorf("处理失败: %w", err)
    }
    
    return nil
}
```

## 类型定义

### ClnameConfig

```go
type ClnameConfig struct {
    Path  string  // 文件或目录路径
    DoTry bool    // 尝试模式，不实际修改
    Debug bool    // 调试模式
    Skip  bool    // 跳过错误
}
```

用于 clname 命令的配置。

### RenameConfig

```go
type RenameConfig struct {
    Debug      bool      // 调试模式
    DoTry      bool      // 尝试模式
    Formats    []string  // 文件格式列表
    Recursive  bool      // 递归搜索
    SourceDir  string    // 源目录
    OutputDir  string    // 输出目录
    Template   string    // 文件名模板
    StartIndex int       // 起始序号
}
```

用于 rename 命令的配置。

## 常量

### 默认值

```go
const (
    DefaultTemplate   = "file-@n"     // 默认文件名模板
    DefaultStartIndex = 1              // 默认起始序号
)
```

## 测试工具

### 测试辅助函数

```go
// 创建临时 EPUB 文件用于测试
func createTestEPUB(t *testing.T, title string) string {
    // 实现
}

// 创建临时目录
func createTestDir(t *testing.T) string {
    return t.TempDir()
}
```

## 最佳实践

### 1. 错误处理

始终检查和处理错误：

```go
if !util.Exists(path) {
    return fmt.Errorf("文件不存在: %s", path)
}
```

### 2. 资源清理

使用 defer 确保资源被释放：

```go
file, err := os.Open(path)
if err != nil {
    return err
}
defer file.Close()
```

### 3. 并发安全

如果在多个 goroutine 中使用，确保线程安全：

```go
var mu sync.Mutex

func safeProcess(path string) {
    mu.Lock()
    defer mu.Unlock()
    // 处理逻辑
}
```

### 4. 路径处理

使用 `filepath` 包处理路径：

```go
import "path/filepath"

// 跨平台路径拼接
fullPath := filepath.Join(dir, "file.epub")

// 获取扩展名
ext := filepath.Ext(path)
```

## 扩展建议

### 添加新的清理规则

在 `pkg/util/cleanname.go` 中添加新函数：

```go
func CleanTitleAdvanced(title string, rules []string) string {
    for _, rule := range rules {
        title = applyRule(title, rule)
    }
    return title
}
```

### 支持插件

定义插件接口：

```go
type Cleaner interface {
    Name() string
    SupportedFormats() []string
    Clean(path string) error
}
```

## 参考资源

- [Cobra 文档](https://github.com/spf13/cobra)
- [epub 库文档](https://github.com/kapmahc/epub)
- [Go 标准库](https://golang.org/pkg/)

---

最后更新: 2025-11-28

