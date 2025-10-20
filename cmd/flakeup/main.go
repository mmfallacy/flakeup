package main

import (
	"github.com/integrii/flaggy"
)

var version = "0.0.1"

var (
	// Global flake flag
	flake string
	// Init subcommand
	initCmd  *flaggy.Subcommand
	template string
)

func init() {
	flaggy.SetName("flakeup")
	flaggy.SetDescription(`
	flakeup is a supercharged 'nix flake init -t' that allows initializing
	Nix flake projects from custom templates with advanced features like
	conflict resolution and argument substitution.
	`)
	flaggy.String(&flake, "f", "flake", `Specify the flake template source (e.g 'github:user/repo', '~/.flake') 

		    Precedence: --flake flag > $FLAKEUP_FLAKE > $FLAKE > $HOME/.nixconfig.`)
	flaggy.SetVersion(version)
	defer flaggy.Parse()

	// Init Subcommand
	initCmd = flaggy.NewSubcommand("init")
	initCmd.Description = "Initialize a new flake project from a template."
	flaggy.AttachSubcommand(initCmd, 1)

	initCmd.AddPositionalValue(&template, "template", 1, true, "Name of the template to initialize.")

}

func main() {
	if initCmd.Used {
		handleInit()
	}
}
