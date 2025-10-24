package main

import (
	"fmt"

	"github.com/mmfallacy/flakeup/internal/nix"
	"github.com/mmfallacy/flakeup/internal/utils"
)

func handleInit() {
	fmt.Println("Cloning %s from %s", template, flake)
	hasOutput, err := nix.HasFlakeOutput(flake, "flakeupTemplates")

	if err != nil {
		utils.Panic("Something went wrong", err)
	}

	fmt.Println("Does my flake have the proper output? %t", hasOutput)

	template, err := nix.GetFlakeOutput[TTemplates](flake, "flakeupTemplates")

	fmt.Println("Got flake output:\n", utils.Prettify(template))
}
