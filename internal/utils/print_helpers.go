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

func Prettify(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		// Return fixed string instead on failure so fn can be used
		return "json marshall failed"
	}
	return string(b)
}
