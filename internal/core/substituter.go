package core

import (
	"bytes"
	"io"
)

type SubRule struct {
	Find    []byte
	Replace []byte
}

type Substituter struct {
	reader io.Reader
	rules  []SubRule
	// Tail contains the previous remaining len(tail) bytes
	// (i.e. last buffer's buf[len(buf)-len(tail)]) to account for cross-chunk reads
	tail []byte
}

func NewSubstituter(r io.Reader, patterns map[string]string) *Substituter {
	rules := make([]SubRule, 0, len(patterns))

	// Map pattern into proper []SubRule while keeping track of max length
	maxlen := 0
	for find, replace := range patterns {
		if len(find) > maxlen {
			maxlen = len(find)
		}

		rules = append(rules, SubRule{
			Find:    []byte(find),
			Replace: []byte(replace),
		})
	}

	return &Substituter{
		reader: r,
		rules:  rules,
		tail:   make([]byte, 0, maxlen),
	}
}

// implement io.Reader interface
func (s *Substituter) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	// Create temporary buffer to not modify p yet
	window := make([]byte, 0, len(p))
	// Apply leftover replacements from last read
	copy(window, s.tail)
	// Prefill remaining bytes of window from Reader
	
	if _, err := s.reader.Read()
	

	
}

// Design choices:
// - Manually reallocate when append() exceeds capacity

// Game plan:
// Consistent behavior:
// Return 0, nil on len(p) == 0
// Flush at most len(p) bytes
// Call 0: 
// Read(p) from reader
// last tail bytes of current p cannot yet be flushed as cross-chunk replacements might happen
// ie. safeFlush := current p - tail
// do bytes.ReplaceAll for safeFlush
// yield len(p) - tail
// push last tail bytes of current p onto window (persists to next calls)
// Subsequent calls:
// If len(window) > len(p),
// --> safeFlush = window[:len(p)]
// --> slide window; window = window[len(p):]
// Otherwise
// Read(p) from reader
// new window = [...window, ...p]
// Now, hopefully len(window) == tail bec of previous clause
// do bytes.ReplaceAll on window
// safeFlush = window[:len(p)]
// window = window[len(p):]

