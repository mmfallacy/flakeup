package cli

import (
	"fmt"

	"github.com/mmfallacy/flakeup/internal/nix"
	"github.com/mmfallacy/flakeup/internal/utils"
)

func HandleInit(opts InitOptions) error {
	fmt.Printf("Cloning template %s from flake %s\n", opts.Template, opts.FlakePath)

	if hasOutput, err := nix.HasFlakeOutput(opts.FlakePath, "flakeupTemplates"); err != nil {
		return fmt.Errorf("init: %w: %w", ErrCliUnexpected, err)
	} else if !hasOutput {
		return fmt.Errorf("init: %w", ErrCliInitMissingFlakeupTemplateOutput)
	}

	template, err := nix.GetFlakeOutput[Templates](opts.FlakePath, "flakeupTemplates")

	if err != nil {
		return fmt.Errorf("init: %s", err)
	}

	fmt.Println("Got flake output:\n", utils.Prettify(template))
	return nil
}
