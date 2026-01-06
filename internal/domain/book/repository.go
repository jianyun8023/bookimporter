package book

// Repository 定义书籍持久化存储接口
type Repository interface {
	// FindAll 查找指定目录下的所有书籍
	// recursive: 是否递归查找子目录
	FindAll(dir string, recursive bool) ([]*Book, error)

	// Rename 重命名书籍文件
	// oldPath: 原文件路径
	// newPath: 新文件路径
	// 返回移动后的实际路径（处理冲突后）和错误
	Rename(oldPath, newPath string) (string, error)

	// Move 将书籍移动到指定目录
	// srcPath: 源文件路径
	// destDir: 目标目录
	// 返回移动后的实际路径（处理冲突后）和错误
	Move(srcPath, destDir string) (string, error)

	// Delete 删除书籍文件
	// force: 是否强制删除（不确认）
	Delete(path string) error

	// UpdateMetadata 更新书籍元数据（如标题）
	UpdateMetadata(path string, title string) error

	// Exists 检查路径是否存在
	Exists(path string) bool

	// IsDir 检查是否为目录
	IsDir(path string) bool

	// IsFile 检查是否为文件
	IsFile(path string) bool

	// MkdirAll 创建目录
	MkdirAll(path string) error
}
