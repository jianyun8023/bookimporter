# 发布流程

本文档描述 BookImporter 项目的发布流程和规范。

## 目录

- [版本规范](#版本规范)
- [发布前检查清单](#发布前检查清单)
- [发布步骤](#发布步骤)
- [发布类型](#发布类型)
- [回滚流程](#回滚流程)
- [发布后任务](#发布后任务)

## 版本规范

BookImporter 遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/)。

### 版本号格式

```
主版本号.次版本号.修订号 (MAJOR.MINOR.PATCH)
```

**示例:**
- `1.0.0` - 第一个稳定版本
- `1.1.0` - 添加新功能
- `1.1.1` - Bug 修复
- `2.0.0` - 不兼容的 API 变更

### 版本号递增规则

1. **主版本号 (MAJOR)**: 不兼容的 API 变更
   - 删除或重命名命令
   - 修改命令参数含义
   - 删除公共 API

2. **次版本号 (MINOR)**: 向下兼容的功能性新增
   - 添加新命令
   - 添加新参数（可选）
   - 添加新的公共 API

3. **修订号 (PATCH)**: 向下兼容的问题修正
   - Bug 修复
   - 性能优化
   - 文档更新

### 先行版本

在正式版本前可以发布先行版本：

- `1.0.0-alpha` - 内部测试版本
- `1.0.0-beta` - 公开测试版本
- `1.0.0-rc.1` - 发布候选版本

## 发布前检查清单

### 代码质量

- [ ] 所有测试通过 (`go test ./...`)
- [ ] 测试覆盖率 >= 80%
- [ ] Linter 检查通过 (`golangci-lint run`)
- [ ] 代码已经过 Code Review
- [ ] 没有已知的严重 Bug

### 文档

- [ ] 更新 `Changelog.md`
- [ ] 更新 `README.md`（如有必要）
- [ ] 更新 API 文档（如有 API 变更）
- [ ] 更新用户指南（如有功能变更）
- [ ] 检查文档中的版本号
- [ ] 检查所有链接有效

### 版本信息

- [ ] 更新 `cmd/version.go` 中的版本号
- [ ] 确认 Git tag 与版本号一致

### 兼容性

- [ ] 验证向后兼容性
- [ ] 测试在不同操作系统上的运行（macOS, Linux, Windows）
- [ ] 测试在不同 Go 版本上的编译（Go 1.18+）

### 依赖

- [ ] 更新依赖到最新稳定版本
- [ ] 运行 `go mod tidy`
- [ ] 检查是否有安全漏洞 (`go list -m all | nancy sleuth`)

## 发布步骤

### 1. 准备发布分支

```bash
# 从 main 创建发布分支
git checkout main
git pull origin main
git checkout -b release/v1.0.0
```

### 2. 更新版本信息

#### 更新 version.go

```go
// cmd/version.go
const (
    Version = "1.0.0"
    BuildDate = "2025-11-28"
)
```

#### 更新 Changelog.md

```markdown
## [1.0.0] - 2025-11-28

### 新增
- 添加 clname 命令
- 添加 rename 命令

### 变更
- 优化性能

### 修复
- 修复路径处理问题
```

### 3. 提交变更

```bash
git add .
git commit -m "chore: prepare release v1.0.0"
git push origin release/v1.0.0
```

### 4. 创建 Pull Request

在 GitHub 上创建 PR，从 `release/v1.0.0` 到 `main`。

等待 CI 通过并进行 Code Review。

### 5. 合并到主分支

```bash
# PR 通过后合并
git checkout main
git merge release/v1.0.0
git push origin main
```

### 6. 创建 Git Tag

```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 7. 构建发布包

#### 使用 Makefile（推荐）

```bash
make release
```

#### 手动构建

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o bookimporter-linux-amd64

# macOS
GOOS=darwin GOARCH=amd64 go build -o bookimporter-darwin-amd64

# Windows
GOOS=windows GOARCH=amd64 go build -o bookimporter-windows-amd64.exe

# 压缩
tar -czf bookimporter-linux-amd64.tar.gz bookimporter-linux-amd64
tar -czf bookimporter-darwin-amd64.tar.gz bookimporter-darwin-amd64
zip bookimporter-windows-amd64.zip bookimporter-windows-amd64.exe
```

### 8. 创建 GitHub Release

1. 访问 https://github.com/jianyun8023/bookimporter/releases/new
2. 选择 tag: `v1.0.0`
3. 填写 Release title: `v1.0.0`
4. 填写 Release notes（从 Changelog 复制）
5. 上传构建的二进制文件
6. 如果是预发布版本，勾选 "This is a pre-release"
7. 点击 "Publish release"

#### Release Notes 模板

```markdown
## BookImporter v1.0.0

### 🎉 新增
- 添加 clname 命令，清理 EPUB 书籍标题
- 添加 rename 命令，批量重命名文件
- 支持递归搜索
- 支持预览模式

### 🔧 变更
- 优化文件处理性能
- 改进错误提示信息

### 🐛 修复
- 修复 Windows 平台路径问题
- 修复空标题导致的崩溃

### 📝 文档
- 添加详细的使用指南
- 添加常见问题文档
- 完善 API 文档

### 📦 下载

- **Linux**: [bookimporter-linux-amd64.tar.gz](...)
- **macOS**: [bookimporter-darwin-amd64.tar.gz](...)
- **Windows**: [bookimporter-windows-amd64.zip](...)

### 📋 完整更新日志

查看 [Changelog.md](https://github.com/jianyun8023/bookimporter/blob/main/Changelog.md)

### 🙏 致谢

感谢所有贡献者！
```

### 9. 发布到包管理器（可选）

#### Homebrew (macOS)

创建或更新 Homebrew formula。

#### apt/yum (Linux)

创建 deb/rpm 包。

## 发布类型

### 补丁版本 (Patch Release)

**场景**: Bug 修复、小的改进

**流程**:
1. 创建 `hotfix/` 分支
2. 修复问题
3. 更新版本号（增加 PATCH）
4. 按正常流程发布

**示例**:
```bash
git checkout -b hotfix/v1.0.1 main
# 修复 bug
git commit -m "fix: resolve crash on empty title"
# 继续发布流程
```

### 次版本 (Minor Release)

**场景**: 新功能、向后兼容的改进

**流程**:
同上述发布步骤。

### 主版本 (Major Release)

**场景**: 不兼容的变更

**特殊注意**:
1. 在 Release Notes 中明确说明不兼容变更
2. 提供迁移指南
3. 考虑提供兼容层（如果可能）

**示例 Release Notes**:

```markdown
## ⚠️ 重大变更

v2.0.0 包含不兼容的 API 变更：

### 变更内容

1. `clname` 命令的 `-t` 参数改为 `--dry-run`
2. 移除了 `--legacy` 选项

### 迁移指南

**旧版本:**
```bash
bookimporter clname -p books -t
```

**新版本:**
```bash
bookimporter clname -p books --dry-run
```

详细迁移指南: [MIGRATION.md](...)
```

### 预发布版本

**Alpha 版本**:
```bash
git tag -a v1.1.0-alpha.1 -m "Alpha release v1.1.0-alpha.1"
```

**Beta 版本**:
```bash
git tag -a v1.1.0-beta.1 -m "Beta release v1.1.0-beta.1"
```

**RC 版本**:
```bash
git tag -a v1.1.0-rc.1 -m "Release candidate v1.1.0-rc.1"
```

## 回滚流程

如果发现严重问题需要回滚：

### 1. 删除有问题的 Release

在 GitHub 上删除 Release（不删除 tag）。

### 2. 通知用户

在 Issues 和文档中通知用户不要使用该版本。

### 3. 快速修复

```bash
# 创建修复分支
git checkout -b hotfix/v1.0.2 v1.0.1

# 修复问题
git commit -m "fix: critical bug"

# 发布新版本
# 按正常流程发布 v1.0.2
```

### 4. 更新文档

在 Changelog 中记录：

```markdown
## [1.0.2] - 2025-11-29

### 修复
- 修复 v1.0.1 中的严重 Bug

## [1.0.1] - 2025-11-28 [已撤回]

此版本因严重 Bug 已撤回，请使用 v1.0.2。
```

## 发布后任务

### 立即执行

- [ ] 验证发布包可以下载
- [ ] 在不同平台测试安装和运行
- [ ] 更新项目网站（如有）
- [ ] 在社交媒体宣布发布

### 24小时内

- [ ] 监控 Issue tracker
- [ ] 回应用户反馈
- [ ] 修复紧急问题

### 一周内

- [ ] 收集用户反馈
- [ ] 规划下一个版本
- [ ] 更新路线图

## 自动化发布

### 使用 GitHub Actions

创建 `.github/workflows/release.yml`:

```yaml
name: Release

on:
  push:
    tags:
      - 'v*'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.18
      
      - name: Build
        run: |
          GOOS=linux GOARCH=amd64 go build -o bookimporter-linux-amd64
          GOOS=darwin GOARCH=amd64 go build -o bookimporter-darwin-amd64
          GOOS=windows GOARCH=amd64 go build -o bookimporter-windows-amd64.exe
      
      - name: Create Release
        uses: softprops/action-gh-release@v1
        with:
          files: |
            bookimporter-*
```

### 使用 GoReleaser

更专业的发布工具。

安装：
```bash
go install github.com/goreleaser/goreleaser@latest
```

配置 `.goreleaser.yml`:

```yaml
builds:
  - env:
      - CGO_ENABLED=0
    goos:
      - linux
      - darwin
      - windows
    goarch:
      - amd64
      - arm64

archives:
  - format: tar.gz
    format_overrides:
      - goos: windows
        format: zip

checksum:
  name_template: 'checksums.txt'

changelog:
  sort: asc
```

发布：
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
goreleaser release --clean
```

## 发布日历

建议的发布节奏：

- **主版本**: 每年 1-2 次
- **次版本**: 每 2-3 个月
- **补丁版本**: 按需发布（通常每月 1-2 次）

## 注意事项

1. **不要在周五发布**: 如果出问题，周末无法及时修复
2. **避免节假日**: 用户和开发者都不在线
3. **预留缓冲时间**: 发布前预留时间处理意外问题
4. **保持沟通**: 在 Issue/Discussion 中告知用户即将发布
5. **备份**: 确保代码和文档都有备份

## 联系信息

发布相关问题请联系：

- **维护者**: @jianyun8023
- **邮箱**: (如有)
- **讨论区**: https://github.com/jianyun8023/bookimporter/discussions

---

最后更新: 2025-11-28

