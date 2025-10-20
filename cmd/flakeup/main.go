package main

import (
	// "flag"
	"fmt"
	"os"
)

var usage string = `
Usage: flakeup [command]

Commands:
  init, i <template>    Initialize a new flake project from a template.

Flags:
  --flake <FLAKE>       Specify the flake template source (e.g., 'github:user/repo').
                        Precedence: $FLAKEUP_FLAKE > $FLAKE > ~/.nixconfig > --flake flag.

Description:
  flakeup is a supercharged 'nix flake init -t' that allows initializing
  Nix flake projects from custom templates with advanced features like
  conflict resolution and argument substitution.
`

func handleInit() {
	fmt.Println("Hello from subcommand init")
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("flakeup: Missing subcommand [init].")
		fmt.Println(usage)
		os.Exit(1)
	}

	subcmd := os.Args[1]

	switch subcmd {
	case "init":
		handleInit()
	default:
		fmt.Printf("flakeup: Unknown subcommand %s\n", subcmd)
		fmt.Println(usage)
	}

}
