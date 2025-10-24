package utils

import (
	"fmt"
	"os"
)

func Panic(msg string, err error) {
	fmt.Printf("Unexpected error: %s\n%v", msg, err)
	os.Exit(1)
}
