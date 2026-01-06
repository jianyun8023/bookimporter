package e2e

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenameCommand(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Prepare files
	files := []string{"file1.txt", "file2.txt", "doc.pdf"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tempDir, f), []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	// 2. Test Rename Pattern
	t.Run("Rename TXT Files", func(t *testing.T) {
		output, err := runCommand("rename", tempDir, "-f", "txt", "-t", "book-@n")
		if err != nil {
			t.Errorf("Rename command failed: %v, output: %s", err, output)
		}

		// Verify changes
		expectedFiles := []string{"book-1.txt", "book-2.txt", "doc.pdf"}
		for _, f := range expectedFiles {
			if _, err := os.Stat(filepath.Join(tempDir, f)); os.IsNotExist(err) {
				t.Errorf("Expected file %s not found", f)
			}
		}
		// Verify old files gone
		for _, f := range []string{"file1.txt", "file2.txt"} {
			if _, err := os.Stat(filepath.Join(tempDir, f)); !os.IsNotExist(err) {
				t.Errorf("Old file %s still exists", f)
			}
		}
	})

	// 3. Test Move to Output Dir
	t.Run("Rename and Move", func(t *testing.T) {
		// Create new set of files
		subDir := filepath.Join(tempDir, "source")
		os.Mkdir(subDir, 0755)
		os.WriteFile(filepath.Join(subDir, "pic.jpg"), []byte("pic"), 0644)

		outDir := filepath.Join(tempDir, "output")

		output, err := runCommand("rename", subDir, "-f", "jpg", "-t", "image-@n", "-o", outDir)
		if err != nil {
			t.Errorf("Rename command failed: %v, output: %s", err, output)
		}

		// Verify
		if _, err := os.Stat(filepath.Join(outDir, "image-1.jpg")); os.IsNotExist(err) {
			t.Errorf("Expected file image-1.jpg in output dir")
		}
		if _, err := os.Stat(filepath.Join(subDir, "pic.jpg")); !os.IsNotExist(err) {
			t.Errorf("Source file should be gone")
		}
	})

	// 4. Test Preview
	t.Run("Preview Mode", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "preview.txt")
		os.WriteFile(testFile, []byte("test"), 0644)

		output, err := runCommand("rename", tempDir, "-f", "txt", "-t", "real-@n", "--do-try")
		if err != nil {
			t.Errorf("Rename command failed: %v", err)
		}

		// Verify output indicates preview
		// Depending on table output...

		// Verify file NOT changed
		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Errorf("File should exist in preview mode")
		}
		if _, err := os.Stat(filepath.Join(tempDir, "real-1.txt")); !os.IsNotExist(err) {
			t.Errorf("New file should NOT exist in preview mode, but found one or other error. Output: %s", output)
		}
	})
}
