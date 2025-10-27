package core

import (
	"fmt"
	"io"
	"os"
)

func Copy(src, dst string) error {
	fmt.Println("Copying ", src, " to ", dst)
	in, err := os.Open(src)

	if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in)

	if err != nil {
		return err
	}

	return out.Sync()
}
