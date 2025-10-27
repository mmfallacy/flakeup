package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// JSON schema as structs
type Templates map[string]Template

type Template struct {
	Root       *string      `json:"root"`
	Parameters *[]Parameter `json:"parameters"`
	Rules      *Rules       `json:"rules"`
}

type Parameter struct {
	Name *string `json:"name"`
	// Nullable
	Prompt  *string `json:"prompt"`
	Default *string `json:"default"`
}

type Rules map[string]Rule

type Rule struct {
	OnConflict *ConflictAction
}

type ConflictAction string

const (
	ConflictPrepend   ConflictAction = "prepend"
	ConflictAppend    ConflictAction = "append"
	ConflictOverwrite ConflictAction = "overwrite"
	ConflictIgnore    ConflictAction = "ignore"
	ConflictAsk       ConflictAction = "ask"
)

func sortMapKeys[M ~map[K]V, K string, V any](m M) []K {
	sorted := make([]K, 0, len(m))

	for k, _ := range m {
		sorted = append(sorted, k)
	}

	sort.Slice(sorted, func(i, j int) bool {
		return len(sorted[i]) > len(sorted[j])
	})

	return sorted
}

func (T Template) Process(outdir string) error {
	root := *T.Root

	if !strings.HasPrefix(root, "/nix/store/") {
		fmt.Println("WARNING: Template path not in the /nix/store/")
	}

	// Ignore if already created.
	if err := os.MkdirAll(outdir, 0o755); err != nil {
		return err
	}

	sortedRuleKeys := sortMapKeys(*T.Rules)

	// Use fs instead of filepath to keep path strings relative to template root path
	return fs.WalkDir(os.DirFS(root), ".", func(path string, d fs.DirEntry, err error) error {
		// Skip root entry
		if path == "." {
			return nil
		}

		fmt.Println("Walking", path)

		rootpath := filepath.Join(root, path)
		outpath := filepath.Join(outdir, path)

		// Skip Path if error arises when opening. Non-nil return crashes the whole walk.
		if err != nil {
			return nil
		}

		switch mode := d.Type(); {
		case mode.IsDir():
			return os.MkdirAll(outpath, 0o755)
		case mode.IsRegular():
			return Copy(rootpath, outpath)
		default:
			fmt.Println("WARNING: Skipping", path, "as it is neither a regular file or a directory!")
			return nil
		}

	})
}
