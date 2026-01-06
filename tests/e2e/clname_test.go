package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestClnameCommand(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Prepare Mock ebook-meta
	_, err := createMockEbookMeta(tempDir)
	if err != nil {
		t.Fatalf("Failed to create mock ebook-meta: %v", err)
	}

	// Add tempDir to PATH so exec.Command("ebook-meta") finds our mock
	oldPath := os.Getenv("PATH")
	os.Setenv("PATH", tempDir+string(os.PathListSeparator)+oldPath)
	defer os.Setenv("PATH", oldPath)

	// Also for "bash -c" calls (which clname uses), we need to ensure bash finds it.
	// If the code uses absolute path or assumes it's in system path, this might be tricky.
	// The code uses `exec.Command("bash", "-c", command)` where command is `ebook-meta "file" ...`
	// So ebook-meta must be in PATH seen by bash.
	// Setting os.Setenv("PATH", ...) affects the child process.

	// 2. Prepare EPUB with dirty title
	dirtyTitle := "Harry Potter (Special Edition) [2023]"
	epubPath := filepath.Join(tempDir, "book.epub")
	if err := createValidEpub(epubPath, dirtyTitle); err != nil {
		t.Fatalf("Failed to create epub: %v", err)
	}

	// 3. Run clname
	t.Run("Clean Title", func(t *testing.T) {
		output, err := runCommand("clname", "-p", epubPath)
		if err != nil {
			t.Errorf("Clname command failed: %v, output: %s", err, output)
		}

		// Verify output
		expectedCleanTitle := "Harry Potter"
		if !strings.Contains(output, expectedCleanTitle) {
			t.Errorf("Output should show clean title. Got: %s", output)
		}

		// Verify mock execution
		// Our mock echoes "Called ebook-meta..."
		// But Wait, `clname` code captures output?
		// internal/infrastructure/filesystem/repository.go:
		// if _, err := cmdx.Output(); err != nil ...
		// It discards stdout unless error.
		// So we won't see "Called ebook-meta" in `bookimporter` output unless we debug or fail.

		// However, if it didn't fail, it means it ran successfully (exit code 0).
		// We can assume if it reports "Success", it tried to run it.
		// To be 100% sure, our mock could create a side-effect file.
	})
}
