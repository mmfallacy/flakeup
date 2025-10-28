package cli

import (
	"fmt"

	"github.com/mmfallacy/flakeup/internal/core"
	"github.com/mmfallacy/flakeup/internal/nix"
	"github.com/mmfallacy/flakeup/internal/utils"
	_ "github.com/mmfallacy/flakeup/internal/utils"
)

type GlobalOptions struct {
	FlakePath string
}

type InitOptions struct {
	GlobalOptions
	Template string
	OutDir   string
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
	} else {
		fmt.Printf("Here are the actions needed to be taken:\n %s\n", utils.Prettify(actions))
	}

	return nil
}
