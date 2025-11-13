package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/integrii/flaggy"
	"github.com/mmfallacy/flakeup/internal/cli"
	"github.com/mmfallacy/flakeup/internal/utils"
)

var version = "0.0.1"

var (
	// Global flake flag
	flake string
	// Init subcommand
	initCmd  *flaggy.Subcommand
	template string
	outdir   string

	dryRun bool
)

var globalOpts = cli.GlobalOptions{}
var initOpts = cli.InitOptions{
	GlobalOptions: &globalOpts,
}

func init() {
	flaggy.SetName("flakeup")
	flaggy.SetDescription(`
	flakeup is a supercharged 'nix flake init -t' that allows initializing
	Nix flake projects from custom templates with advanced features like
	conflict resolution and argument substitution.
	`)
	flaggy.String(&globalOpts.FlakePath, "f", "flake", `Specify the flake template source (e.g 'github:user/repo', '~/.flake') 

		    Precedence: --flake flag > $FLAKEUP_FLAKE > $FLAKE > $HOME/.nixconfig.`)
	flaggy.SetVersion(version)
	defer flaggy.Parse()

	// Init Subcommand
	initCmd = flaggy.NewSubcommand("init")
	initCmd.Description = "Initialize a new flake project from a template."
	flaggy.AttachSubcommand(initCmd, 1)

	initCmd.AddPositionalValue(&initOpts.Template, "template", 1, true, "Name of the template to initialize.")

	outdir = "."
	initCmd.AddPositionalValue(&initOpts.OutDir, "outdir", 2, false, "Directory to put the initialized template")

	initCmd.Bool(&initOpts.DryRun, "", "dry-run", "Show changes only, do not apply.")

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
		utils.Panic("Cannot get user home dir", err)
	}

	return filepath.Join(home, ".nixconfig")
}

func main() {
	if flake == "" {
		globalOpts.FlakePath = getFlakePath()
	}

	if initCmd.Used {
		err := cli.HandleInit(initOpts)
		fmt.Println(err)
	}
}
