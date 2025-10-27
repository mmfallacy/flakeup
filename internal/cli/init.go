package cli

import (
	"fmt"

	"github.com/mmfallacy/flakeup/internal/core"
	"github.com/mmfallacy/flakeup/internal/nix"
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

	templates[opts.Template].Process(opts.OutDir)
	return nil
}
