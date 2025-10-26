package nix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
)

var ErrNixEvalFailed = errors.New("nix eval failed")
var ErrNixJsonUnmarshallFailed = errors.New("nix eval resulted in an unexpected json schema")

func HasFlakeOutput(flake, output string) (bool, error) {
	path, err := filepath.Abs(flake)
	if err != nil {
		return false, fmt.Errorf("flake path resolution failed: %w", err)
	}
	expr := fmt.Sprintf("(builtins.getFlake \"%s\").outputs ? %s", path, output)
	cmd := exec.Command("nix", "eval", "--impure", "--expr", expr)
	out, err := cmd.Output()

	if err != nil {
		return false, fmt.Errorf("%w: cannot check if %s output exists in flake", ErrNixEvalFailed, output)
	}

	return bytes.Equal(bytes.TrimSpace(out), []byte("true")), nil
}

// T should be a struct that would contain the JSON contents
func GetFlakeOutput[T any](flake, output string) (T, error) {
	var unmarshalled T

	flakeArg := fmt.Sprintf("%s#%s", flake, output)
	cmd := exec.Command("nix", "eval", flakeArg, "--json")

	out, err := cmd.Output()

	if err != nil {
		return unmarshalled, fmt.Errorf("%w: %w", ErrNixEvalFailed, err)
	}

	if err := json.Unmarshal(out, &unmarshalled); err != nil {
		return unmarshalled, fmt.Errorf("%w: %w", ErrNixJsonUnmarshallFailed, err)
	}

	return unmarshalled, nil
}
