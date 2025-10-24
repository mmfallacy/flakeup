package nix

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mmfallacy/flakeup/internal/utils"
)

func HasFlakeOutput(flake string, output string) bool {
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
