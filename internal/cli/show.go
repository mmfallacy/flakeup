package cli

import (
	"fmt"
	"strings"

	"github.com/mmfallacy/flakeup/internal/core"
	"github.com/mmfallacy/flakeup/internal/nix"
	s "github.com/mmfallacy/flakeup/internal/style"
	u "github.com/mmfallacy/flakeup/internal/utils"
)

type ShowOptions struct {
	GlobalOptions *GlobalOptions
	// Positionals
	Template string
	// Flags
	Source bool
	Desc   bool
	Rules  bool
}

var HR = strings.Repeat("=", 80)

func HandleShow(opts *ShowOptions) error {
	if hasOutput, err := nix.HasFlakeOutput(opts.GlobalOptions.FlakePath, "flakeup"); err != nil {
		return fmt.Errorf("show: %w: %w", ErrCliUnexpected, err)
	} else if !hasOutput {
		return fmt.Errorf("show: %w", ErrCliInitMissingFlakeupOutput)
	}

	conf, err := nix.GetFlakeOutput[core.Config](opts.GlobalOptions.FlakePath, "flakeup")

	if err != nil {
		return fmt.Errorf("show: %s", err)
	}

	templates := conf.Templates

	// By default, show Source Desc Rules when given specific template
	if opts.Template != "" {
		val, ok := templates[opts.Template]
		if !ok {
			fmt.Println(s.Err(s.Icons.Err, " Cannot find specified template with name ", opts.Template))
		}
		return showTemplate(opts.Template, val, &ShowOptions{
			Source: true,
			Desc:   true,
			Rules:  true,
		})
	}

	fmt.Println(s.Info("List of available flakeup templates:"))

	// Short print
	if !opts.Desc && !opts.Rules {
		return showShort(templates, opts)
	}

	// Full Show
	return showFull(templates, opts)
}

func showShort(templates core.Templates, opts *ShowOptions) error {
	for template, val := range templates {
		fmt.Print(s.Success("‣ ", template))
		if opts.Source {
			source := u.UnwrapOrDefault(val.Root, "nil")
			fmt.Print(" @ ", u.Path{Root: source, Rel: ""}.ShortenTo(8, 0))
		}
		fmt.Println()
	}
	return nil
}

func showFull(templates core.Templates, opts *ShowOptions) error {
	for template, val := range templates {
		if err := showTemplate(template, val, opts); err != nil {
			return err
		}
	}
	return nil
}

func showTemplate(template string, val core.Template, opts *ShowOptions) error {
	fmt.Println(HR)
	fmt.Println(s.Success("‣ ", template))
	if opts.Source {
		source := u.UnwrapOrDefault(val.Root, "nil")
		fmt.Println(" @ ", u.Path{Root: source, Rel: ""}.ShortenTo(8, 0))
	}
	fmt.Print(HR)

	if opts.Desc && val.Description != nil {
		out, err := s.Markdown.Render(*val.Description)
		if err != nil {
			return fmt.Errorf("show: error rendering description: %w", err)
		}

		fmt.Print(out)
	}

	if opts.Rules && val.Rules != nil {
		fmt.Println(s.Info("  Rules:"))
		for pattern, rule := range *val.Rules {
			fmt.Println(s.Infof("  ‣ %s : %s", pattern, *rule.OnConflict))
		}
	}

	return nil
}
