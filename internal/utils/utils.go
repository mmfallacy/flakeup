package utils

import (
	"errors"
	"sort"
)

var ErrAssertion = errors.New("Assertion error")

func GetKeys[M ~map[K]V, K string, V any](m M) []K {
	keys := make([]K, 0, len(m))

	for k, _ := range m {
		keys = append(keys, k)
	}

	return keys
}

func SortKeysByLength(keys []string) []string {
	sort.Slice(keys, func(i, j int) bool {
		return len(keys[i]) > len(keys[j])
	})
	return keys
}

func AssertEach[T any](slice []T, fn func(T) bool) error {
	for _, v := range slice {
		if ok := fn(v); !ok {
			return ErrAssertion
		}
	}
	return nil
}

var ErrNotImplemented = errors.New("Function not yet implemented")

func LooseMapStringToType[T ~string](val string, choices []T) (T, bool) {
	for _, choice := range choices {
		if len(val) < 1 {
			return "", false
		}
		if val == string(choice) {
			return choice, true
		}
		// prefix matching
		if len(choice) >= len(val) && string(choice)[:len(val)] == val {
			return choice, true
		}
	}

	return "", false
}
