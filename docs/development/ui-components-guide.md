# UI 组件使用指南

## 概述

BookImporter 提供了一套完整的 UI 组件库（`pkg/ui/`），用于创建美观、一致的终端用户界面。

## 组件库架构

```
pkg/ui/
├── styles.go       # 颜色、文本样式、状态指示符
├── components.go   # 可复用的 UI 组件
├── progress.go     # 进度跟踪器和 Spinner
├── table.go        # 表格生成器
└── terminal.go     # 终端能力检测
```

## 样式系统 (styles.go)

### 颜色主题

```go
import "github.com/jianyun8023/bookimporter/pkg/ui"

// 6种主题颜色
ui.RenderSuccess("操作成功")    // 绿色
ui.RenderError("操作失败")      // 红色
ui.RenderWarning("警告信息")    // 黄色
ui.RenderInfo("提示信息")       // 青色
ui.RenderPrimary("主要内容")    // 蓝色
ui.RenderSecondary("次要信息")  // 灰色
```

### 文本样式

```go
// 基础样式
ui.StyleBold("粗体文本")
ui.StyleItalic("斜体文本")
ui.StyleUnderline("下划线文本")
ui.StyleStrikethrough("删除线文本")

// 组合样式
ui.StyleSuccess("成功") + " " + ui.StyleBold("粗体")
```

### 状态指示符

```go
// Unicode 图标（自动降级为 ASCII）
ui.IconSuccess    // ✓ (OK)
ui.IconError      // ✗ (X)
ui.IconInfo       // → (>)
ui.IconWarning    // ⚠ (!)
ui.IconSkipped    // ⋯ (-)
```

## 可复用组件 (components.go)

### 统计表格

```go
stats := map[string]interface{}{
    "总计": 100,
    "成功": 85,
    "失败": 15,
}
fmt.Println(ui.RenderStatsSummary(stats))
```

### 消息框

```go
// 错误消息框
ui.RenderMessageBox("error", "操作失败", "详细错误信息")

// 警告消息框
ui.RenderMessageBox("warning", "警告", "请注意...")

// 信息消息框
ui.RenderMessageBox("info", "提示", "这是一个提示信息")
```

### 文件操作展示

```go
// 显示文件操作（old → new）
ui.FormatFileOperation("/path/to/file.txt", "book-1.txt")

// 重命名预览
ui.FormatRenamePreview("/old/path/file.txt", "/new/path/book-1.txt")
```

### 命令头部和分隔线

```go
// 命令头部
ui.RenderHeader("命令标题", "命令描述说明")

// 分隔线
ui.RenderSeparator()
```

## 进度跟踪 (progress.go)

### 标准进度条

```go
// 创建进度跟踪器
progress := ui.NewProgressTracker(total)

// 更新进度
progress.Increment()
fmt.Print("\r" + progress.Render())

// 完成
progress.Finish()
```

**输出**:
```
[████████████░░░░] 75% (15/20)
已用 3s | 预计剩余 1s
```

### 紧凑模式进度条

```go
// 创建紧凑模式进度跟踪器
progress := ui.NewCompactProgressTracker(total)

// 更新进度和统计
progress.IncrementSuccess()  // 成功 +1
progress.IncrementFailure()  // 失败 +1
progress.IncrementSkipped()  // 跳过 +1

fmt.Print("\r" + progress.RenderCompact("current-file.txt"))
```

**输出**:
```
[████████░░] 80% (8/10) ✓7 ✗1 current-file.txt
```

### 带统计的进度展示

```go
progress := ui.NewProgressTracker(total)
// ... 处理文件 ...
progress.IncrementSuccess()

// 显示进度条 + 统计
fmt.Println(progress.RenderWithStats())
```

**输出**:
```
[████████████████░░░░] 80% (8/10)
✓ 7  ✗ 1  ⋯ 2
```

### Spinner 动画

```go
// 基础 Spinner
spinner := ui.NewSpinner()
fmt.Print("\r" + spinner.Tick() + " 处理中...")

// 带时间的 Spinner
spinner := ui.NewSpinnerWithTime("处理中...")
fmt.Print("\r" + spinner.Render())  // ⠋ 处理中... (3s)

// 多任务 Spinner
tasks := []string{"任务 1", "任务 2", "任务 3"}
multiSpinner := ui.NewMultiSpinner(tasks)
multiSpinner.CompleteTask(0)     // 任务 1 完成
multiSpinner.SetActiveTask(1)    // 激活任务 2
fmt.Println(multiSpinner.Render())
```

**输出**:
```
✓ 任务 1
⠋ 任务 2
  任务 3
```

## 表格生成器 (table.go)

### 快速创建表格

```go
headers := []string{"列1", "列2", "列3"}
rows := [][]string{
    {"值1", "值2", "值3"},
    {"值4", "值5", "值6"},
}

table := ui.QuickTable(headers, rows)
fmt.Println(table)
```

### 自定义配置

```go
config := ui.NewTableConfig()
config.Headers = []string{"状态", "数量", "百分比"}
config.Rows = [][]string{
    {"✓ 成功", "85", "85.0%"},
    {"✗ 失败", "15", "15.0%"},
}
config.BorderStyle = "rounded"   // normal, rounded, double, thick, none
config.AlignRight = []int{1, 2}  // 列1和列2右对齐
config.CompactMode = true        // 紧凑模式

table := ui.RenderTable(config)
fmt.Println(table)
```

**输出**:
```
╭────────────┬───────┬──────────╮
│   状态     │ 数量  │ 百分比   │
├────────────┼───────┼──────────┤
│ ✓ 成功     │   85  │  85.0%   │
│ ✗ 失败     │   15  │  15.0%   │
╰────────────┴───────┴──────────╯
```

### 边框样式

```go
// 4种边框样式
config.BorderStyle = "normal"   // ┌─┐
config.BorderStyle = "rounded"  // ╭─╮
config.BorderStyle = "double"   // ╔═╗
config.BorderStyle = "thick"    // ┏━┓
config.BorderStyle = "none"     // 无边框
```

## 终端能力检测 (terminal.go)

```go
import "github.com/jianyun8023/bookimporter/pkg/ui"

// 检测终端能力
caps := ui.DetectTerminalCapabilities()

if caps.SupportsColor {
    // 使用彩色输出
} else {
    // 使用纯文本
}

if caps.SupportsUnicode {
    fmt.Println(ui.IconSuccess)  // ✓
} else {
    fmt.Println("[OK]")           // ASCII 备选
}

// 自动降级示例
icon := ui.GetIcon("success")  // 自动选择 ✓ 或 [OK]
```

### 环境变量支持

```go
// 检测 NO_COLOR 环境变量
if ui.ShouldDisableColor() {
    // 禁用彩色输出
}

// 在代码中使用
text := ui.StyleSuccess("成功")  // 自动根据环境决定是否着色
```

## 最佳实践

### 1. 批量操作进度展示

```go
total := len(files)
progress := ui.NewCompactProgressTracker(total)

for _, file := range files {
    // 处理文件
    result := processFile(file)
    
    if result.Success {
        progress.IncrementSuccess()
    } else if result.Skipped {
        progress.IncrementSkipped()
    } else {
        progress.IncrementFailure()
    }
    
    // 显示进度（单行刷新）
    fmt.Print("\r" + progress.RenderCompact(file.Name))
}

// 清除进度行
fmt.Print("\r" + strings.Repeat(" ", 120) + "\r")

// 显示最终统计
fmt.Println(progress.RenderWithStats())
```

### 2. 命令输出结构

```go
// 1. 显示命令头部
fmt.Println(ui.RenderHeader("命令标题", "命令描述"))

// 2. 显示找到的文件数
fmt.Println(ui.RenderInfo(fmt.Sprintf("找到 %d 个文件", total)))

// 3. 处理文件（显示进度）
// ... 处理逻辑 ...

// 4. 显示分隔线
fmt.Println(ui.RenderSeparator())

// 5. 显示统计表格
stats := map[string]interface{}{
    "成功": successCount,
    "失败": failureCount,
}
fmt.Println(ui.RenderStatsSummary(stats))

// 6. 显示完成消息
if failureCount > 0 {
    fmt.Println(ui.RenderError(fmt.Sprintf("有 %d 个文件处理失败", failureCount)))
} else {
    fmt.Println(ui.RenderSuccess("所有文件处理完成"))
}
```

### 3. 预览模式

```go
if doTry {
    // 使用表格展示预览
    headers := []string{"#", "原文件名", "→", "新文件名"}
    rows := [][]string{}
    
    for i, item := range previewItems {
        rows = append(rows, []string{
            fmt.Sprintf("%d", i+1),
            item.OldName,
            "→",
            item.NewName,
        })
    }
    
    config := ui.NewTableConfig()
    config.Headers = headers
    config.Rows = rows
    config.BorderStyle = "rounded"
    
    fmt.Println(ui.RenderInfo("📋 重命名预览"))
    fmt.Println()
    fmt.Println(ui.RenderTable(config))
    fmt.Println()
    fmt.Println(ui.RenderInfo(fmt.Sprintf("📝 [试运行] 将重命名 %d 个文件", len(rows))))
}
```

### 4. 错误处理

```go
if err != nil {
    if skipOnError {
        // 显示警告并继续
        fmt.Println(ui.RenderWarning(fmt.Sprintf("跳过: %v", err)))
        stats.Skipped++
    } else {
        // 显示错误并停止
        fmt.Println(ui.RenderError(fmt.Sprintf("处理失败: %v", err)))
        fmt.Println()
        fmt.Println(ui.RenderInfo("提示: 使用 -i 参数可以忽略错误退出码"))
        os.Exit(1)
    }
}
```

## 性能注意事项

1. **进度条更新**: 使用 `\r` 单行刷新，避免大量输出
2. **表格渲染**: 纯字符串操作，零运行时开销
3. **颜色检测**: 一次检测，全局缓存
4. **批量操作**: 限制更新频率（如每处理10个文件更新一次）

## 终端兼容性

- ✅ macOS Terminal
- ✅ iTerm2
- ✅ Windows Terminal
- ✅ Alacritty
- ✅ Tmux/Screen
- ✅ CI/CD 环境（自动降级）

## 示例代码

完整的使用示例可以参考：
- `cmd/check.go` - check 命令的 UI 实现
- `cmd/clname.go` - clname 命令的 UI 实现
- `cmd/rename.go` - rename 命令的 UI 实现

