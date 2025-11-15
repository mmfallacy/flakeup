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

func push(s *[]ActionEntry, el *ActionEntry) error {
	*s = append(*s, *el)
	return nil
}

type Template struct {
	Description *string      `json:"description"`
	Root        *string      `json:"root"`
	Parameters  *[]Parameter `json:"parameters"`
	Rules       *Rules       `json:"rules"`
}

func (T Template) Process(outdir string) ([]ActionEntry, error) {
	if T.Root == nil {
		return nil, fmt.Errorf("template processing: no specified root")
	}
	root := *T.Root

	if !strings.HasPrefix(root, "/nix/store/") {
		fmt.Println("WARNING: Template path not in the /nix/store/")
	}

	sortedRuleKeys := make([]string, 0)
	if T.Rules != nil {
		sortedRuleKeys = u.SortKeysByLength(u.GetKeys(*T.Rules))
	}

	if err := u.AssertEach(sortedRuleKeys, func(el string) bool {
		return doublestar.ValidatePattern(el)
	}); err != nil {
		return nil, fmt.Errorf("template processing: invalid pattern glob! %v", doublestar.ErrBadPattern)
	}

	// Create dynamic list with default capacity 10
	actions := make([]ActionEntry, 0, 10)

	// Add action to create output directory
	push(&actions, &ActionEntry{
		Desc:    "mkdir output directory",
		Pattern: "",
		Action:  &Mkdir{Dest: u.Path{Root: outdir, Rel: ""}},
	})

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
			return push(&actions, &ActionEntry{
				Desc:    fmt.Sprintf("mkdir %s", path),
				Pattern: "",
				Action:  &Mkdir{Dest: u.Path{Root: outdir, Rel: path}},
			})
		default:
			fmt.Println("WARNING: Skipping", path, "as it is neither a regular file or a directory!")
			return nil

		case mode.IsRegular():
			// continue
		}

		var pattern string
		var match Rule

		if _, err := os.Stat(filepath.Join(outdir, path)); os.IsNotExist(err) {
			return push(&actions, &ActionEntry{
				Desc:    fmt.Sprintf("exact %s", path),
				Pattern: "",
				Action: &Exact{
					Src:  u.Path{Root: root, Rel: path},
					Dest: u.Path{Root: outdir, Rel: path},
				},
			})
		}

		for _, key := range sortedRuleKeys {
			if ok := doublestar.MatchUnvalidated(key, path); ok {
				// When nil, the for loop already terminates so no need to check before dereferencing
				match = (*T.Rules)[key]
				pattern = key
				break
			}
		}

		// Ask by default on no matching rules
		if match == (Rule{}) && pattern == "" {
			return push(&actions, &ActionEntry{
				Desc:    fmt.Sprintf("ask %s", path),
				Pattern: "",
				Action: &Ask{
					Src:  u.Path{Root: root, Rel: path},
					Dest: u.Path{Root: outdir, Rel: path},
				},
			})
		}

		// Handle onConflict Rules
		switch *match.OnConflict {
		case ConflictPrepend:
			return push(&actions, &ActionEntry{
				Desc:    fmt.Sprintf("prepend %s", path),
				Pattern: pattern,
				Action: &Prepend{
					Base:   u.Path{Root: root, Rel: path},
					Prefix: u.Path{Root: outdir, Rel: path},
					Dest:   u.Path{Root: outdir, Rel: path},
				},
			})
		case ConflictAppend:
			return push(&actions, &ActionEntry{
				Desc:    fmt.Sprintf("append %s", path),
				Pattern: pattern,
				Action: &Append{
					Base:   u.Path{Root: root, Rel: path},
					Suffix: u.Path{Root: outdir, Rel: path},
					Dest:   u.Path{Root: outdir, Rel: path},
				},
			})
		case ConflictOverwrite:
			return push(&actions, &ActionEntry{
				Desc:    fmt.Sprintf("overwrite %s", path),
				Pattern: pattern,
				Action: &Overwrite{
					Src:  u.Path{Root: root, Rel: path},
					Dest: u.Path{Root: outdir, Rel: path},
				},
			})
		case ConflictIgnore:
			return push(&actions, &ActionEntry{
				Desc:    fmt.Sprintf("ignore %s", path),
				Pattern: pattern,
				Action: &Ignore{
					Src:  u.Path{Root: root, Rel: path},
					Dest: u.Path{Root: outdir, Rel: path},
				},
			})
		// Ask by default
		// case ConflictAsk:
		default:
			return push(&actions, &ActionEntry{
				Desc:    fmt.Sprintf("ask %s", path),
				Pattern: pattern,
				Action: &Ask{
					Src:  u.Path{Root: root, Rel: path},
					Dest: u.Path{Root: outdir, Rel: path},
				},
			})
		}

	})
	return actions, nil
}
