package main

import (
	"fmt"
	"os"
)

func panic(msg string, err error) {
	fmt.Printf("Unexpected error: %s\n%v", msg, err)
	os.Exit(2)
}
