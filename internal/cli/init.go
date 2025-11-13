package cli

import (
	"fmt"
	"os"

	"github.com/mmfallacy/flakeup/internal/core"
	"github.com/mmfallacy/flakeup/internal/nix"
	s "github.com/mmfallacy/flakeup/internal/style"
	"github.com/mmfallacy/flakeup/internal/utils"
)

type GlobalOptions struct {
	FlakePath string
}

type InitOptions struct {
	GlobalOptions *GlobalOptions
	Template      string
	OutDir        string

	DryRun bool
}

var conflictActionChoices = []core.ConflictAction{
	core.ConflictPrepend,
	core.ConflictAppend,
	core.ConflictOverwrite,
	core.ConflictIgnore,
}

func HandleInit(opts InitOptions) error {
	fmt.Println(s.Infof("Cloning template %s from flake %s onto %s", opts.Template, opts.GlobalOptions.FlakePath, opts.OutDir))

	if hasOutput, err := nix.HasFlakeOutput(opts.GlobalOptions.FlakePath, "flakeupTemplates"); err != nil {
		return fmt.Errorf("init: %w: %w", ErrCliUnexpected, err)
	} else if !hasOutput {
		return fmt.Errorf("init: %w", ErrCliInitMissingFlakeupTemplateOutput)
	}

	templates, err := nix.GetFlakeOutput[core.Templates](opts.GlobalOptions.FlakePath, "flakeupTemplates")

	if err != nil {
		return fmt.Errorf("init: %s", err)
	}

	actions, err := templates[opts.Template].Process(opts.OutDir)

	if err != nil {
		fmt.Printf("Encountered an error! %w\n", err)
		return err
	}

	// Resolve all asks first
	for i := range actions {
		if action, ok := actions[i].Action.(*core.Ask); ok {
			answer, err := ask(s.Warnf("%s Conflict at %s", s.Icons.Warn, action.Dest.Resolve()), conflictActionChoices)

			if err != nil {
				return err
			}

			resolved := action.Resolve(answer)

			prev := actions[i]
			actions[i] = core.ActionEntry{
				Desc:    "resolved ask",
				Pattern: prev.Pattern,
				Action:  resolved,
			}
		}
	}

	fmt.Println()
	fmt.Println(s.Info("Summary of changes:"))
	// Summarize changes
	for _, action := range actions {
		switch action := action.Action.(type) {
		default:
			return fmt.Errorf("init: %w: unsupported action type", ErrCliUnexpected)
		case *core.Mkdir:
			fmt.Println(s.Mkdir(action))
		case *core.Exact:
			fmt.Println(s.Clean(action))
		case *core.Overwrite:
			fmt.Println(s.Conflict(action))
		case *core.Append:
			fmt.Println(s.Conflict(action))
		case *core.Prepend:
			fmt.Println(s.Conflict(action))
		case *core.Ignore:
			fmt.Println(s.Ignore(action))
		// This should have already been resolved
		case *core.Ask:

		}
	}

	if opts.DryRun {
		fmt.Println()
		fmt.Println(s.Info("No changes applied as --dry-run is supplied"))
		return nil
	}

	// Create temporary directory to host changes first to not mutate the target directory prematurely
	dir, err := os.MkdirTemp("", "flakeup-")
	if err != nil {
		fmt.Printf("Encountered an error! %w\n", err)
		return err
	}

	// Cleanup
	// NOTE: Checking for errors aren't really actionable here. On error, rely on tmpfs to delete tmpdir on reboot
	defer os.RemoveAll(dir)

	// Apply changes to tempdir first
	for i := range actions {
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

	fmt.Println()
	//Ask user if they want to apply the template changes
	answer, err := ask(s.Info("Apply the changes? "), []string{"yes", "no"})

	if err != nil {
		return err
	}

	if answer == "no" {
		fmt.Println(s.Errf("%s User cancelled", s.Icons.Err))
		return nil
	}

	//On confirm, apply template from tempdir
	if err := core.CopyRecursiveOverwrite(dir, opts.OutDir); err != nil {
		return err
	}

	fmt.Println(s.Successf("%s Succesfully applied template %s onto directory %s", s.Icons.Success, opts.Template, opts.OutDir))
	return nil
}
