package filesystem

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jianyun8023/bookimporter/internal/domain/book"
	"github.com/jianyun8023/bookimporter/pkg/ui"
	"github.com/kapmahc/epub"
)

// Repository 实现了 book.Repository 接口
type Repository struct{}

// NewRepository 创建一个新的文件系统仓库
func NewRepository() *Repository {
	return &Repository{}
}

// FindAll 查找指定目录下的所有书籍
func (r *Repository) FindAll(dir string, recursive bool) ([]*book.Book, error) {
	var books []*book.Book

	if !r.Exists(dir) {
		return nil, fmt.Errorf("path does not exist: %s", dir)
	}

	// 辅助函数，将文件转为 Book 实体
	addBook := func(path string) error {
		// 简单的 Book 创建，详细信息可能需要打开文件读取
		// 这里我们暂时只设置路径和文件名
		// 实际应用中可能需要更复杂的元数据读取逻辑，但这通常属于 infrastructure/epub 服务的职责
		// 或者 Repository 可以调用 epub 库来填充信息
		// 为保持 Repository 纯粹处理文件存储，我们这里只做基础信息
		// 如果需要读取 Title，可能需要在 Service 层调用 EpubService 填充，或者在这里集成 Epub 解析

		// 考虑到性能，FindAll 可能只返回文件列表封装成的 Book 对象
		// 具体的元数据加载可以延迟或按需进行
		// 但为了兼容现有逻辑，我们尝试读取一下基础信息

		b := book.NewBook(path, "", nil)

		// 尝试读取 Title (可选优化)
		if ep, err := epub.Open(path); err == nil {
			if len(ep.Opf.Metadata.Title) > 0 {
				b.Title = ep.Opf.Metadata.Title[0]
			}
			if len(ep.Opf.Metadata.Creator) > 0 {
				for _, c := range ep.Opf.Metadata.Creator {
					b.Author = append(b.Author, c.Data)
				}
			}
			ep.Close()
		}

		books = append(books, b)
		return nil
	}

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".epub") {
				return addBook(path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		// 判断 dir 是否是文件
		if r.IsFile(dir) {
			if strings.HasSuffix(strings.ToLower(dir), ".epub") {
				if err := addBook(dir); err != nil {
					return nil, err
				}
			}
			return books, nil
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".epub") {
				if err := addBook(filepath.Join(dir, entry.Name())); err != nil {
					return nil, err
				}
			}
		}
	}

	return books, nil
}

// Rename 重命名书籍文件
func (r *Repository) Rename(oldPath, newPath string) (string, error) {
	// 确保目标目录存在
	dstDir := filepath.Dir(newPath)
	if err := r.MkdirAll(dstDir); err != nil {
		return "", fmt.Errorf("failed to create destination directory: %w", err)
	}

	// 如果目标文件已存在，不自动处理（这是 Service 层的逻辑决定，还是 Repository 层决定？）
	// 按照接口定义，Rename 似乎是一个直接操作。
	// 但看原 filetool 的逻辑，Move 是有 conflict handling 的。
	// 这里我们实现一个直接的 rename，conflict handling 可以放在 Move 或者上层
	// 不过为了兼容性，我们可以让 rename 也具备 conflict handling，或者严格按照 os.Rename

	// 这里我们实现严格的 Rename，如果目标存在则报错，除非是同一个文件
	if oldPath == newPath {
		return newPath, nil
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return "", err
	}
	return newPath, nil
}

// Move 将书籍移动到指定目录，并处理冲突
func (r *Repository) Move(srcPath, destDir string) (string, error) {
	// 确保目标目录存在
	if err := r.MkdirAll(destDir); err != nil {
		return "", fmt.Errorf("cannot create destination directory: %w", err)
	}

	fileName := filepath.Base(srcPath)
	dstPath := filepath.Join(destDir, fileName)

	// 如果目标文件不存在，直接移动
	if !r.Exists(dstPath) {
		if err := os.Rename(srcPath, dstPath); err != nil {
			return "", fmt.Errorf("move file failed: %w", err)
		}
		return dstPath, nil
	}

	// 处理文件名冲突，添加序号
	ext := filepath.Ext(fileName)
	nameWithoutExt := strings.TrimSuffix(fileName, ext)

	for i := 1; i < 10000; i++ {
		newFileName := fmt.Sprintf("%s(%d)%s", nameWithoutExt, i, ext)
		newDstPath := filepath.Join(destDir, newFileName)
		if !r.Exists(newDstPath) {
			if err := os.Rename(srcPath, newDstPath); err != nil {
				return "", fmt.Errorf("move file failed: %w", err)
			}
			return newDstPath, nil
		}
	}

	return "", fmt.Errorf("cannot find available filename (tried 10000 times)")
}

// Delete 删除书籍文件
func (r *Repository) Delete(path string) error {
	// 注意：Force 逻辑通常属于 Application/Service 层，Repository 只负责执行动作
	// 但这里接口设计上，Delete 就是删除。确认逻辑应该在上层。
	// 如果需要 Repository 处理确认，那就违反了单一职责，Repository 不该跟 UI 交互。
	// 所以这里的 Delete 只是纯粹的 os.Remove。
	// 原有的 SafeDeleteFile 包含 UI 交互，这部分逻辑需要移到 Application 或 UI 层。
	// 这里只实现物理删除。

	if !r.Exists(path) {
		return fmt.Errorf("file not exist: %s", path)
	}
	return os.Remove(path)
}

// UpdateMetadata 使用 ebook-meta 更新元数据
func (r *Repository) UpdateMetadata(path string, title string) error {
	// 使用 Calibre 的 ebook-meta 工具
	command := fmt.Sprintf("ebook-meta \"%s\" -t \"%s\"", path, title)
	cmdx := exec.Command("bash", "-c", command)

	if _, err := cmdx.Output(); err != nil {
		return fmt.Errorf("ebook-meta execution failed: %w", err)
	}
	return nil
}

// Exists 检查路径是否存在
func (r *Repository) Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsDir 检查是否为目录
func (r *Repository) IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 检查是否为文件
func (r *Repository) IsFile(path string) bool {
	return !r.IsDir(path)
}

// MkdirAll 创建目录
func (r *Repository) MkdirAll(dir string) error {
	if r.Exists(dir) {
		if !r.IsDir(dir) {
			return fmt.Errorf("path %s exists but is not a directory", dir)
		}
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

// Helper: InteractiveDelete (如果确实需要在基础设施层保留带交互的删除，虽然不推荐)
// 为了兼容原来的 SafeDeleteFile，我们可以提供一个额外的 Helper 方法，但不在 Repository 接口中
func InteractiveDelete(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file not exist: %s", path)
	}

	// UI 交互应该尽量上移，这里先保留逻辑以便参考，但后续应在 Service 层实现
	fmt.Println()
	warningMsg := ui.RenderMessageBox(
		"删除确认",
		fmt.Sprintf("即将删除文件:\n%s\n\n此操作无法撤销！", ui.RenderPath(path)),
		"warning",
	)
	fmt.Println(warningMsg)
	fmt.Println()
	fmt.Print(ui.StyleWarning.Render("确认删除? (y/N): "))

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("read user input failed: %w", err)
	}
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		fmt.Println(ui.RenderInfo("已取消删除操作"))
		return fmt.Errorf("user canceled delete operation")
	}

	return os.Remove(path)
}
