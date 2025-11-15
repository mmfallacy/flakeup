package cli

import (
	"fmt"

	"github.com/mmfallacy/flakeup/internal/core"
	"github.com/mmfallacy/flakeup/internal/nix"
	s "github.com/mmfallacy/flakeup/internal/style"
	"github.com/mmfallacy/flakeup/internal/utils"
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
	for template, val := range templates {
		fmt.Println("============================================")
		fmt.Println(s.Success("‣ ", template))
		fmt.Println("@ ", utils.Path{Root: *val.Root, Rel: ""}.ShortenTo(8, 0))
		fmt.Print("============================================")

		out, err := s.Markdown.Render(*val.Description)
		if err != nil {
			return fmt.Errorf("show: error rendering description: %w", err)
		}

		fmt.Print(out)

		fmt.Println(s.Info("  Rules:"))
		for pattern, rule := range *val.Rules {
			fmt.Println(s.Infof("  ‣ %s : %s", pattern, *rule.OnConflict))
		}
	}

	return nil
}
