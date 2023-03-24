package epub

// Metadata Define EPUB metadata
type Metadata struct {
	// 书籍标题
	Title string
	// 作者
	Author string
	// 书籍描述
	Description string
	// 出版社
	Publisher string
	// 出版日期
	Date string
	// 书籍标识
	Identifier []Identifier
	// 语言
	Language string
	// ISBN
	Isbn string
}

// Identifier 定义epub元数据标识
type Identifier struct {
	// 标识ID
	ID string
	// 标识类型 eg: ISBN\ASIN\MOBI-ASIN
	Scheme string
	// 标识值
	Value string
}
