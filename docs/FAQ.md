# 常见问题 (FAQ)

本文档汇总了 BookImporter 使用过程中的常见问题和解决方案。

## 目录

- [安装相关](#安装相关)
- [clname 命令](#clname-命令)
- [rename 命令](#rename-命令)
- [性能相关](#性能相关)
- [平台相关](#平台相关)

## 安装相关

### Q1: 如何验证安装是否成功？

运行以下命令：

```bash
bookimporter version
```

如果显示版本信息，说明安装成功。

### Q2: 提示 "command not found: bookimporter"

**可能原因:**
1. 未正确安装
2. 安装路径不在 PATH 环境变量中

**解决方案:**

检查文件是否存在：

```bash
ls -l /usr/local/bin/bookimporter
```

检查 PATH：

```bash
echo $PATH
```

如果 `/usr/local/bin` 不在 PATH 中，添加它：

```bash
export PATH=$PATH:/usr/local/bin
```

或将可执行文件移动到已在 PATH 中的目录：

```bash
sudo mv bookimporter /usr/bin/
```

### Q3: macOS 提示 "无法打开，因为无法验证开发者"

**解决方案:**

```bash
# 移除隔离属性
xattr -d com.apple.quarantine /usr/local/bin/bookimporter

# 或在系统偏好设置 > 安全性与隐私中允许
```

### Q4: 如何更新到最新版本？

**使用二进制文件:**

下载最新版本并替换旧文件。

**使用源码:**

```bash
cd bookimporter
git pull origin main
go build -o bookimporter
sudo mv bookimporter /usr/local/bin/
```

## clname 命令

### Q5: 提示 "ebook-meta: command not found"

**原因:** 未安装 Calibre 或 ebook-meta 不在 PATH 中。

**解决方案:**

安装 Calibre：

```bash
# macOS
brew install calibre

# Ubuntu/Debian
sudo apt-get install calibre

# 验证安装
which ebook-meta
ebook-meta --version
```

### Q6: 处理时提示 "无法获得书籍标题"

**可能原因:**
1. EPUB 文件损坏
2. EPUB 格式不标准
3. 文件没有标题元数据

**解决方案:**

1. 程序默认会跳过错误文件继续处理其他文件：

```bash
bookimporter clname -p /path/to/books -r
```

如果希望即使有失败也返回退出码 0，可以使用 `-i` 选项：

```bash
bookimporter clname -p /path/to/books -r -i
```

2. 使用 Calibre 手动检查文件：

```bash
ebook-meta problem.epub
```

3. 如果文件确实损坏，尝试重新下载或使用 Calibre 修复。

### Q7: 清理后标题没有变化？

**可能原因:**
1. 标题本来就没有括号等标记
2. 使用了 `-t` 预览模式

**解决方案:**

确保没有使用 `-t` 选项：

```bash
bookimporter clname -p /path/to/book.epub
```

检查原标题是否确实包含需要清理的字符。

### Q8: 能否自定义清理规则？

目前不支持通过命令行参数自定义，但可以修改源码。

编辑 `pkg/util/cleanname.go` 中的 `TryCleanTitle` 函数：

```go
func TryCleanTitle(title string) string {
    // 添加自定义规则
    title = strings.ReplaceAll(title, "custom", "")
    return title
}
```

然后重新编译。

### Q9: clname 会修改文件名吗？

**不会**。clname 只修改 EPUB 文件内部的元数据（标题字段），不会改变文件名。

如果要修改文件名，需要使用其他工具或 `rename` 命令。

### Q10: 批量处理时如何查看处理进度？

使用管道和 grep 统计：

```bash
bookimporter clname -p /path/to/books | tee process.log
```

或编写脚本显示进度：

```bash
#!/bin/bash
total=$(find /path/to/books -name "*.epub" | wc -l)
count=0
for file in /path/to/books/*.epub; do
    count=$((count+1))
    echo "[$count/$total] Processing $file"
    bookimporter clname -p "$file"
done
```

## rename 命令

### Q11: 模板中的 @n 是什么？

`@n` 是序号占位符，会被替换为实际的序列号。

例如，模板 `book-@n` 会生成：
- `book-1.epub`
- `book-2.epub`
- `book-3.epub`

### Q12: 如何保持原文件扩展名？

扩展名会自动保留，无需在模板中指定。

例如：
- 原文件: `old.txt`
- 模板: `new-@n`
- 结果: `new-1.txt`

### Q13: 可以使用多位数编号吗（如 001, 002）？

目前不直接支持，但可以通过脚本实现：

```bash
# 使用 printf 格式化
for i in {1..100}; do
    printf "book-%03d\n" $i
done
```

或提交 Feature Request 建议添加此功能。

### Q14: rename 会覆盖已存在的文件吗？

如果目标文件名已存在，程序会报错并停止。建议先使用 `--do-try` 预览。

### Q15: 如何批量重命名特定目录下的文件？

使用路径参数指定目录：

```bash
bookimporter rename /specific/directory -f txt -t "file-@n"
```

如果需要递归处理子目录，添加 `-r` 选项：

```bash
bookimporter rename /specific/directory -f txt -t "file-@n" -r
```

### Q16: 文件顺序是怎样的？

文件按照文件系统返回的顺序处理，通常是按名称排序。

如果需要特定顺序，建议先用其他工具排序，然后逐个处理。

### Q17: 可以只移动文件不重命名吗？

可以，使用原文件名作为模板（需要脚本实现）。

或者直接使用 `mv` 命令：

```bash
mv *.epub /destination/directory/
```

## 性能相关

### Q18: 处理大量文件很慢怎么办？

**优化建议:**

1. 分批处理：

```bash
find . -name "*.epub" | head -1000 | xargs -I {} bookimporter clname -p {}
```

2. 使用并行处理（GNU Parallel）：

```bash
find . -name "*.epub" | parallel bookimporter clname -p {}
```

3. 使用 SSD 而非机械硬盘

### Q19: 处理大文件时内存占用很高？

这是正常现象，特别是处理 EPUB 文件时。

**缓解措施:**
- 分批处理
- 关闭其他占用内存的程序
- 程序默认会跳过问题文件继续处理

### Q20: 可以后台运行吗？

可以，使用 `nohup` 或 `screen`：

```bash
nohup bookimporter clname -p /large/library -r > process.log 2>&1 &
```

或使用 screen：

```bash
screen -S bookimporter
bookimporter clname -p /large/library -r
# 按 Ctrl+A 然后 D 分离会话
```

## 平台相关

### Q21: Windows 上如何使用？

在 PowerShell 或 CMD 中：

```cmd
bookimporter.exe clname -p C:\Books
```

路径使用反斜杠或正斜杠都可以。

### Q22: Windows 上 clname 命令不工作？

Windows 版 Calibre 的 `ebook-meta.exe` 可能不在 PATH 中。

**解决方案:**

1. 找到 Calibre 安装路径（通常是 `C:\Program Files\Calibre2`）

2. 添加到 PATH 或使用完整路径：

```cmd
set PATH=%PATH%;C:\Program Files\Calibre2
```

### Q23: macOS Big Sur 及以上版本权限问题？

**解决方案:**

授予终端完全磁盘访问权限：

1. 系统偏好设置 > 安全性与隐私 > 隐私
2. 选择"完全磁盘访问权限"
3. 添加终端应用

### Q24: Linux 上提示权限不足？

**解决方案:**

1. 使用 `chmod` 添加执行权限：

```bash
chmod +x bookimporter
```

2. 或使用 `sudo`：

```bash
sudo bookimporter clname -p /path/to/books
```

## 其他问题

### Q25: 支持哪些电子书格式？

**clname 命令:** 仅支持 EPUB

**rename 命令:** 支持所有文件格式，使用 `-f` 参数指定

### Q26: 会损坏原文件吗？

正常情况下不会，但建议：
1. 首次使用前备份重要文件
2. 使用预览模式（`-t` 或 `--do-try`）测试
3. 分批处理，先测试小批量

### Q27: 如何报告 Bug？

1. 访问 https://github.com/jianyun8023/bookimporter/issues
2. 搜索是否已有类似问题
3. 如果没有，创建新 Issue，包含：
   - 系统信息（OS、版本）
   - 完整的命令
   - 错误信息
   - 重现步骤

### Q28: 如何贡献代码？

参见 [CONTRIBUTING.md](CONTRIBUTING.md)

### Q29: 有 GUI 版本吗？

目前只有命令行版本。GUI 版本在计划中，欢迎贡献。

### Q30: 支持多语言书籍吗？

支持，只要是标准的 EPUB 格式，无论书籍内容是什么语言都可以处理。

清理规则对中文和英文的括号都有效。

## 还有问题？

如果你的问题没有在此列出：

1. 查看 [使用指南](USER_GUIDE.md)
2. 搜索 [GitHub Issues](https://github.com/jianyun8023/bookimporter/issues)
3. 提交新的 Issue
4. 加入讨论区交流

---

最后更新: 2025-11-28

