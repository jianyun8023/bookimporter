package book

// Book 代表一本电子书实体
type Book struct {
	Path     string   // 文件绝对路径
	Title    string   // 书籍标题
	Author   []string // 作者列表
	Format   string   // 文件格式 (如 epub)
	IsValid  bool     // 文件是否完整有效
	FileSize int64    // 文件大小
}

// NewBook 创建一个新的 Book 实例
func NewBook(path, title string, authors []string) *Book {
	return &Book{
		Path:    path,
		Title:   title,
		Author:  authors,
		Format:  "epub", // 目前默认为 epub，后续可扩展
		IsValid: true,
	}
}
