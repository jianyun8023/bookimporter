package epub

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
)

func TestService_Validate_FileNotExists(t *testing.T) {
	s := NewService()
	err := s.Validate("/path/to/nonexistent/file.epub")
	if err == nil {
		t.Error("期望返回错误，但得到 nil")
	}
	if !isEpubError(err) {
		t.Error("期望返回 EpubError")
	}
	epubErr := err.(*EpubError)
	if epubErr.Type != ErrorTypeCorrupted {
		t.Errorf("期望错误类型为 ErrorTypeCorrupted，得到 %v", epubErr.Type)
	}
}

func TestService_Validate_CorruptedZip(t *testing.T) {
	s := NewService()
	// 创建一个损坏的 zip 文件
	tmpDir := t.TempDir()
	corruptedFile := filepath.Join(tmpDir, "corrupted.epub")

	// 写入无效的 ZIP 数据
	err := os.WriteFile(corruptedFile, []byte("This is not a valid ZIP file"), 0644)
	if err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}

	err = s.Validate(corruptedFile)
	if err == nil {
		t.Error("期望返回错误，但得到 nil")
	}
	if !isEpubError(err) {
		t.Error("期望返回 EpubError")
	}
	epubErr := err.(*EpubError)
	if epubErr.Type != ErrorTypeCorrupted {
		t.Errorf("期望错误类型为 ErrorTypeCorrupted，得到 %v", epubErr.Type)
	}
}

func TestService_Validate_MissingRequiredFiles(t *testing.T) {
	s := NewService()
	// 创建一个缺少必需文件的 EPUB（有效的 ZIP，但缺少 mimetype）
	tmpDir := t.TempDir()
	epubFile := filepath.Join(tmpDir, "missing_files.epub")

	// 创建一个只有一个普通文件的 ZIP
	f, err := os.Create(epubFile)
	if err != nil {
		t.Fatalf("无法创建测试文件: %v", err)
	}
	defer f.Close()

	w := zip.NewWriter(f)
	// 只添加一个普通文件，不添加必需的 mimetype 和 container.xml
	fileWriter, err := w.Create("test.txt")
	if err != nil {
		t.Fatalf("无法创建 ZIP 条目: %v", err)
	}
	_, err = fileWriter.Write([]byte("test content"))
	if err != nil {
		t.Fatalf("无法写入 ZIP 条目: %v", err)
	}
	w.Close()

	err = s.Validate(epubFile)
	if err == nil {
		t.Error("期望返回错误，但得到 nil")
	}
	if !isEpubError(err) {
		t.Error("期望返回 EpubError")
	}
	epubErr := err.(*EpubError)
	if epubErr.Type != ErrorTypeMissing {
		t.Errorf("期望错误类型为 ErrorTypeMissing，得到 %v", epubErr.Type)
	}
}

func TestEpubError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *EpubError
		expected string
	}{
		{
			name: "With detail",
			err: &EpubError{
				Type:    ErrorTypeCorrupted,
				Message: "文件损坏",
				Detail:  "无法读取",
			},
			expected: "文件损坏: 无法读取",
		},
		{
			name: "Without detail",
			err: &EpubError{
				Type:    ErrorTypeMissing,
				Message: "文件缺失",
			},
			expected: "文件缺失",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("期望 %q，得到 %q", tt.expected, result)
			}
		})
	}
}

// Helper to check error type locally since IsEpubError was in pkg/util
func isEpubError(err error) bool {
	_, ok := err.(*EpubError)
	return ok
}
