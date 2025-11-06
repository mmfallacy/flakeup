package core

import (
	"io"
	"os"
)

// b is nil-able
func MergeInto(a string, b *string, c string) error {
	pre, err := os.Open(a)

	if err != nil {
		return err
	}

	defer pre.Close()

	// Presume file at c does not exist, throw otherwise
	out, err := os.OpenFile(c, os.O_CREATE|os.O_EXCL, 0o644)

	if err != nil {
		return err
	}

	defer out.Close()

	if _, err = io.Copy(out, pre); err != nil {
		return err
	}

	// When b is nil, no more stuff needs to be appended
	if b == nil {
		return out.Sync()
	}

	post, err := os.Open(*b)

	if err != nil {
		return err
	}

	defer post.Close()

	if _, err = io.Copy(out, pre); err != nil {
		return err
	}

	return out.Sync()
}
