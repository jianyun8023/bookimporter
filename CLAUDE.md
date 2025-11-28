# BookImporter - AI 开发指南

## 项目概述

BookImporter 是一个用 Go 语言开发的书籍导入助手工具，旨在帮助用户管理和整理电子书库。

### 核心功能

1. **书籍标题清理 (clname)**
   - 自动清理 EPUB 书籍标题中的无用描述符
   - 移除常见的括号标记：`（）【】()[]`
   - 使用 calibre 的 `ebook-meta` 工具修改元数据
   - 支持批量处理目录中的所有 EPUB 文件

2. **文件批量重命名 (rename)**
   - 按自定义模板批量重命名文件
   - 支持递归搜索子目录
   - 支持移动文件到指定输出目录
   - 序列号自动编号功能

## 技术栈

- **语言**: Go 1.18+
- **CLI框架**: [Cobra](https://github.com/spf13/cobra) - 强大的命令行应用框架
- **EPUB处理**: [epub](https://github.com/kapmahc/epub) - EPUB 文件解析库
- **外部依赖**: Calibre 的 ebook-meta 工具（用于修改 EPUB 元数据）

## 项目结构

```
bookimporter/
├── cmd/                    # 命令定义
│   ├── root.go            # 根命令
│   ├── clname.go          # 清理标题命令
│   ├── rename.go          # 重命名命令
│   └── version.go         # 版本命令
├── pkg/                   # 包和工具
│   └── util/              # 工具函数
│       ├── cleanname.go   # 标题清理逻辑
│       ├── cleanname_test.go
│       └── filetool.go    # 文件操作工具
├── docs/                  # 文档目录
├── main.go               # 程序入口
├── go.mod                # Go 模块定义
└── README.md             # 用户文档
```

## 命令使用示例

### 1. 清理书籍标题 (clname)

```bash
# 清理单个文件
bookimporter clname -p /path/to/book.epub

# 批量清理目录下所有 EPUB 文件
bookimporter clname -p /path/to/books/

# 尝试运行（不实际修改）
bookimporter clname -p /path/to/books/ -t

# 跳过无法解析的书籍
bookimporter clname -p /path/to/books/ -j

# 调试模式
bookimporter clname -p /path/to/books/ -d
```

**参数说明:**
- `-p, --path`: 目录或文件路径
- `-t, --dotry`: 尝试运行，不实际修改
- `-j, --skip`: 跳过无法解析的书籍
- `-d, --debug`: 调试模式

### 2. 批量重命名 (rename)

```bash
# 基本用法：重命名当前目录下的 txt 文件
bookimporter rename . -f txt -t "book-@n"

# 递归搜索并移动到输出目录
bookimporter rename /source/path -f epub -t "novel-@n" -r -o /output/path

# 从指定序号开始
bookimporter rename . -f pdf -t "doc-@n" --start-num 100

# 尝试运行（预览操作）
bookimporter rename . -f txt -t "file-@n" --do-try
```

**参数说明:**
- `-f, --format`: 文件格式过滤（可多次指定）
- `-t, --template`: 文件名模板，`@n` 为序号占位符
- `-r, --recursive`: 递归搜索子目录
- `-o, --output`: 输出目录（移动文件）
- `--start-num`: 起始序号（默认为1）
- `--do-try`: 仅预览，不实际执行
- `--debug`: 显示调试信息

## 开发指南

### 本地构建

```bash
# 克隆项目
git clone https://github.com/jianyun8023/bookimporter.git
cd bookimporter

# 安装依赖
go mod download

# 构建
go build -o bookimporter

# 运行测试
go test ./...

# 查看测试覆盖率
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 添加新命令

1. 在 `cmd/` 目录创建新的命令文件
2. 定义 cobra.Command 结构
3. 在 `cmd/root.go` 的 `init()` 函数中注册命令

示例：
```go
var newCmd = &cobra.Command{
    Use:   "newcmd",
    Short: "新命令的简短描述",
    Run: func(cmd *cobra.Command, args []string) {
        // 命令逻辑
    },
}

func init() {
    rootCmd.AddCommand(newCmd)
}
```

### 代码规范

- 遵循 Go 官方代码风格
- 使用 `gofmt` 格式化代码
- 为公共函数添加注释
- 编写单元测试
- 错误处理要完善

## 依赖说明

### 运行时依赖

- **Calibre ebook-meta**: clname 命令需要此工具来修改 EPUB 元数据
  - macOS: `brew install calibre`
  - Linux: `sudo apt-get install calibre`
  - Windows: 从 [Calibre 官网](https://calibre-ebook.com/download) 下载安装

### Go 模块依赖

- `github.com/spf13/cobra`: CLI 框架
- `github.com/kapmahc/epub`: EPUB 文件解析
- `github.com/spf13/pflag`: 命令行参数解析

## 常见问题

### clname 命令相关

**Q: 为什么执行失败，提示找不到 ebook-meta？**
A: 需要安装 Calibre 软件包，它包含了 ebook-meta 命令行工具。

**Q: 如何自定义清理规则？**
A: 修改 `pkg/util/cleanname.go` 中的 `TryCleanTitle` 函数，添加自定义的清理逻辑。

### rename 命令相关

**Q: 模板中的 @n 是什么？**
A: @n 是序号占位符，会被替换为实际的序列号。

**Q: 如何处理多种文件格式？**
A: 可以多次使用 `-f` 参数，例如：`-f epub -f pdf -f mobi`

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [docs/LICENSE.md](docs/LICENSE.md)

## 联系方式

- 项目地址: https://github.com/jianyun8023/bookimporter
- 问题反馈: https://github.com/jianyun8023/bookimporter/issues

## AI 助手提示

在协助开发此项目时，请注意：

1. **保持代码简洁**: Go 语言强调简洁和可读性
2. **错误处理**: 确保所有可能的错误都被妥善处理
3. **测试覆盖**: 为新功能编写相应的测试用例
4. **向后兼容**: 修改现有功能时保持向后兼容性
5. **文档更新**: 功能变更时同步更新文档
6. **性能考虑**: 处理大量文件时注意性能优化

