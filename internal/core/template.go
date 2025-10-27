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

func (T Template) Process(outdir string) error {
	root := *T.Root

	if !strings.HasPrefix(root, "/nix/store/") {
		fmt.Println("WARNING: Template path not in the /nix/store/")
	}

	sortedRuleKeys := u.SortKeysByLength(u.GetKeys(*T.Rules))

	if err := u.AssertEach(sortedRuleKeys, func(el string) bool {
		return doublestar.ValidatePattern(el)
	}); err != nil {
		return fmt.Errorf("template processing: invalid pattern glob! %v", doublestar.ErrBadPattern)
	}

	// Ignore if already created.
	if err := os.MkdirAll(outdir, 0o755); err != nil {
		return err
	}

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
		default:
			fmt.Println("WARNING: Skipping", path, "as it is neither a regular file or a directory!")
			return nil

		case mode.IsRegular():
			// continue
		}

		var pattern string
		var match Rule

		for _, key := range sortedRuleKeys {
			if ok := doublestar.MatchUnvalidated(key, path); ok {
				match = (*T.Rules)[key]
				pattern = key
				break
			}
		}

		// Raw copy on no matching rules
		if match == (Rule{}) && pattern == "" {
			return Copy(rootpath, outpath)
		}

		fmt.Printf("Here's the matching rule for path %s:\n %s:\n%s\n", path, pattern, u.Prettify(match))
		return Copy(rootpath, outpath)
	})
}
