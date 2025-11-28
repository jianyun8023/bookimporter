# 使用指南

本指南详细介绍 BookImporter 的所有功能和使用方法。

## 目录

- [基础概念](#基础概念)
- [clname 命令](#clname-命令)
- [rename 命令](#rename-命令)
- [高级用法](#高级用法)
- [最佳实践](#最佳实践)

## 基础概念

BookImporter 是一个命令行工具，提供两个主要命令：

- **clname**: 清理 EPUB 书籍标题中的无用描述符
- **rename**: 批量重命名和移动文件

### 通用选项

所有命令都支持以下选项：

```bash
-h, --help      # 显示帮助信息
```

## clname 命令

清理 EPUB 书籍标题中的括号和其他无用标记。

### 语法

```bash
bookimporter clname [选项]
```

### 选项

| 选项 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| --path | -p | ./ | 目标文件或目录路径 |
| --dotry | -t | false | 尝试运行，不实际修改 |
| --skip | -j | false | 跳过无法解析的书籍 |
| --debug | -d | false | 启用调试模式 |

### 使用示例

#### 1. 清理单个文件

```bash
bookimporter clname -p /path/to/book.epub
```

**输出示例:**
```
路径: /path/to/book.epub 新名称: 三体 旧名称: 三体【刘慈欣】
```

#### 2. 批量清理目录

```bash
bookimporter clname -p /path/to/books/
```

程序会自动处理目录下所有 `.epub` 文件。

#### 3. 预览模式（不实际修改）

```bash
bookimporter clname -p /path/to/books/ -t
```

这会显示将要进行的修改，但不会实际修改文件。

#### 4. 跳过错误

```bash
bookimporter clname -p /path/to/books/ -j
```

遇到无法解析的 EPUB 文件时继续处理其他文件，而不是终止。

#### 5. 调试模式

```bash
bookimporter clname -p /path/to/books/ -d
```

显示详细的调试信息。

### 清理规则

clname 命令会移除以下标记：

- 中文括号：`（内容）`、`【内容】`
- 英文括号：`(content)`、`[content]`

**示例转换:**

| 原标题 | 清理后 |
|--------|--------|
| 三体【刘慈欣】 | 三体 |
| Python编程（第3版） | Python编程 |
| 数据结构[C语言版] | 数据结构 |
| 哈利波特(全集) | 哈利波特 |

### 注意事项

1. **依赖 Calibre**: 此命令需要安装 Calibre 的 `ebook-meta` 工具
2. **仅支持 EPUB**: 目前只支持 EPUB 格式
3. **元数据修改**: 修改的是 EPUB 文件内部的元数据，不是文件名
4. **备份建议**: 首次使用建议先备份文件或使用 `-t` 选项预览

## rename 命令

按照模板批量重命名和移动文件。

### 语法

```bash
bookimporter rename <源目录> [选项]
```

### 选项

| 选项 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| --format | -f | * | 文件格式过滤（可多次使用） |
| --template | -t | file-@n | 文件名模板 |
| --recursive | -r | false | 递归搜索子目录 |
| --output | -o | 无 | 输出目录（移动文件） |
| --start-num | 无 | 1 | 起始序号 |
| --do-try | 无 | false | 预览模式，不实际执行 |
| --debug | 无 | false | 显示调试信息 |

### 使用示例

#### 1. 基础重命名

将当前目录的 txt 文件重命名为 `book-1.txt`, `book-2.txt`, ...

```bash
bookimporter rename . -f txt -t "book-@n"
```

#### 2. 指定起始序号

从 100 开始编号：

```bash
bookimporter rename . -f pdf -t "doc-@n" --start-num 100
```

输出: `doc-100.pdf`, `doc-101.pdf`, ...

#### 3. 移动并重命名

将文件移动到新目录并重命名：

```bash
bookimporter rename /source/dir -f epub -t "novel-@n" -o /output/dir
```

#### 4. 递归搜索

搜索子目录：

```bash
bookimporter rename . -f txt -t "file-@n" -r
```

#### 5. 多种格式

处理多种文件格式：

```bash
bookimporter rename . -f epub -f pdf -f mobi -t "book-@n"
```

#### 6. 预览模式

查看将要执行的操作：

```bash
bookimporter rename . -f txt -t "file-@n" --do-try
```

**输出示例:**
```
3 files found.
3 files renamed:
  - ./old1.txt -> ./file-1.txt
  - ./old2.txt -> ./file-2.txt
  - ./old3.txt -> ./file-3.txt
```

### 模板语法

- **@n**: 序号占位符，会被替换为实际序号
- 文件扩展名自动保留

**模板示例:**

| 模板 | 结果 |
|------|------|
| `book-@n` | `book-1.epub`, `book-2.epub` |
| `novel-@n-zh` | `novel-1-zh.pdf`, `novel-2-zh.pdf` |
| `第@n章` | `第1章.txt`, `第2章.txt` |
| `@n` | `1.epub`, `2.epub` |

### 高级示例

#### 整理下载的书籍

```bash
# 将 Downloads 中的 epub 文件移动到 Books 目录并统一命名
bookimporter rename ~/Downloads -f epub -t "imported-@n" -o ~/Books
```

#### 批量整理章节文件

```bash
# 将小说章节文件重命名为统一格式
bookimporter rename ./novel -f txt -t "chapter-@n" -r --start-num 1
```

#### 多格式电子书整理

```bash
# 整理各种格式的电子书
bookimporter rename ~/ebooks -f epub -f pdf -f mobi -f azw3 -t "book-@n" -o ~/Library/Books
```

## 高级用法

### 结合 Shell 脚本

#### 批量处理多个目录

```bash
#!/bin/bash
for dir in ~/Books/*/; do
    echo "Processing $dir"
    bookimporter clname -p "$dir" -j
done
```

#### 处理前自动备份

```bash
#!/bin/bash
SOURCE_DIR="$1"
BACKUP_DIR="${SOURCE_DIR}_backup_$(date +%Y%m%d_%H%M%S)"

# 备份
cp -r "$SOURCE_DIR" "$BACKUP_DIR"

# 处理
bookimporter clname -p "$SOURCE_DIR"

echo "备份保存在: $BACKUP_DIR"
```

### 工作流示例

#### 完整的书籍导入流程

```bash
# 1. 预览将要重命名的文件
bookimporter rename ~/Downloads -f epub -t "new-@n" --do-try

# 2. 确认无误后，移动到书库目录
bookimporter rename ~/Downloads -f epub -t "new-@n" -o ~/Books/ToProcess

# 3. 清理书籍标题
bookimporter clname -p ~/Books/ToProcess -j

# 4. 移动到最终位置
mv ~/Books/ToProcess/* ~/Books/Library/
```

## 最佳实践

### 1. 使用预览模式

在实际执行前，总是先使用预览模式：

```bash
# clname 使用 -t
bookimporter clname -p /path/to/books -t

# rename 使用 --do-try
bookimporter rename . -f txt -t "new-@n" --do-try
```

### 2. 定期备份

在批量操作前备份重要文件：

```bash
tar -czf books_backup_$(date +%Y%m%d).tar.gz ~/Books/
```

### 3. 使用 skip 选项

处理大量文件时使用 `-j` 跳过错误：

```bash
bookimporter clname -p /large/library -j > errors.log 2>&1
```

### 4. 分批处理

对于大量文件，分批处理更安全：

```bash
# 处理前 100 个文件
find . -name "*.epub" | head -100 | xargs -I {} bookimporter clname -p {}
```

### 5. 保持命名一致

使用统一的命名模板：

```bash
# 图书：book-001, book-002
bookimporter rename . -f epub -t "book-@n" --start-num 1

# 论文：paper-001, paper-002
bookimporter rename . -f pdf -t "paper-@n" --start-num 1
```

### 6. 记录操作日志

```bash
bookimporter clname -p ~/Books | tee cleanup_$(date +%Y%m%d).log
```

## 常见问题

参见 [FAQ.md](FAQ.md)

## 获取帮助

- 命令帮助: `bookimporter [command] --help`
- GitHub Issues: https://github.com/jianyun8023/bookimporter/issues
- 项目文档: https://github.com/jianyun8023/bookimporter/docs

---

最后更新: 2025-11-28

