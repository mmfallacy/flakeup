package core

import (
	"io"
	"os"
)

func MergeIntoWithSubstitutions(a string, b *string, c string, sub Substitutions) error {
	pre, err := os.Open(a)
	if err != nil {
		return err
	}
	defer pre.Close()

	out, err := os.OpenFile(c, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()

	presub := NewSubstituter(pre, sub)
	if _, err = io.Copy(out, presub); err != nil {
		return err
	}

	if b == nil {
		return out.Sync()
	}

	// Append contents of b
	post, err := os.Open(*b)
	if err != nil {
		return err
	}
	defer post.Close()

	postsub := NewSubstituter(post, sub)
	if _, err = io.Copy(out, postsub); err != nil {
		return err
	}

	return out.Sync()

}

// b is nil-able
func MergeInto(a string, b *string, c string) error {
	pre, err := os.Open(a)

	if err != nil {
		return err
	}

	defer pre.Close()

	// Presume file at c does not exist, throw otherwise
	out, err := os.OpenFile(c, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)

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

	if _, err = io.Copy(out, post); err != nil {
		return err
	}

	return out.Sync()
}
