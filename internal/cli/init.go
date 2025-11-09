package cli

import (
	"fmt"
	"os"

	"github.com/mmfallacy/flakeup/internal/core"
	"github.com/mmfallacy/flakeup/internal/nix"
	"github.com/mmfallacy/flakeup/internal/utils"
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
		// Resolve asks first
		if action, ok := actions[i].Action.(*core.Ask); ok {
			answer, err := ask(fmt.Sprintf("conflict at %s", action.Dest.Resolve()), conflictActionChoices)

			if err != nil {
				return err
			}

			resolved, err := action.Resolve(answer)

			if err != nil || resolved == nil {
				return err
			}

			prev := actions[i]
			actions[i] = core.ActionEntry{
				Desc:    "resolved ask",
				Pattern: prev.Pattern,
				Action:  resolved,
			}
		}
		// For all resolved actions, temporarily set Dest to tempdir to enable rollbacks on failure
		switch action := actions[i].Action.(type) {
		default:
			return fmt.Errorf("init: %w: unsupported action type", ErrCliUnexpected)
		case *core.Mkdir:
			action.Dest = utils.Path{Root: dir, Rel: action.Dest.Rel}
		case *core.Exact:
			action.Dest = utils.Path{Root: dir, Rel: action.Dest.Rel}
		case *core.Overwrite:
			action.Dest = utils.Path{Root: dir, Rel: action.Dest.Rel}
		case *core.Append:
			action.Dest = utils.Path{Root: dir, Rel: action.Dest.Rel}
		case *core.Prepend:
			action.Dest = utils.Path{Root: dir, Rel: action.Dest.Rel}
		// noop, so don't bother resetting Dest
		case *core.Ignore:
			continue

		// This should have already been resolved
		case *core.Ask:
			continue
		}

		// Process action entry
		if err = actions[i].Process(); err != nil {
			return err
		}
	}

	// Summarize changes
	for _, action := range actions {
		// Reset tempdir back to outdir
		switch action := action.Action.(type) {
		default:
			return fmt.Errorf("init: %w: unsupported action type", ErrCliUnexpected)
		case *core.Mkdir:
			action.Dest = utils.Path{Root: opts.OutDir, Rel: action.Dest.Rel}
		case *core.Exact:
			action.Dest = utils.Path{Root: opts.OutDir, Rel: action.Dest.Rel}
		case *core.Overwrite:
			action.Dest = utils.Path{Root: opts.OutDir, Rel: action.Dest.Rel}
		case *core.Append:
			action.Dest = utils.Path{Root: opts.OutDir, Rel: action.Dest.Rel}
		case *core.Prepend:
			action.Dest = utils.Path{Root: opts.OutDir, Rel: action.Dest.Rel}
		// noop, so don't bother resetting Dest
		case *core.Ignore:
		// This should have already been resolved
		case *core.Ask:
		}
		fmt.Println(action)
	}

	//TODO: Ask user if they want to apply the template changes

	//TODO: On confirm, apply template from tempdir

	return nil
}
