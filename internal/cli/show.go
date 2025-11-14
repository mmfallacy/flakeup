package cli

import (
	"fmt"

	"github.com/mmfallacy/flakeup/internal/core"
	"github.com/mmfallacy/flakeup/internal/nix"
	s "github.com/mmfallacy/flakeup/internal/style"
)

type ShowOptions struct {
	GlobalOptions *GlobalOptions
	// Positionals
	Template string
}

func HandleShow(opts *ShowOptions) error {
	if hasOutput, err := nix.HasFlakeOutput(opts.GlobalOptions.FlakePath, "flakeupTemplates"); err != nil {
		return fmt.Errorf("show: %w: %w", ErrCliUnexpected, err)
	} else if !hasOutput {
		return fmt.Errorf("show: %w", ErrCliInitMissingFlakeupTemplateOutput)
	}

	conf, err := nix.GetFlakeOutput[core.Config](opts.GlobalOptions.FlakePath, "flakeupTemplates")

	if err != nil {
		return fmt.Errorf("show: %s", err)
	}

	templates := conf.Templates

	fmt.Println(s.Info("List of available flakeup templates:"))
	for template := range templates {
		fmt.Println("‣ ", template)
	}

	return nil
}
