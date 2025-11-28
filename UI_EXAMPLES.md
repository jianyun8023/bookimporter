# BookImporter 新 UI 使用示例

## 快速开始

### 1. 检测 EPUB 文件完整性

```bash
# 检测单个文件
./bookimporter check -p /path/to/book.epub

# 检测目录中的所有文件
./bookimporter check -p /path/to/books/

# 递归检测（包括子目录）
./bookimporter check -p /path/to/books/ -r

# 只显示有问题的文件
./bookimporter check -p /path/to/books/ --only-errors
```

**新 UI 特性展示：**
- ✓ 彩色的命令头部
- ✓ 实时进度显示
- ✓ 统计表格总结
- ✓ 清晰的成功/失败标识

### 2. 清理书籍标题

```bash
# 清理单个文件
./bookimporter clname -p /path/to/book.epub

# 批量清理目录
./bookimporter clname -p /path/to/books/

# 试运行模式（预览但不执行）
./bookimporter clname -p /path/to/books/ -t

# 跳过无法解析的文件
./bookimporter clname -p /path/to/books/ -j
```

**新 UI 特性展示：**
- ✓ 清晰的 old → new 标题对比
- ✓ 批量处理进度条
- ✓ 样式化的文件路径显示
- ✓ 彩色的统计报告

### 3. 批量重命名文件

```bash
# 基本重命名
./bookimporter rename /path/to/files -f txt -t "book-@n"

# 预览模式（不实际执行）
./bookimporter rename /path/to/files -f txt -t "book-@n" --do-try

# 递归搜索并移动到新目录
./bookimporter rename /path/to/files -f epub -t "novel-@n" -r -o /output/path

# 自定义起始序号
./bookimporter rename /path/to/files -f pdf -t "doc-@n" --start-num 100
```

**新 UI 特性展示：**
- ✓ 预览模式的表格展示
- ✓ 重命名操作的进度跟踪
- ✓ 文件变更的清晰对比
- ✓ 完成后的统计汇总

## UI 元素说明

### 状态指示符

- `✓` 绿色 - 成功/通过
- `✗` 红色 - 错误/失败
- `→` 青色 - 信息/操作
- `⚠` 黄色 - 警告
- `⋯` 灰色 - 跳过/无需处理

### 颜色含义

- **绿色** - 成功、完成、新值
- **红色** - 错误、失败
- **黄色** - 警告、旧值（删除线）
- **青色** - 信息、提示
- **蓝色** - 标题、高亮
- **灰色** - 次要信息、路径

### 进度显示

批量操作时会显示实时进度：

```
进度: 75% (15/20)
```

### 统计表格

操作完成后会显示美观的统计表格：

```
┌────────────┬───────┐
│   状态     │ 数量  │
├────────────┼───────┤
│ ✓ 通过     │   18  │
│ ✗ 失败     │   2   │
└────────────┴───────┘
```

## 高级用法

### 处理损坏的文件

```bash
# 检测并移动损坏的文件
./bookimporter check -p /path/to/books/ --move-to /corrupted/

# 删除损坏的文件（需要确认）
./bookimporter check -p /path/to/books/ --delete

# 强制删除（不需要确认）
./bookimporter check -p /path/to/books/ --delete --force
```

**删除确认提示**会显示美观的警告框：
```
╭────────────────────────────────────╮
│  ⚠ 删除确认                        │
│                                    │
│  即将删除文件:                     │
│  /path/to/file.epub                │
│                                    │
│  此操作无法撤销！                  │
╰────────────────────────────────────╯

确认删除? (y/N): 
```

### 调试模式

所有命令都支持调试模式，显示详细信息：

```bash
./bookimporter check -p /path/to/books/ -d
./bookimporter clname -p /path/to/books/ -d
./bookimporter rename /path/to/files -f txt -t "book-@n" --debug
```

## 性能说明

- UI 组件基于 Lipgloss，性能开销极小
- 进度显示采用单行刷新，不会产生大量输出
- 批量操作时的 UI 更新频率经过优化

## 兼容性

- **终端支持**: 支持所有现代终端（iTerm2, Terminal, Windows Terminal, Alacritty 等）
- **颜色支持**: 自动检测终端颜色能力
- **字符集**: 使用 Unicode 字符（✓ ✗ → 等）
- **降级方案**: 在不支持颜色的终端中自动降级为纯文本

## 提示和技巧

1. **预览模式**: 建议在执行批量操作前先使用 `--do-try` 或 `-t` 预览
2. **递归搜索**: 使用 `-r` 可以处理嵌套目录结构
3. **过滤输出**: 使用 `--only-errors` 只关注问题文件
4. **保护操作**: 删除文件默认需要确认，使用 `--force` 跳过

## 反馈和贡献

如果你发现任何 UI 问题或有改进建议，欢迎提交 Issue 或 Pull Request！

