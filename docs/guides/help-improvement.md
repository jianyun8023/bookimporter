# 帮助信息改进指南

## 概述

本文档说明了如何为 BookImporter 命令编写高质量的帮助信息。

## 帮助信息结构

### 1. 命令描述 (Short & Long)

#### Short - 简短描述
- 一行简明概括，显示在命令列表中
- 25-50 字符为宜

```go
Short: "清理书籍标题中的无用描述"
```

#### Long - 详细描述
- 多行详细说明
- 包含功能介绍、实现说明、支持特性

```go
Long: `清理 EPUB 书籍标题中的无用描述符和标记

自动移除书籍标题中的各种括号标记，如：（）【】()[]
使用 Calibre 的 ebook-meta 工具修改元数据。

支持：
  • 单个文件或批量目录处理
  • 递归搜索子目录
  • 预览模式（不实际修改）
  • 自动处理损坏的 EPUB 文件`,
```

### 2. 使用示例 (Example)

提供 3-5 个常见使用场景的实用示例：

```go
Example: `  # 清理单个文件
  bookimporter clname -p book.epub

  # 批量清理当前目录（不包含子目录）
  bookimporter clname -p /path/to/books/

  # 递归清理所有子目录
  bookimporter clname -p /path/to/books/ -r

  # 预览模式（不实际修改）
  bookimporter clname -p /path/to/books/ -r -t

  # 跳过错误，移动损坏文件
  bookimporter clname -p /path/to/books/ -r -j --move-corrupted-to /corrupted/`,
```

### 3. 参数说明 (Flags)

#### 基本要求
- 清晰、准确的描述
- 说明参数的作用和影响
- 注明互斥关系和依赖关系

#### 示例

```go
// 基础参数
clnameCmd.Flags().StringVarP(&config.Path, "path", "p", "./", 
    "指定要处理的 EPUB 文件或目录路径")
clnameCmd.Flags().BoolVarP(&config.Recursive, "recursive", "r", false, 
    "递归搜索子目录中的所有 EPUB 文件")

// 运行模式
clnameCmd.Flags().BoolVarP(&config.DoTry, "dotry", "t", false, 
    "预览模式，显示将要进行的修改但不实际执行")
clnameCmd.Flags().BoolVarP(&config.Skip, "skip", "j", false, 
    "遇到无法解析或损坏的文件时跳过，继续处理其他文件")

// 损坏文件处理（互斥选项）
clnameCmd.Flags().StringVar(&config.MoveCorruptedTo, "move-corrupted-to", "", 
    "将损坏的 EPUB 文件移动到指定目录（与 --delete-corrupted 互斥）")
clnameCmd.Flags().BoolVar(&config.DeleteCorrupted, "delete-corrupted", false, 
    "直接删除损坏的 EPUB 文件（与 --move-corrupted-to 互斥）")
clnameCmd.Flags().BoolVar(&config.ForceDelete, "force-delete", false, 
    "删除损坏文件时不需要用户确认（需配合 --delete-corrupted 使用）")

// 调试选项
clnameCmd.Flags().BoolVarP(&config.Debug, "debug", "d", false, 
    "启用调试模式，显示详细的执行信息")
```

## 参数验证

### 互斥参数检查

```go
func ValidateConfig(c *ClnameConfig) error {
    // 互斥参数检查
    if c.MoveCorruptedTo != "" && c.DeleteCorrupted {
        fmt.Println(ui.RenderError("--move-corrupted-to 和 --delete-corrupted 参数不能同时使用"))
        fmt.Println(ui.RenderInfo("提示: 请选择移动或删除损坏文件，不能同时进行"))
        os.Exit(1)
    }
    
    return nil
}
```

### 参数依赖警告

```go
func ValidateConfig(c *ClnameConfig) error {
    // 参数依赖警告
    if c.ForceDelete && !c.DeleteCorrupted {
        fmt.Println(ui.RenderWarning("--force-delete 参数需要配合 --delete-corrupted 使用"))
        fmt.Println(ui.RenderInfo("提示: --force-delete 用于在删除损坏文件时跳过确认步骤"))
    }
    
    return nil
}
```

### 路径安全检查

```go
func ValidateConfig(c *ClnameConfig) error {
    // 路径安全检查
    if c.MoveCorruptedTo != "" {
        absPath, _ := filepath.Abs(c.Path)
        absDst, _ := filepath.Abs(c.MoveCorruptedTo)
        
        // 确保目标目录不是源目录的子目录
        if strings.HasPrefix(absDst, absPath) {
            fmt.Println(ui.RenderError("错误: 目标目录不能是源目录的子目录"))
            fmt.Println(ui.RenderInfo("提示: 这样会导致循环处理"))
            os.Exit(1)
        }
    }
    
    return nil
}
```

## 帮助信息最佳实践

### 1. 描述要点

- ✅ **清晰明确**: 避免模糊的表述
- ✅ **用户视角**: 关注用户需要知道什么
- ✅ **实际用途**: 说明实际应用场景
- ❌ **避免术语**: 不要使用技术黑话

**示例对比**:

| 不好 | 好 |
|------|-----|
| "递归搜索子目录" | "递归搜索子目录中的所有 EPUB 文件" |
| "尝试运行" | "预览模式，显示将要进行的修改但不实际执行" |
| "跳过" | "遇到无法解析或损坏的文件时跳过，继续处理其他文件" |

### 2. 示例要点

- ✅ **实用性**: 选择最常见的使用场景
- ✅ **渐进性**: 从简单到复杂
- ✅ **注释清晰**: 每个示例都有清晰的注释
- ✅ **可复制**: 用户可以直接复制使用

### 3. 参数分组

按功能逻辑分组参数：

```
基础参数:
  -p, --path          文件或目录路径
  -r, --recursive     递归搜索

运行模式:
  -t, --dotry         预览模式
  -j, --skip          跳过错误

损坏文件处理 (互斥):
  --move-corrupted-to 移动到指定目录
  --delete-corrupted  直接删除
  --force-delete      强制删除（无需确认）

调试选项:
  -d, --debug         调试模式
```

### 4. 错误提示

良好的错误提示应该：

- ✅ 清楚说明错误原因
- ✅ 提供解决建议
- ✅ 使用友好的语气

**示例**:

```go
// 不好
fmt.Println("参数错误")

// 好
fmt.Println(ui.RenderError("--move-corrupted-to 和 --delete-corrupted 参数不能同时使用"))
fmt.Println(ui.RenderInfo("提示: 请选择移动或删除损坏文件，不能同时进行"))
```

## 完整示例

查看 `cmd/clname.go` 获取完整的帮助信息实现示例。

## 测试帮助信息

```bash
# 查看命令帮助
./bookimporter clname --help

# 测试互斥参数
./bookimporter clname --move-corrupted-to ./corrupted --delete-corrupted

# 测试依赖参数
./bookimporter clname --force-delete

# 测试完整功能
./bookimporter clname -p ./test -r -j --move-corrupted-to ./corrupted
```

## 文档同步

帮助信息修改后，需要同步更新：

1. **README.md** - 基础使用说明
2. **USER_GUIDE.md** - 详细用户指南
3. **CLAUDE.md** - AI 开发指南
4. **Changelog.md** - 版本变更记录

