package utils

import (
	"encoding/json"
)

func Prettify(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		// Return fixed string instead on failure so fn can be used
		return "json marshall failed"
	}
	return string(b)
}
