package util

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"github.com/kapmahc/epub"
)

// EpubError 表示 EPUB 文件检测错误
type EpubError struct {
	Type    ErrorType // 错误类型
	Message string    // 错误信息
	Detail  string    // 详细描述
}

// ErrorType 定义错误类型
type ErrorType int

const (
	ErrorTypeCorrupted ErrorType = iota // 文件损坏
	ErrorTypeMissing                    // 结构缺失
	ErrorTypeFormat                     // 格式错误
	ErrorTypeMetadata                   // 元数据缺失
)

// Error 实现 error 接口
func (e *EpubError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Detail)
	}
	return e.Message
}

// ValidateEpubFile 检测 EPUB 文件完整性
// 返回 nil 表示文件正常，返回 EpubError 表示检测到问题
func ValidateEpubFile(filePath string) error {
	// 1. 检查文件是否存在
	if !Exists(filePath) {
		return &EpubError{
			Type:    ErrorTypeCorrupted,
			Message: "文件不存在",
			Detail:  filePath,
		}
	}

	// 2. 检查 ZIP 文件完整性
	if err := checkZipIntegrity(filePath); err != nil {
		return err
	}

	// 3. 检查必需文件存在性
	if err := checkRequiredFiles(filePath); err != nil {
		return err
	}

	// 4. 检查元数据可解析性
	if err := checkMetadata(filePath); err != nil {
		return err
	}

	return nil
}

// checkZipIntegrity 检查 ZIP 文件完整性
func checkZipIntegrity(filePath string) error {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return &EpubError{
			Type:    ErrorTypeCorrupted,
			Message: "ZIP 文件损坏",
			Detail:  "无法打开或读取文件",
		}
	}
	defer r.Close()

	// 尝试读取所有文件条目以确保 ZIP 结构完整
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return &EpubError{
				Type:    ErrorTypeCorrupted,
				Message: "ZIP 文件损坏",
				Detail:  fmt.Sprintf("无法读取文件条目: %s", f.Name),
			}
		}
		// 读取一小部分数据以验证可读性
		buf := make([]byte, 1024)
		_, err = rc.Read(buf)
		if err != nil && err != io.EOF {
			rc.Close()
			return &EpubError{
				Type:    ErrorTypeCorrupted,
				Message: "ZIP 文件损坏",
				Detail:  fmt.Sprintf("文件条目数据损坏: %s", f.Name),
			}
		}
		rc.Close()
	}

	return nil
}

// checkRequiredFiles 检查 EPUB 必需文件
func checkRequiredFiles(filePath string) error {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return &EpubError{
			Type:    ErrorTypeCorrupted,
			Message: "无法打开文件",
			Detail:  err.Error(),
		}
	}
	defer r.Close()

	// EPUB 必需文件列表
	requiredFiles := map[string]bool{
		"mimetype":               false,
		"META-INF/container.xml": false,
	}

	// 检查文件是否存在
	for _, f := range r.File {
		name := strings.TrimPrefix(f.Name, "./")
		if _, required := requiredFiles[name]; required {
			requiredFiles[name] = true
		}
	}

	// 验证所有必需文件都存在
	for file, found := range requiredFiles {
		if !found {
			return &EpubError{
				Type:    ErrorTypeMissing,
				Message: "缺少必需文件",
				Detail:  file,
			}
		}
	}

	return nil
}

// checkMetadata 检查元数据可解析性
func checkMetadata(filePath string) error {
	book, err := epub.Open(filePath)
	if err != nil {
		return &EpubError{
			Type:    ErrorTypeFormat,
			Message: "无法解析 EPUB 元数据",
			Detail:  err.Error(),
		}
	}

	// 检查基本元数据
	if book == nil {
		return &EpubError{
			Type:    ErrorTypeFormat,
			Message: "EPUB 解析失败",
			Detail:  "无法读取书籍结构",
		}
	}

	if len(book.Opf.Metadata.Title) == 0 {
		return &EpubError{
			Type:    ErrorTypeMetadata,
			Message: "缺少书籍标题",
			Detail:  "元数据中未找到标题信息",
		}
	}

	return nil
}

// IsEpubError 判断是否为 EpubError
func IsEpubError(err error) bool {
	_, ok := err.(*EpubError)
	return ok
}

// GetErrorType 获取 EpubError 的错误类型
func GetErrorType(err error) ErrorType {
	if epubErr, ok := err.(*EpubError); ok {
		return epubErr.Type
	}
	return ErrorTypeCorrupted
}
