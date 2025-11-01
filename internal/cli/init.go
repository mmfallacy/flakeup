package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mmfallacy/flakeup/internal/core"
	"github.com/mmfallacy/flakeup/internal/nix"
)

type GlobalOptions struct {
	FlakePath string
}

type InitOptions struct {
	GlobalOptions
	Template string
	OutDir   string
}

var conflictActionChoices = []core.ConflictAction{
	core.ConflictPrepend,
	core.ConflictAppend,
	core.ConflictOverwrite,
	core.ConflictIgnore,
}

func HandleInit(opts InitOptions) error {
	fmt.Printf("Cloning template %s from flake %s onto %s\n", opts.Template, opts.FlakePath, opts.OutDir)

	if hasOutput, err := nix.HasFlakeOutput(opts.FlakePath, "flakeupTemplates"); err != nil {
		return fmt.Errorf("init: %w: %w", ErrCliUnexpected, err)
	} else if !hasOutput {
		return fmt.Errorf("init: %w", ErrCliInitMissingFlakeupTemplateOutput)
	}

	templates, err := nix.GetFlakeOutput[core.Templates](opts.FlakePath, "flakeupTemplates")

	if err != nil {
		return fmt.Errorf("init: %s", err)
	}

	dir, err := os.MkdirTemp("", "flakeup-")
	if err != nil {
		fmt.Printf("Encountered an error! %w\n", err)
		return err
	}

	fmt.Println("Created tmp dir at", dir)

	// Cleanup
	defer func() {
		return
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("Encountered an error! %w\n", err)
		}
	}()

	actions, err := templates[opts.Template].Process(opts.OutDir)

	if err != nil {
		fmt.Printf("Encountered an error! %w\n", err)
		return err
	}

	for i := range actions {
		// For all actions, temporarily set Dest to tempdir to enable rollbacks on failures
		// Use index to access as actual struct, not by a copy
		switch action := actions[i].(type) {

		default:
			return fmt.Errorf("init: %w: unhandled action type", ErrCliUnexpected)

		// Unfortunately, we can't merge these two cases in the same area as go cannot validate a Dest property
		case *core.ActionApply:
			action.Dest = dir
			continue

		case *core.ActionMkdir:
			action.Dest = dir
			continue

		case *core.ActionAsk:
			answer, err := ask(fmt.Sprintf("conflict at %s", filepath.Join(action.Dest, action.Path)), conflictActionChoices)

			if err != nil {
				return err
			}

			actions[i] = core.ActionApply{
				Desc:    string(answer),
				Src:     action.Src,
				Dest:    dir,
				Path:    action.Path,
				Pattern: action.Pattern,
				Rule: core.Rule{
					OnConflict: &answer,
				},
				Write: true,
			}
		}

	}

	for _, a := range actions {
		_ = a.Process()
	}

	return nil
}
