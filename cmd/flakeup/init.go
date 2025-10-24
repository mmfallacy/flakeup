package main

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/mmfallacy/flakeup/internal/nix"
	"github.com/mmfallacy/flakeup/internal/utils"
)

func getTemplates() {
	flakeArg := fmt.Sprintf("%s#flakeupTemplates", flake)
	cmd := exec.Command("nix", "eval", flakeArg, "--json")

	out, err := cmd.Output()

	if err != nil {
		utils.Panic("flakeup: nix eval failed to run!", err)
	}

	var templates TTemplates
	if err := json.Unmarshal(out, &templates); err != nil {
		utils.Panic("flakeup: failed to unmarshal json from nix eval", err)
	}

	sout, err := utils.Prettify(templates)
	if err != nil {
		utils.Panic("", err)
	}

	fmt.Println("Unmarshalled: %s", sout)
}

func handleInit() {
	fmt.Println("Cloning %s from %s", template, flake)
	hasOutput := nix.HasFlakeOutput(flake, "flakeupTemplates")
	fmt.Println("Does my flake have the proper output? %t", hasOutput)
	getTemplates()
}
