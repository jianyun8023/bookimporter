package e2e

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	binaryName string
)

func TestMain(m *testing.M) {
	// 1. Build the binary
	if runtime.GOOS == "windows" {
		binaryName = "bookimporter.exe"
	} else {
		binaryName = "bookimporter"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryName, "../../main.go")
	if err := buildCmd.Run(); err != nil {
		fmt.Printf("Failed to build binary: %v\n", err)
		os.Exit(1)
	}

	// 2. Run tests
	code := m.Run()

	// 3. Cleanup
	os.Remove(binaryName)

	os.Exit(code)
}

// Helper: Run command
func runCommand(args ...string) (string, error) {
	path, err := filepath.Abs(binaryName)
	if err != nil {
		return "", err
	}
	cmd := exec.Command(path, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Helper: Create valid EPUB
func createValidEpub(path string, title string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// mimetype
	mime, err := w.Create("mimetype")
	if err != nil {
		return err
	}
	if _, err := mime.Write([]byte("application/epub+zip")); err != nil {
		return err
	}

	// container.xml
	container, err := w.Create("META-INF/container.xml")
	if err != nil {
		return err
	}
	containerContent := `<?xml version="1.0"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
   <rootfiles>
      <rootfile full-path="content.opf" media-type="application/oebps-package+xml"/>
   </rootfiles>
</container>`
	if _, err := container.Write([]byte(containerContent)); err != nil {
		return err
	}

	// content.opf
	opf, err := w.Create("content.opf")
	if err != nil {
		return err
	}
	opfContent := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<package xmlns="http://www.idpf.org/2007/opf" unique-identifier="uuid_id" version="2.0">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>%s</dc:title>
  </metadata>
</package>`, title)
	if _, err := opf.Write([]byte(opfContent)); err != nil {
		return err
	}

	return nil
}

// Helper: Create corrupted EPUB
func createCorruptedEpub(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("This is not a zip file")
	return err
}

// Helper: Create mock ebook-meta
func createMockEbookMeta(dir string) (string, error) {
	name := "ebook-meta"
	content := `#!/bin/bash
echo "Called ebook-meta with args: $@"
`
	if runtime.GOOS == "windows" {
		name = "ebook-meta.bat"
		content = `@echo off
echo Called ebook-meta with args: %*
`
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0755); err != nil {
		return "", err
	}
	return path, nil
}

// Helper: Copy file
func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
