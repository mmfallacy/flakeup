package core

import (
	"bytes"
	"fmt"
	"io"
)

type SubRule struct {
	Find    []byte
	Replace []byte
}

type Substituter struct {
	reader io.Reader
	rules  []SubRule
	// window is a buffer that contains previous tail and current read bytes.
	// Read(p) updates window and does the substitutions.
	window []byte
	nTail  int
}

// Window grows depending on len(p) on Read(p) call.
// Preallocate 4 * 1024 bytes instead to reduce overhead on reallocations
const MIN_WINDOW_LENGTH = 4 * 1024

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
		window: make([]byte, 0, MIN_WINDOW_LENGTH),
		nTail:  maxlen,
	}
}

func (s *Substituter) Debug() error {
	for {
		buf := make([]byte, 32*1024)
		nr, err := s.Read(buf)
		fmt.Printf("Read %d bytes from reader. err %+v\n", nr, err)

		if err != nil && err != io.EOF {
			return err
		}

		if err == io.EOF {
			break
		}

		if nr > 0 {
			fmt.Print(string(buf[:nr]))
		}
	}
	return nil
}

// implement io.Reader interface
func (s *Substituter) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	// Flush safe bytes from window
	// Always keep at least s.nTail bytes to accomodate cross-chunk replacements
	if len(s.window) >= len(p) {
		keep := min(len(p), len(s.window)-s.nTail)
		n := copy(p, s.window[:keep])
		s.window = s.window[n:]
		return n, nil
	}

	tmp := make([]byte, len(p))
	if _, err := s.reader.Read(tmp); err != nil {
		return 0, err
	}

	s.window = append(s.window, tmp...)

	// Apply substitution rules
	for _, rule := range s.rules {
		s.window = bytes.ReplaceAll(s.window, rule.Find, rule.Replace)
	}

	// Strictly len(window) > len(p) as we're appending tmp (of length len(p))
	// TODO: check if at this step len(window) = len(p) + s.nTail
	n := copy(p, s.window[:len(tmp)])
	s.window = s.window[n:]
	return n, nil
}

// Design choices:
// - ~~Manually reallocate when append() exceeds capacity~~
// - preallocate decent amount to save overhead for reallocation

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
