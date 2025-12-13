package core_test

import (
	"os"
	"path/filepath"
	"testing"

	c "github.com/mmfallacy/flakeup/internal/core"
	u "github.com/mmfallacy/flakeup/internal/utils"
)

func TestMkdirAction_Do(t *testing.T) {
	t.Run("create new directory", func(t *testing.T) {
		tempDir := t.TempDir()
		newDirName := "test_dir"
		expectedPath := filepath.Join(tempDir, newDirName)

		action := &c.Mkdir{
			Dest: u.Path{Root: tempDir, Rel: newDirName},
		}

		err := action.Do(nil) // Mkdir doesn't use substitutions
		if err != nil {
			t.Fatalf("Mkdir.Do failed: %v", err)
		}

		// Verify directory was created
		info, err := os.Stat(expectedPath)
		if os.IsNotExist(err) {
			t.Errorf("directory %s was not created", expectedPath)
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", expectedPath)
		}
	})

	t.Run("create nested directory", func(t *testing.T) {
		tempDir := t.TempDir()
		newDirPath := "parent/child/grandchild"
		expectedPath := filepath.Join(tempDir, newDirPath)

		action := &c.Mkdir{
			Dest: u.Path{Root: tempDir, Rel: newDirPath},
		}

		err := action.Do(nil)
		if err != nil {
			t.Fatalf("Mkdir.Do failed for nested directory: %v", err)
		}

		// Verify nested directory was created
		info, err := os.Stat(expectedPath)
		if os.IsNotExist(err) {
			t.Errorf("nested directory %s was not created", expectedPath)
		}
		if !info.IsDir() {
			t.Errorf("%s is not a directory", expectedPath)
		}
	})
}
