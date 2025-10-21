package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

func prettify(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func getTemplates() {
	flakeArg := fmt.Sprintf("%s#flakeupTemplates", flake)
	cmd := exec.Command("nix", "eval", flakeArg, "--json")

	out, err := cmd.Output()

	if err != nil {
		panic("flakeup: nix eval failed to run!")
	}

	var templates TTemplates
	if err := json.Unmarshal(out, &templates); err != nil {
		panic("flakeup: nix eval failed to run!")
	}

	sout, err := prettify(templates)
	if err != nil {
		panic("")
	}
	fmt.Println("Unmarshalled: %s", sout)
}

func handleInit() {
	fmt.Printf("Cloning %s from %s\n", template, flake)
	getTemplates()
}
