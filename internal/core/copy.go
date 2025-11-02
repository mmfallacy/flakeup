package core

import (
	"fmt"
	"io"
	"os"
)

// Assumption: file at dst already exists
func CopyPrepend(src, dst string) error {
	pre, err := os.Open(src)

	if err != nil {
		return err
	}

	defer pre.Close()

	// Backup existing file
	bak := fmt.Sprintf("%s.bak")
	if err := os.Rename(dst, bak); err != nil {
		return err
	}

	post, err := os.Open(bak)

	if err != nil {
		return err
	}

	defer post.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}

	defer out.Close()

	if _, err = io.Copy(out, pre); err != nil {
		return err
	}

	if _, err = io.Copy(out, post); err != nil {
		return err
	}

	// Remove backup
	os.Remove(bak)
	return out.Sync()
}

// Assumption: file at dst already exists
func CopyAppend(src, dst string) error {
	in, err := os.Open(src)

	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}

	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

// Assumption: file at dst does not exist
func CopyClean(src, dst string) error {
	fmt.Println("Copying ", src, " to ", dst)
	in, err := os.Open(src)

	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_EXCL, 0o644)

	if err != nil {
		return err
	}

	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

// Assumption: file at dst exists
func CopyOverwrite(src, dst string) error {
	fmt.Println("Copying ", src, " to ", dst)
	in, err := os.Open(src)

	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC, 0o644)

	if err != nil {
		return err
	}

	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}
