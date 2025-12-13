package core_test

import (
	"os"
	"path/filepath"
	"testing"

	c "github.com/mmfallacy/flakeup/internal/core"
)

func TestCopyRecursiveOverwrite(t *testing.T) {
	// Test case 1: Copy a simple file
	t.Run("copy simple file", func(t *testing.T) {
		srcDir := t.TempDir()
		destDir := t.TempDir()

		// Create a source file
		srcFilePath := filepath.Join(srcDir, "test.txt")
		err := os.WriteFile(srcFilePath, []byte("hello world"), 0644)
		if err != nil {
			t.Fatalf("failed to create source file: %v", err)
		}

		err = c.CopyRecursiveOverwrite(srcDir, destDir)
		if err != nil {
			t.Fatalf("CopyRecursiveOverwrite failed: %v", err)
		}

		// Verify the file exists in the destination
		destFilePath := filepath.Join(destDir, "test.txt")
		content, err := os.ReadFile(destFilePath)
		if err != nil {
			t.Fatalf("failed to read destination file: %v", err)
		}
		if string(content) != "hello world" {
			t.Errorf("unexpected content in destination file: %s", string(content))
		}
	})

	// Test case 2: Copy a directory with subdirectories and files
	t.Run("copy directory with subdirectories and files", func(t *testing.T) {
		srcDir := t.TempDir()
		destDir := t.TempDir()

		// Create source structure
		os.MkdirAll(filepath.Join(srcDir, "subdir1", "subsubdir"), 0755)
		os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
		os.WriteFile(filepath.Join(srcDir, "subdir1", "file2.txt"), []byte("content2"), 0644)
		os.WriteFile(filepath.Join(srcDir, "subdir1", "subsubdir", "file3.txt"), []byte("content3"), 0644)

		err := c.CopyRecursiveOverwrite(srcDir, destDir)
		if err != nil {
			t.Fatalf("CopyRecursiveOverwrite failed: %v", err)
		}

		// Verify destination structure and content
		if _, err := os.Stat(filepath.Join(destDir, "file1.txt")); os.IsNotExist(err) {
			t.Errorf("file1.txt not copied")
		}
		if _, err := os.Stat(filepath.Join(destDir, "subdir1", "file2.txt")); os.IsNotExist(err) {
			t.Errorf("subdir1/file2.txt not copied")
		}
		if _, err := os.Stat(filepath.Join(destDir, "subdir1", "subsubdir", "file3.txt")); os.IsNotExist(err) {
			t.Errorf("subdir1/subsubdir/file3.txt not copied")
		}

		content, _ := os.ReadFile(filepath.Join(destDir, "file1.txt"))
		if string(content) != "content1" {
			t.Errorf("unexpected content for file1.txt")
		}
	})

	// Test case 3: Overwrite existing file
	t.Run("overwrite existing file", func(t *testing.T) {
		srcDir := t.TempDir()
		destDir := t.TempDir()

		// Create source file
		srcFilePath := filepath.Join(srcDir, "overwrite.txt")
		os.WriteFile(srcFilePath, []byte("new content"), 0644)

		// Create existing destination file
		destFilePath := filepath.Join(destDir, "overwrite.txt")
		os.WriteFile(destFilePath, []byte("old content"), 0644)

		err := c.CopyRecursiveOverwrite(srcDir, destDir)
		if err != nil {
			t.Fatalf("CopyRecursiveOverwrite failed: %v", err)
		}

		content, err := os.ReadFile(destFilePath)
		if err != nil {
			t.Fatalf("failed to read destination file: %v", err)
		}
		if string(content) != "new content" {
			t.Errorf("file not overwritten, got: %s", string(content))
		}
	})

	// Test case 4: Copy empty directory
	t.Run("copy empty directory", func(t *testing.T) {
		srcDir := t.TempDir()
		destDir := t.TempDir()

		err := c.CopyRecursiveOverwrite(srcDir, destDir)
		if err != nil {
			t.Fatalf("CopyRecursiveOverwrite failed: %v", err)
		}

		// Verify destination directory exists and is empty
		entries, err := os.ReadDir(destDir)
		if err != nil {
			t.Fatalf("failed to read destination directory: %v", err)
		}
		if len(entries) != 0 {
			t.Errorf("destination directory not empty, found %d entries", len(entries))
		}
	})
}
