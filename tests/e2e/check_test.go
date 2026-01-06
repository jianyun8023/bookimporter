package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckCommand(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Prepare files
	validEpub := filepath.Join(tempDir, "valid.epub")
	if err := createValidEpub(validEpub, "Valid Book"); err != nil {
		t.Fatalf("Failed to create valid epub: %v", err)
	}

	corruptedEpub := filepath.Join(tempDir, "corrupted.epub")
	if err := createCorruptedEpub(corruptedEpub); err != nil {
		t.Fatalf("Failed to create corrupted epub: %v", err)
	}

	// 2. Test Check Valid
	t.Run("Check Valid File", func(t *testing.T) {
		output, err := runCommand("check", "-p", validEpub)
		if err != nil {
			t.Errorf("Check command failed: %v, output: %s", err, output)
		}
		if !strings.Contains(output, "通过") {
			t.Errorf("Expected output to contain '通过', got: %s", output)
		}
	})

	// 3. Test Check Corrupted
	t.Run("Check Corrupted File", func(t *testing.T) {
		output, err := runCommand("check", "-p", corruptedEpub)
		if err != nil {
			// check command currently does not return non-zero exit code for corruption found unless it crashes?
			// Actually svc.Run returns nil even if failures found, unless exit code logic is implemented in cmd.
			// In cmd/check.go: if err := svc.Run(...); err != nil { os.Exit(1) }
			// check.Service.Run only returns error if scan failed or critical error. Individual file check error is logged.
			// So exit code should be 0.
			// But wait, existing cmd/check.go logic:
			// stats.Failed > 0 does NOT trigger error return in check.Service.Run (it returns nil).
			// So exit code is 0. This is expected behavior for now.
		}
		if !strings.Contains(output, "失败") {
			t.Errorf("Expected output to contain '失败', got: %s", output)
		}
	})

	// 4. Test Move Corrupted
	t.Run("Move Corrupted File", func(t *testing.T) {
		moveToDir := filepath.Join(tempDir, "corrupted_files")
		output, err := runCommand("check", "-p", corruptedEpub, "--move-to", moveToDir)
		if err != nil {
			t.Errorf("Check command failed: %v", err)
		}

		if !strings.Contains(output, "已移动") {
			t.Errorf("Expected output to contain '已移动', got: %s", output)
		}

		// Verify file moved
		movedFile := filepath.Join(moveToDir, "corrupted.epub")
		if _, err := os.Stat(movedFile); os.IsNotExist(err) {
			t.Errorf("File was not moved to %s", movedFile)
		}
		if _, err := os.Stat(corruptedEpub); !os.IsNotExist(err) {
			t.Errorf("Original file %s still exists", corruptedEpub)
		}
	})

	// Prepare another corrupted file for Delete test
	corruptedEpub2 := filepath.Join(tempDir, "corrupted2.epub")
	createCorruptedEpub(corruptedEpub2)

	// 5. Test Delete Corrupted
	t.Run("Delete Corrupted File", func(t *testing.T) {
		// Use --force to skip interactive confirmation (we assume force works as intended in CLI arg parsing)
		// But in our refactoring, we removed interactive confirmation temporarily?
		// No, I added it back in the last step.
		// So we MUST use --force to avoid getting stuck on stdin read.

		output, err := runCommand("check", "-p", corruptedEpub2, "--delete", "--force")
		if err != nil {
			t.Errorf("Check command failed: %v, output: %s", err, output)
		}

		if !strings.Contains(output, "已删除") {
			t.Errorf("Expected output to contain '已删除', got: %s", output)
		}

		if _, err := os.Stat(corruptedEpub2); !os.IsNotExist(err) {
			t.Errorf("File %s was not deleted", corruptedEpub2)
		}
	})
}
