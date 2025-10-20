package main

import (
	"fmt"
	"github.com/integrii/flaggy"
)

func handleInit() {
	fmt.Println("Hello from subcommand init")
}

func main() {
	flaggy.SetName("flakeup")
	flaggy.SetDescription(`
	flakeup is a supercharged 'nix flake init -t' that allows initializing
	Nix flake projects from custom templates with advanced features like
	conflict resolution and argument substitution.
	`)

	init := flaggy.NewSubcommand("init")
	init.Description = "Initialize a new flake project from a template."

	var flake string
	flaggy.String(&flake, "f", "flake", `Specify the flake template source (e.g 'github:user/repo', '~/.flake') 
		    Precedence: --flake flag > $FLAKEUP_FLAKE > $FLAKE > $HOME/.nixconfig.`)

	flaggy.AttachSubcommand(init, 1)
	flaggy.Parse()
}
