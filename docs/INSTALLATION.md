# 安装指南

本文档介绍如何在不同操作系统上安装 BookImporter。

## 系统要求

- **操作系统**: macOS, Linux, Windows
- **Go 版本**: 1.18 或更高（仅源码安装需要）
- **Calibre**: 用于 clname 命令（可选）

## 安装方式

### 方式一：下载预编译二进制文件

从 [GitHub Releases](https://github.com/jianyun8023/bookimporter/releases) 页面下载适合你系统的预编译版本。

#### macOS

```bash
# 下载（替换 VERSION 为实际版本号）
curl -LO https://github.com/jianyun8023/bookimporter/releases/download/vVERSION/bookimporter-darwin-amd64

# 添加执行权限
chmod +x bookimporter-darwin-amd64

# 移动到系统路径
sudo mv bookimporter-darwin-amd64 /usr/local/bin/bookimporter

# 验证安装
bookimporter version
```

#### Linux

```bash
# 下载（替换 VERSION 为实际版本号）
wget https://github.com/jianyun8023/bookimporter/releases/download/vVERSION/bookimporter-linux-amd64

# 添加执行权限
chmod +x bookimporter-linux-amd64

# 移动到系统路径
sudo mv bookimporter-linux-amd64 /usr/local/bin/bookimporter

# 验证安装
bookimporter version
```

#### Windows

1. 从 Releases 页面下载 `bookimporter-windows-amd64.exe`
2. 将文件重命名为 `bookimporter.exe`
3. 将文件移动到 PATH 环境变量包含的目录
4. 打开命令提示符，运行 `bookimporter version` 验证

### 方式二：从源码编译

#### 前置要求

确保已安装 Go 1.18+：

```bash
go version
```

#### 克隆并编译

```bash
# 克隆仓库
git clone https://github.com/jianyun8023/bookimporter.git
cd bookimporter

# 下载依赖
go mod download

# 编译
go build -o bookimporter

# 安装到系统路径（可选）
sudo mv bookimporter /usr/local/bin/

# 验证安装
bookimporter version
```

#### 使用 Makefile（如果可用）

```bash
# 编译
make build

# 安装
sudo make install

# 运行测试
make test
```

### 方式三：使用 Go Install

```bash
go install github.com/jianyun8023/bookimporter@latest
```

确保 `$GOPATH/bin` 或 `$HOME/go/bin` 在你的 PATH 中。

## 安装 Calibre（clname 命令依赖）

如果你需要使用 `clname` 命令清理 EPUB 书籍标题，需要安装 Calibre。

### macOS

使用 Homebrew：

```bash
brew install calibre
```

或者从 [Calibre 官网](https://calibre-ebook.com/download) 下载安装包。

### Linux

#### Ubuntu/Debian

```bash
sudo apt-get update
sudo apt-get install calibre
```

#### Fedora

```bash
sudo dnf install calibre
```

#### Arch Linux

```bash
sudo pacman -S calibre
```

### Windows

从 [Calibre 官网](https://calibre-ebook.com/download) 下载 Windows 安装程序并安装。

安装后，确保 `ebook-meta` 命令可用：

```bash
ebook-meta --version
```

## 验证安装

安装完成后，运行以下命令验证：

```bash
# 查看版本
bookimporter version

# 查看帮助
bookimporter --help

# 查看子命令帮助
bookimporter clname --help
bookimporter rename --help
```

## 配置（可选）

### Shell 自动补全

#### Bash

```bash
bookimporter completion bash > /etc/bash_completion.d/bookimporter
```

#### Zsh

```bash
bookimporter completion zsh > "${fpath[1]}/_bookimporter"
```

#### Fish

```bash
bookimporter completion fish > ~/.config/fish/completions/bookimporter.fish
```

### 环境变量

可以设置以下环境变量：

```bash
# Calibre ebook-meta 工具路径（如果不在 PATH 中）
export EBOOK_META_PATH=/path/to/ebook-meta
```

## 升级

### 预编译版本

重新下载最新版本并替换旧文件。

### 源码编译

```bash
cd bookimporter
git pull origin main
go build -o bookimporter
sudo mv bookimporter /usr/local/bin/
```

### Go Install

```bash
go install github.com/jianyun8023/bookimporter@latest
```

## 卸载

### 删除二进制文件

```bash
sudo rm /usr/local/bin/bookimporter
```

### 删除配置文件（如有）

```bash
rm -rf ~/.bookimporter
```

## 故障排除

### 命令未找到

确保安装目录在 PATH 中：

```bash
echo $PATH
```

### Permission Denied

添加执行权限：

```bash
chmod +x bookimporter
```

### ebook-meta 未找到

确保 Calibre 已正确安装并在 PATH 中：

```bash
which ebook-meta
```

### Go 版本过低

升级 Go 到 1.18 或更高版本：

```bash
# 使用官方安装包或包管理器升级
# 参考: https://golang.org/doc/install
```

## 获取帮助

如果遇到安装问题，请：

1. 查看 [常见问题](FAQ.md)
2. 搜索 [GitHub Issues](https://github.com/jianyun8023/bookimporter/issues)
3. 提交新的 Issue

---

最后更新: 2025-11-28

