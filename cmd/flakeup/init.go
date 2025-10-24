package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mmfallacy/flakeup/internal/utils"
)

func prettify(v any) (string, error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ensureFlakeupTemplatesOutputExists() bool {
	path, err := filepath.Abs(flake)

	if err != nil {
		utils.Panic("flakeup: cannot normalize flake path", err)
	}

	expr := fmt.Sprintf("(builtins.getFlake \"%s\").outputs ? flakeupTemplates", path)

	cmd := exec.Command("nix", "eval", "--impure", "--expr", expr)

	fmt.Println(cmd.Args)

	out, err := cmd.Output()

	if err != nil {
		utils.Panic("flakeup: nix eval failed to check existence of flakeupTemplates output", err)
	}

	result := strings.TrimSpace(string(out))

	return result == "true"
}

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

	sout, err := prettify(templates)
	if err != nil {
		utils.Panic("", err)
	}

	fmt.Println("Unmarshalled: %s", sout)
}

func handleInit() {
	fmt.Println("Cloning %s from %s", template, flake)
	hasOutput := ensureFlakeupTemplatesOutputExists()
	fmt.Println("Does my flake have the proper output? %t", hasOutput)
	getTemplates()
}
