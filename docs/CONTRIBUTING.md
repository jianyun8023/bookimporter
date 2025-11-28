# 贡献指南

感谢你对 BookImporter 项目的兴趣！我们欢迎各种形式的贡献。

## 目录

- [行为准则](#行为准则)
- [如何贡献](#如何贡献)
- [开发环境设置](#开发环境设置)
- [提交指南](#提交指南)
- [代码规范](#代码规范)
- [测试要求](#测试要求)
- [文档贡献](#文档贡献)

## 行为准则

### 我们的承诺

为了营造开放和友好的环境，我们承诺：

- 使用友好和包容的语言
- 尊重不同的观点和经验
- 优雅地接受建设性批评
- 关注对社区最有利的事情
- 对其他社区成员表示同理心

### 不可接受的行为

- 使用性暗示的语言或图像
- 发表侮辱性/贬损性评论，进行人身攻击
- 公开或私下骚扰
- 未经许可发布他人的私人信息
- 其他在专业场合可能被认为不当的行为

## 如何贡献

### 报告 Bug

提交 Bug 前，请：

1. 检查 [Issues](https://github.com/jianyun8023/bookimporter/issues) 确认问题未被报告
2. 使用最新版本确认问题仍存在
3. 收集相关信息

提交 Bug 时，请包含：

- **标题**: 简洁描述问题
- **环境信息**:
  - 操作系统和版本
  - Go 版本
  - BookImporter 版本
- **重现步骤**: 详细的步骤
- **期望行为**: 你期望发生什么
- **实际行为**: 实际发生了什么
- **错误信息**: 完整的错误输出
- **附加信息**: 日志、截图等

**模板:**

```markdown
## 环境
- OS: macOS 13.0
- Go: 1.20
- BookImporter: v0.1.0

## 重现步骤
1. 运行命令 `bookimporter clname -p test.epub`
2. ...

## 期望行为
应该清理标题

## 实际行为
程序崩溃

## 错误信息
```
panic: ...
```

## 附加信息
（截图、日志等）
```

### 建议功能

提交功能建议时：

1. 使用清晰描述性的标题
2. 详细描述建议的功能
3. 解释为什么这个功能有用
4. 提供使用示例（如果可能）

**模板:**

```markdown
## 功能描述
添加对 PDF 文件的支持

## 使用场景
用户经常需要清理 PDF 文件的元数据

## 建议实现
可以使用 pdfcpu 库...

## 替代方案
或者使用 ...
```

### 提交 Pull Request

1. Fork 项目
2. 创建特性分支
3. 进行更改
4. 提交 Pull Request

详细步骤见下文。

## 开发环境设置

### 前置要求

- Go 1.18 或更高版本
- Git
- make（可选）

### 克隆仓库

```bash
# Fork 项目后，克隆你的 fork
git clone https://github.com/YOUR_USERNAME/bookimporter.git
cd bookimporter

# 添加上游仓库
git remote add upstream https://github.com/jianyun8023/bookimporter.git
```

### 安装依赖

```bash
go mod download
```

### 构建项目

```bash
go build -o bookimporter
```

或使用 Makefile（如果可用）：

```bash
make build
```

### 运行测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./pkg/util

# 查看测试覆盖率
go test -cover ./...

# 生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 运行 Linter

```bash
# 安装 golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 运行 linter
golangci-lint run
```

## 提交指南

### 工作流程

1. **同步上游仓库**

```bash
git fetch upstream
git checkout main
git merge upstream/main
```

2. **创建特性分支**

```bash
git checkout -b feature/my-new-feature
```

分支命名规范：
- `feature/` - 新功能
- `fix/` - Bug 修复
- `docs/` - 文档更新
- `refactor/` - 代码重构
- `test/` - 测试相关

3. **进行更改**

- 编写代码
- 添加测试
- 更新文档

4. **提交更改**

```bash
git add .
git commit -m "feat: add PDF support"
```

提交信息规范见下文。

5. **推送到你的 Fork**

```bash
git push origin feature/my-new-feature
```

6. **创建 Pull Request**

在 GitHub 上创建 Pull Request，描述你的更改。

### 提交信息规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 格式：

```
<类型>[可选范围]: <描述>

[可选正文]

[可选脚注]
```

**类型:**

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式（不影响代码运行）
- `refactor`: 重构（既不是新功能也不是修复）
- `test`: 添加测试
- `chore`: 构建过程或辅助工具的变动

**示例:**

```
feat: add support for PDF files

Implement PDF metadata cleaning using pdfcpu library.

Closes #123
```

```
fix: handle empty EPUB titles

Previously crashed when title was empty, now returns error gracefully.
```

### Pull Request 要求

PR 应该：

1. **解决单一问题**: 每个 PR 专注于一个功能或修复
2. **包含测试**: 新功能或修复应有相应测试
3. **更新文档**: 如有必要，更新 README 或其他文档
4. **通过 CI**: 所有测试和 linter 检查都应通过
5. **清晰的描述**: 说明做了什么，为什么这么做

**PR 模板:**

```markdown
## 更改类型
- [ ] Bug 修复
- [ ] 新功能
- [ ] 重构
- [ ] 文档更新

## 描述
简要描述这个 PR 做了什么

## 相关 Issue
Closes #123

## 更改清单
- 添加了 X 功能
- 修复了 Y 问题
- 更新了 Z 文档

## 测试
- [ ] 添加了单元测试
- [ ] 手动测试通过
- [ ] 所有测试通过

## 截图（如果适用）

## 额外说明
```

## 代码规范

### Go 代码风格

遵循 [Effective Go](https://golang.org/doc/effective_go.html) 和 Go 社区最佳实践。

**基本规则:**

1. 使用 `gofmt` 格式化代码
2. 使用有意义的变量名
3. 添加必要的注释
4. 导出的函数必须有文档注释
5. 保持函数简短，单一职责
6. 正确处理错误

**示例:**

```go
// CleanTitle removes unwanted characters from book titles.
// It removes brackets and their content: ()[]（）【】
func CleanTitle(title string) string {
    // Implementation
}
```

### 错误处理

```go
// 好的做法
if err != nil {
    return fmt.Errorf("failed to process file %s: %w", file, err)
}

// 不好的做法
if err != nil {
    panic(err)  // 避免在库代码中使用 panic
}
```

### 命名规范

- **包名**: 简短、小写、单数
- **文件名**: 小写，使用下划线分隔
- **函数名**: 驼峰命名，导出函数首字母大写
- **变量名**: 驼峰命名，简短但有意义
- **常量**: 驼峰命名或全大写

### 代码组织

```go
package cmd

import (
    // 标准库
    "fmt"
    "os"
    
    // 第三方库
    "github.com/spf13/cobra"
    
    // 本地包
    "github.com/jianyun8023/bookimporter/pkg/util"
)

// 常量
const (
    DefaultTemplate = "file-@n"
)

// 变量
var (
    config *Config
)

// 类型定义
type Config struct {
    Path string
}

// 函数
func Execute() {
    // Implementation
}
```

## 测试要求

### 单元测试

为所有公共函数编写测试：

```go
func TestCleanTitle(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "remove chinese brackets",
            input:    "三体【刘慈欣】",
            expected: "三体",
        },
        {
            name:     "no brackets",
            input:    "三体",
            expected: "三体",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := CleanTitle(tt.input)
            if result != tt.expected {
                t.Errorf("got %q, want %q", result, tt.expected)
            }
        })
    }
}
```

### 测试覆盖率

- 新代码应有至少 80% 的测试覆盖率
- 核心功能应有更高的覆盖率

### 集成测试

对于复杂功能，考虑添加集成测试：

```go
func TestRenameIntegration(t *testing.T) {
    // 设置测试环境
    tempDir := t.TempDir()
    
    // 创建测试文件
    createTestFiles(tempDir)
    
    // 执行操作
    err := Rename(tempDir, config)
    
    // 验证结果
    assertFilesRenamed(t, tempDir)
}
```

## 文档贡献

文档同样重要！你可以：

1. **修正错误**: 拼写、语法、技术错误
2. **改进说明**: 让文档更清晰易懂
3. **添加示例**: 提供更多使用示例
4. **翻译**: 帮助翻译文档

### 文档规范

- 使用清晰简洁的语言
- 提供代码示例
- 保持格式一致
- 更新时同步修改相关文档

## 发布流程

（维护者专用）

1. 更新 `Changelog.md`
2. 更新版本号
3. 创建 Git tag
4. 构建发布包
5. 发布到 GitHub Releases

## 获取帮助

有问题？

- 查看 [FAQ](FAQ.md)
- 搜索 [Issues](https://github.com/jianyun8023/bookimporter/issues)
- 在 [Discussions](https://github.com/jianyun8023/bookimporter/discussions) 提问

## 致谢

感谢所有贡献者！

你的名字将出现在贡献者列表中。

---

最后更新: 2025-11-28

