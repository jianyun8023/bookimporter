# v2025-11-28 命令优化更新

> 更新日期: 2025-11-28  
> 相关版本: Unreleased  
> 影响命令: clname, rename

## 概述

本次更新主要优化了 `clname` 和 `rename` 两个命令的用户体验，包括功能增强、帮助信息优化和错误处理改进。

## Clname 命令优化

### 新增功能

#### 1. 递归搜索 (`-r/--recursive`)
- 支持递归搜索子目录中的所有 EPUB 文件
- 解决了之前只能扫描当前目录的问题

**使用示例：**
```bash
# 递归搜索所有子目录
bookimporter clname -p ./books -r

# 结合其他参数使用
bookimporter clname -p ./books -r -i
```

#### 2. 忽略错误退出码 (`-i/--ignore-errors`)
- 允许即使有失败也返回退出码 0
- 适用于 CI/CD 脚本或批处理场景

**使用示例：**
```bash
# 忽略错误，总是返回 0
bookimporter clname -p ./books -r -i
```

### 行为变更

#### 默认错误处理改变
- **旧行为**: 遇到错误会停止处理（除非使用 `-j` 参数）
- **新行为**: 默认跳过错误继续处理其他文件（更实用）
- **退出码**: 
  - 有失败时返回 `1`（除非使用 `-i` 参数）
  - 所有成功或使用 `-i` 时返回 `0`

#### 移除参数
- **移除 `-j/--skip` 参数**（破坏性变更）
- 原因: 其功能已成为默认行为

### 帮助信息优化

1. **详细的功能描述**
   - 清晰说明命令的核心功能
   - 列举所有主要特性

2. **5 个实用示例**
   ```bash
   # 1. 单个文件
   bookimporter clname -p book.epub
   
   # 2. 当前目录
   bookimporter clname -p ./
   
   # 3. 递归搜索
   bookimporter clname -p ./books -r
   
   # 4. 预览模式
   bookimporter clname -p ./books -r -t
   
   # 5. 处理损坏文件
   bookimporter clname -p ./books -r --move-corrupted-to ./corrupted
   ```

3. **参数分组说明**
   - 基础参数
   - 搜索选项
   - 错误处理
   - 其他选项

### 参数验证增强

1. **互斥参数检查**
   - `--move-corrupted-to` 与 `--delete-corrupted` 不能同时使用

2. **参数依赖警告**
   - `--force-delete` 需要配合 `--delete-corrupted` 使用

3. **路径安全检查**
   - 防止将文件移动到源目录的子目录

### 错误处理改进

1. 不再使用 `panic`，优雅处理错误
2. 显示详细的文件路径和错误信息
3. 统计信息区分"跳过"和"失败"

## Rename 命令优化

### 帮助信息全面升级

#### 1. 命令描述优化

**优化前:**
```
Use:   "rename"
Short: "Rename or move files according to a template"
```

**优化后:**
```
Use:   "rename <目录路径>"
Short: "按自定义模板批量重命名或移动文件"
Long:  详细的中文说明，包括：
  • 核心功能列表
  • 序号占位符说明
  • 重要说明事项
  • 使用建议
```

#### 2. 新增 6 个实用示例

```bash
# 1. 基础用法
bookimporter rename . -f txt -t "book-@n"

# 2. 递归搜索
bookimporter rename /path/to/books -f epub -t "novel-@n" -r

# 3. 多格式处理
bookimporter rename . -f epub -f pdf -f mobi -t "ebook-@n"

# 4. 移动文件
bookimporter rename /source -f jpg -t "photo-@n" -o /photos

# 5. 自定义起始序号
bookimporter rename . -f txt -t "doc-@n" --start-num 100

# 6. 预览模式
bookimporter rename . -f pdf -t "file-@n" --do-try
```

#### 3. 参数说明优化

**改进点:**
- ✅ 全部改为中文说明
- ✅ 添加更详细的功能描述
- ✅ 在说明中给出具体示例
- ✅ 按重要性重新排序

**优化对比:**

| 参数 | 优化前 | 优化后 |
|------|--------|--------|
| -f | File format to match | 指定要处理的文件格式（如 'txt', 'epub'），可多次使用 |
| -t | Template for new filename | 文件名模板，@n 为序号占位符（如 'book-@n' → book-1.epub） |
| -r | Recursively search | 递归搜索子目录中的所有匹配文件 |

### 错误提示优化

#### 1. 模板格式错误

**优化前:**
```
模板 [book] 中不存在占位符 @n
```

**优化后:**
```
✗ 错误: 模板 'book' 中缺少序号占位符 @n

→ 模板必须包含 @n 作为序号占位符，例如:
  ✓ 正确: -t "book-@n"    → book-1.epub, book-2.epub
  ✓ 正确: -t "file_@n"    → file_1.txt, file_2.txt
  ✗ 错误: -t "book"       (没有 @n)
```

#### 2. 目录不存在

**新增友好提示:**
```
✗ 错误: 目录不存在: /path/to/dir
```

### 参数验证增强

1. **模板验证**
   - 检查是否包含 `@n` 占位符
   - 提供详细的正确/错误示例

2. **目录验证**
   - 检查源目录是否存在
   - 显示友好的错误提示

3. **输出目录处理**
   - 自动创建输出目录（如果不存在）
   - 避免用户手动创建的麻烦

## 用户体验提升

### 整体改进

| 项目 | 优化前 | 优化后 |
|------|--------|--------|
| 界面语言 | 英文为主 | 全中文 |
| 使用示例 | 缺少或简单 | 详细丰富（5-6个） |
| 参数说明 | 简单 | 详细，带示例 |
| 错误提示 | 简单 | 友好，带解决方案 |
| 参数验证 | 基础 | 完善 |

### 设计原则

1. **用户友好**: 全中文界面，清晰的说明
2. **示例驱动**: 提供丰富的实用示例
3. **渐进式学习**: 从简单到复杂的示例排列
4. **错误友好**: 不仅指出错误，还提供解决方案
5. **预防式设计**: 通过验证避免常见错误

## 迁移指南

### 从旧版本升级

#### Clname 命令

如果你之前使用了 `-j/--skip` 参数：

**旧用法（不再可用）:**
```bash
./bookimporter clname -p ./books/ -j
```

**新用法（-j 已是默认行为）:**
```bash
# 默认就会跳过错误
./bookimporter clname -p ./books/

# 如果需要忽略退出码
./bookimporter clname -p ./books/ -i
```

#### 退出码变更

**新行为:**
- 所有文件处理成功：退出码 `0`
- 有文件处理失败（默认）：退出码 `1`
- 有文件处理失败但使用 `-i` 参数：退出码 `0`

**脚本适配建议:**
```bash
# 如果你的脚本依赖退出码
if ./bookimporter clname -p ./books/ -r; then
    echo "处理成功"
else
    echo "有失败（但已处理完所有文件）"
fi

# 如果你不关心是否有失败
./bookimporter clname -p ./books/ -r -i
```

## 技术细节

### 代码改动

#### Clname 命令
- 文件: `cmd/clname.go`
- 主要修改:
  1. 添加 `Recursive` 和 `IgnoreErrors` 字段
  2. 使用 `filepath.Walk` 实现递归搜索
  3. 优化错误处理逻辑
  4. 增强参数验证
  5. 完善帮助信息

#### Rename 命令
- 文件: `cmd/rename.go`
- 主要修改:
  1. 优化 `Use`, `Short`, `Long`, `Example` 字段
  2. 改进参数说明（全中文）
  3. 增强 `validateConfig()` 函数
  4. 优化错误提示信息

### 测试验证

所有改动已通过以下测试：
- ✅ 帮助信息显示测试
- ✅ 参数验证测试
- ✅ 错误提示测试
- ✅ 功能执行测试

## 相关文档

- `Changelog.md` - 完整的变更记录
- `docs/USER_GUIDE.md` - 用户指南
- `docs/guides/help-improvement.md` - 帮助信息改进指南
- `README.md` - 项目说明

## 后续计划

### 短期改进
- [ ] 为 `check` 命令应用相同的优化模式
- [ ] 为 `version` 命令添加详细信息
- [ ] 添加更多边界情况的测试

### 长期规划
- [ ] 考虑添加交互式模式
- [ ] 支持更多的占位符（如日期、时间等）
- [ ] 添加撤销功能
- [ ] 支持批量操作的进度保存和恢复

---

最后更新: 2025-11-28

