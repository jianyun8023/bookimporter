package epub

import (
	"path/filepath"
	"testing"
)

func TestFindIsbn(t *testing.T) {
	// 查找目录下的epub文件

	matches, err := filepath.Glob(filepath.Join("/Users/jianyun/Downloads/", "*.epub"))
	if err != nil {
		t.Fatal(err)
	}

	for _, match := range matches {

		println(match)
		isbn, err := FindIsbn(match)
		if err != nil {
			t.Fatal(err)
		}
		if isbn == "" {
			t.Log("isbn is empty")
		}
		t.Log(isbn)
	}
}
