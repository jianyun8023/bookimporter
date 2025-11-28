# Changelog

本文档记录 BookImporter 项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [Unreleased]

### 新增
- 创建 CLAUDE.md 文档，为 AI 助手提供项目开发指南
- 创建 Changelog.md，规范版本变更记录
- 创建 docs 目录，整理项目文档结构

### 变更
- 将 LICENSE 移至 docs/LICENSE.md
- 将 SECURITY.md 移至 docs/SECURITY.md

## [0.1.0] - 2023-XX-XX

### 新增
- 实现 `clname` 命令：清理 EPUB 书籍标题中的无用描述符
  - 支持单个文件和批量处理
  - 支持尝试运行模式（--dotry）
  - 支持跳过错误（--skip）
  - 集成 Calibre 的 ebook-meta 工具
  
- 实现 `rename` 命令：批量重命名和移动文件
  - 支持自定义文件名模板
  - 支持递归搜索子目录
  - 支持移动文件到指定目录
  - 支持自定义起始序号
  - 支持多种文件格式过滤
  - 支持预览模式（--do-try）

- 实现 `version` 命令：显示程序版本信息

- 添加工具函数库
  - EPUB 标题清理逻辑
  - 文件操作工具函数
  - 单元测试覆盖

### 技术细节
- 使用 Go 1.18+ 开发
- 集成 Cobra CLI 框架
- 集成 epub 解析库
- 完善的命令行参数支持

## 版本说明

### 版本号规则
- **主版本号 (Major)**: 不兼容的 API 变更
- **次版本号 (Minor)**: 向下兼容的功能性新增
- **修订号 (Patch)**: 向下兼容的问题修正

### 变更类型
- **新增 (Added)**: 新功能
- **变更 (Changed)**: 对现有功能的变更
- **弃用 (Deprecated)**: 即将移除的功能
- **移除 (Removed)**: 已移除的功能
- **修复 (Fixed)**: 任何 bug 修复
- **安全 (Security)**: 修复安全问题

---

## 如何维护此文档

### 添加新版本
1. 在 `[Unreleased]` 下方添加新版本标题
2. 使用格式：`## [版本号] - YYYY-MM-DD`
3. 按类型组织变更内容

### 记录变更
在开发过程中，将变更添加到 `[Unreleased]` 部分：

```markdown
## [Unreleased]

### 新增
- 添加了某某新功能

### 修复
- 修复了某某问题
```

### 发布新版本
发布时，将 `[Unreleased]` 的内容移到新版本标题下：

```markdown
## [Unreleased]

## [1.0.0] - 2025-11-28

### 新增
- 添加了某某新功能

### 修复
- 修复了某某问题
```

### 示例

```markdown
## [Unreleased]

### 新增
- 添加对 PDF 文件的元数据清理支持

## [1.1.0] - 2025-12-01

### 新增
- 添加 `export` 命令，支持导出书籍列表
- 支持 JSON 和 CSV 格式导出

### 变更
- 优化 `rename` 命令的性能
- 改进错误提示信息

### 修复
- 修复 Windows 平台路径分隔符问题
- 修复大文件处理时的内存泄漏

## [1.0.0] - 2025-11-28

### 新增
- 首个稳定版本发布
```

---

最后更新: 2025-11-28

