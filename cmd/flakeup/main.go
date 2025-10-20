package main

import (
	"fmt"
	"os"
	"path/filepath"

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

func panic(msg string) {
	fmt.Printf("Unexpected error: %s\n", msg)
	os.Exit(2)
}

// Get Flake Path from other source if flag is unset
func getFlakePath() string {
	if flake := os.Getenv("FLAKEUP_FLAKE"); flake != "" {
		return flake
	}
	if flake := os.Getenv("FLAKE"); flake != "" {
		return flake
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic("Cannot get user home dir")
	}

	return filepath.Join(home, ".nixconfig")
}

func main() {
	if flake == "" {
		flake = getFlakePath()
	}
	if initCmd.Used {
		handleInit()
	}
}
