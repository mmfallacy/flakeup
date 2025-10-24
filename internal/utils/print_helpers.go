package utils

import (
	"encoding/json"
	"fmt"
	"os"
)

func Panic(msg string, err error) {
	fmt.Printf("Unexpected error: %s\n%v", msg, err)
	os.Exit(1)
}

func Prettify(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
