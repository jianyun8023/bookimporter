# BookImporter

一个用 Go 语言编写的书籍导入助手工具，帮助你轻松管理和整理电子书库。

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://golang.org/dl/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](docs/LICENSE.md)

## ✨ 功能特性

### 1. 清理书籍标题 (clname)

自动清理 EPUB 书籍标题中的无用描述符，让你的电子书库更整洁。

- ✅ 移除各种括号标记：`（）【】()[]`
- ✅ 批量处理目录下所有 EPUB 文件
- ✅ 支持预览模式，安全可靠
- ✅ 使用 Calibre 的 ebook-meta 工具修改元数据
- ✅ 自动检测并处理损坏的 EPUB 文件

### 2. 批量重命名 (rename)

按照自定义模板批量重命名和移动文件，提高文件管理效率。

- ✅ 自定义文件名模板（支持序号占位符）
- ✅ 递归搜索子目录
- ✅ 移动文件到指定目录
- ✅ 多种文件格式支持
- ✅ 预览模式，确保操作正确

### 3. EPUB 文件完整性检测 (check)

检测 EPUB 文件是否损坏，帮助你维护健康的电子书库。

- ✅ ZIP 文件完整性验证
- ✅ 必需文件存在性检查（mimetype、container.xml 等）
- ✅ 元数据可解析性验证
- ✅ 批量检测，支持递归搜索
- ✅ 自动移动或删除损坏的文件
- ✅ 详细的错误报告和统计信息

## 🚀 快速开始

### 安装

**下载预编译版本**（推荐）

从 [Releases](https://github.com/jianyun8023/bookimporter/releases) 页面下载适合你系统的版本。

**从源码编译**

```bash
git clone https://github.com/jianyun8023/bookimporter.git
cd bookimporter
go build -o bookimporter
```

详细安装说明请查看 [安装指南](docs/INSTALLATION.md)。

### 基本使用

#### 清理书籍标题

```bash
# 清理单个文件
bookimporter clname -p book.epub

# 批量清理目录
bookimporter clname -p /path/to/books/

# 预览模式（不实际修改）
bookimporter clname -p /path/to/books/ -t

# 清理时自动移动损坏的文件
bookimporter clname -p /path/to/books/ --move-corrupted-to /path/to/corrupted/
```

#### 批量重命名

```bash
# 重命名当前目录的 txt 文件
bookimporter rename . -f txt -t "book-@n"

# 递归搜索并移动到输出目录
bookimporter rename /source/path -f epub -t "novel-@n" -r -o /output/path
```

#### 检测 EPUB 文件完整性

```bash
# 检测单个文件
bookimporter check -p book.epub

# 批量检测目录（递归）
bookimporter check -p /path/to/books/ -r

# 只显示有问题的文件
bookimporter check -p /path/to/books/ -r --only-errors

# 将损坏的文件移动到指定目录
bookimporter check -p /path/to/books/ -r --move-to /path/to/corrupted/

# 删除损坏的文件（会要求确认）
bookimporter check -p /path/to/books/ -r --delete

# 强制删除损坏的文件（不需要确认）
bookimporter check -p /path/to/books/ -r --delete --force
```

## 📚 文档

### 用户文档

- **[安装指南](docs/INSTALLATION.md)** - 详细的安装步骤
- **[使用指南](docs/USER_GUIDE.md)** - 完整的使用说明和示例
- **[常见问题](docs/FAQ.md)** - 常见问题解答

### 开发文档

- **[开发指南](CLAUDE.md)** - AI 开发助手指南
- **[API 文档](docs/API.md)** - 代码 API 参考
- **[架构设计](docs/ARCHITECTURE.md)** - 项目架构说明
- **[贡献指南](docs/CONTRIBUTING.md)** - 如何为项目贡献

### 项目管理

- **[变更日志](Changelog.md)** - 版本变更记录
- **[发布流程](docs/RELEASE.md)** - 版本发布流程
- **[文档中心](docs/README.md)** - 所有文档索引

## 💻 系统要求

- **操作系统**: macOS, Linux, Windows
- **Calibre**: clname 命令需要（用于修改 EPUB 元数据）

## 🤝 贡献

欢迎贡献！请查看 [贡献指南](docs/CONTRIBUTING.md) 了解如何参与项目。

## 📄 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](docs/LICENSE.md) 文件。

## 🙏 致谢

- [Cobra](https://github.com/spf13/cobra) - 强大的 CLI 框架
- [epub](https://github.com/kapmahc/epub) - EPUB 文件解析库
- [Calibre](https://calibre-ebook.com/) - 优秀的电子书管理软件

## 📞 联系方式

- **项目主页**: https://github.com/jianyun8023/bookimporter
- **问题反馈**: https://github.com/jianyun8023/bookimporter/issues
- **讨论区**: https://github.com/jianyun8023/bookimporter/discussions

---

如果这个项目对你有帮助，请给一个 ⭐️ Star！
