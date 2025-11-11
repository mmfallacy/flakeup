package core

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyRecursiveOverwrite(src, dest string) error {
	return fs.WalkDir(os.DirFS(src), ".", func(path string, d fs.DirEntry, err error) error {
		// non-nil return value crashes walk
		if err != nil {
			return nil
		}
		if path == "." || path == ".." {
			return nil
		}

		srcpath := filepath.Join(src, path)
		destpath := filepath.Join(dest, path)

		switch mode := d.Type(); {
		case mode.IsDir():
			return os.MkdirAll(destpath, 0o755)
		default:
			fmt.Println("WARNING: Skipping", path, "as it is neither a regular file or a directory!")
			return nil
		case mode.IsRegular():
			// noop
		}

		in, err := os.Open(srcpath)

		// Crash whole walk to signal irrecoverable state
		// Ideally, this should not happen as srcpath is created by flakeup in a tempdir
		if err != nil {
			return err
		}

		defer in.Close()

		out, err := os.OpenFile(destpath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)

		// Ideally this should not happen as destpath as already been accounted for while processing template
		// TODO: Investigate when and how can MkdirAll and io.Copy fail and account it within template processing
		if err != nil {
			return err
		}

		defer out.Close()

		// This should also not happen
		// TODO: Consider never cases to panic instead of returning errors
		if _, err = io.Copy(out, in); err != nil {
			return err
		}

		return out.Sync()
	})
}
