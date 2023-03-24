package epub

import "testing"

// 测试ReadMetadata
func TestReadMetadata(t *testing.T) {
	m, err := ReadMetadata("/Users/jianyun/Downloads/Google系统架构解密 构建安全可靠的系统 2021 - [美]希瑟·阿德金斯 [美]贝齐·拜尔_224663.epub")
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("metadata is nil")
	}
	t.Log(m)
}
