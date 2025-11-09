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
		// For all actions, temporarily set Dest to tempdir to enable rollbacks on failure

		switch action := actions[i].Action.(type) {
		default:
			return fmt.Errorf("init: %w: unsupported action type", ErrCliUnexpected)
		case *core.Mkdir:
			action.Dest = utils.Path{Root: dir, Rel: action.Dest.Rel}
		case *core.Exact:
			action.Dest = utils.Path{Root: dir, Rel: action.Dest.Rel}
		case *core.Append:
			action.Dest = utils.Path{Root: dir, Rel: action.Dest.Rel}
		case *core.Prepend:
			action.Dest = utils.Path{Root: dir, Rel: action.Dest.Rel}
		// noop, so don't bother resetting Dest
		case *core.Ignore:
			continue

		case *core.Ask:
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

	}

	fmt.Println(utils.Prettify(actions))
	for _, a := range actions {
		_ = a.Process()
	}

	fmt.Println("Created tmp dir copy at", dir)

	return nil
}
