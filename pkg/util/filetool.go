package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// EnsureDir 确保目录存在，不存在则创建
func EnsureDir(dir string) error {
	if Exists(dir) {
		if !IsDir(dir) {
			return fmt.Errorf("路径 %s 存在但不是目录", dir)
		}
		return nil
	}
	return os.MkdirAll(dir, 0755)
}

// MoveFileWithConflictHandling 移动文件并处理重名冲突
// 如果目标文件已存在，会自动添加序号后缀，如 file(1).epub, file(2).epub
func MoveFileWithConflictHandling(srcPath, dstDir string) (string, error) {
	// 确保目标目录存在
	if err := EnsureDir(dstDir); err != nil {
		return "", fmt.Errorf("无法创建目标目录: %w", err)
	}

	fileName := filepath.Base(srcPath)
	dstPath := filepath.Join(dstDir, fileName)

	// 如果目标文件不存在，直接移动
	if !Exists(dstPath) {
		if err := os.Rename(srcPath, dstPath); err != nil {
			return "", fmt.Errorf("移动文件失败: %w", err)
		}
		return dstPath, nil
	}

	// 处理文件名冲突，添加序号
	ext := filepath.Ext(fileName)
	nameWithoutExt := strings.TrimSuffix(fileName, ext)

	for i := 1; i < 10000; i++ {
		newFileName := fmt.Sprintf("%s(%d)%s", nameWithoutExt, i, ext)
		newDstPath := filepath.Join(dstDir, newFileName)
		if !Exists(newDstPath) {
			if err := os.Rename(srcPath, newDstPath); err != nil {
				return "", fmt.Errorf("移动文件失败: %w", err)
			}
			return newDstPath, nil
		}
	}

	return "", fmt.Errorf("无法找到可用的文件名（尝试了 10000 次）")
}

// SafeDeleteFile 安全删除文件
// needConfirm 为 true 时会要求用户确认
func SafeDeleteFile(filePath string, needConfirm bool) error {
	if !Exists(filePath) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}

	if needConfirm {
		fmt.Printf("确认删除文件: %s? (y/N): ", filePath)
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取用户输入失败: %w", err)
		}
		input = strings.TrimSpace(strings.ToLower(input))
		if input != "y" && input != "yes" {
			return fmt.Errorf("用户取消删除操作")
		}
	}

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除文件失败: %w", err)
	}

	return nil
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("无法打开源文件: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("无法创建目标文件: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("复制文件失败: %w", err)
	}

	return dstFile.Sync()
}
