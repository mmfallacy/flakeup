package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	u "github.com/mmfallacy/flakeup/internal/utils"
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

func push(s *[]Action, el Action) error {
	*s = append(*s, el)
	return nil
}

func (T Template) Process(outdir string) ([]Action, error) {
	root := *T.Root

	if !strings.HasPrefix(root, "/nix/store/") {
		fmt.Println("WARNING: Template path not in the /nix/store/")
	}

	sortedRuleKeys := u.SortKeysByLength(u.GetKeys(*T.Rules))

	if err := u.AssertEach(sortedRuleKeys, func(el string) bool {
		return doublestar.ValidatePattern(el)
	}); err != nil {
		return nil, fmt.Errorf("template processing: invalid pattern glob! %v", doublestar.ErrBadPattern)
	}

	// Create dynamic list with default capacity 10
	actions := make([]Action, 0, 10)

	// Add action to create output directory
	push(&actions, &ActionMkdir{Desc: "mkdir output directory", Dest: outdir, Path: ""})

	// Use fs instead of filepath to keep path strings relative to template root path
	fs.WalkDir(os.DirFS(root), ".", func(path string, d fs.DirEntry, err error) error {
		// Skip root entry or cases where opening path results into an error
		if path == "." || err != nil {
			return nil
		}

		switch mode := d.Type(); {
		// On directory type, add action to mimic template dir tree
		// This should always push a parent dir before pushing children actions
		case mode.IsDir():
			return push(&actions, &ActionMkdir{Desc: fmt.Sprintf("mkdir %s", outdir), Dest: outdir, Path: path})
		default:
			fmt.Println("WARNING: Skipping", path, "as it is neither a regular file or a directory!")
			return nil

		case mode.IsRegular():
			// continue
		}

		var pattern string
		var match Rule

		if _, err := os.Stat(filepath.Join(outdir, path)); os.IsNotExist(err) {
			return push(&actions, &ActionApply{
				Desc:    "no existing file",
				Src:     root,
				Dest:    outdir,
				Path:    path,
				Pattern: "",
				Rule:    Rule{},
				Write:   true,
			})
		}

		for _, key := range sortedRuleKeys {
			if ok := doublestar.MatchUnvalidated(key, path); ok {
				match = (*T.Rules)[key]
				pattern = key
				break
			}
		}

		// Raw copy on no matching rules
		if match == (Rule{}) && pattern == "" {
			return push(&actions, &ActionApply{
				Desc:    "no matching rule",
				Src:     root,
				Dest:    outdir,
				Path:    path,
				Pattern: "",
				Rule:    Rule{},
				Write:   true,
			})
		}

		// Handle onConflict Rules
		// TODO: Refactor! This awfully seems too verbose with multiple sources of truth and unnecessary branching.
		switch *match.OnConflict {
		case ConflictPrepend:
			return push(&actions, &ActionApply{
				Desc:    "prepend",
				Src:     root,
				Dest:    outdir,
				Path:    path,
				Pattern: pattern,
				Rule:    match,
				Write:   true,
			})
		case ConflictAppend:
			return push(&actions, &ActionApply{
				Desc:    "append",
				Src:     root,
				Dest:    outdir,
				Path:    path,
				Pattern: pattern,
				Rule:    match,
				Write:   true,
			})
		case ConflictOverwrite:
			return push(&actions, &ActionApply{
				Desc:    "overwrite",
				Src:     root,
				Dest:    outdir,
				Path:    path,
				Pattern: pattern,
				Rule:    match,
				Write:   true,
			})
		case ConflictIgnore:
			return push(&actions, &ActionApply{
				Desc:    "ignore",
				Src:     root,
				Dest:    outdir,
				Path:    path,
				Pattern: pattern,
				Rule:    match,
				Write:   false,
			})
		case ConflictAsk:
			return push(&actions, &ActionAsk{
				Desc:    "ask",
				Src:     root,
				Dest:    outdir,
				Path:    path,
				Pattern: pattern,
				Rule:    match,
			})
		default:
			return push(&actions, &ActionAsk{
				Desc:    "ask by default",
				Src:     root,
				Dest:    outdir,
				Path:    path,
				Pattern: pattern,
				Rule:    match,
			})
		}

	})
	return actions, nil
}
