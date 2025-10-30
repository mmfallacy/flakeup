package cli

import (
	"fmt"

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

	actions, err := templates[opts.Template].Process(opts.OutDir)

	if err != nil {
		fmt.Printf("Encountered an error! %w\n", err)
		return err
	}

	for i, a := range actions {
		action, ok := a.(core.ActionAsk)
		if !ok {
			continue
		}

		answer, err := ask(fmt.Sprintf("conflict at %s", action.Dest), conflictActionChoices)

		if err != nil {
			return err
		}

		actions[i] = core.ActionApply{
			Desc:    string(answer),
			Src:     action.Src,
			Dest:    action.Dest,
			Pattern: action.Pattern,
			Rule: core.Rule{
				OnConflict: &answer,
			},
			Write: true,
		}
	}

	for _, a := range actions {
		a.Process()
	}

	return nil
}
