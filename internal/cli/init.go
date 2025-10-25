package cli

import (
	"fmt"

	"github.com/mmfallacy/flakeup/internal/nix"
	"github.com/mmfallacy/flakeup/internal/utils"
)

func HandleInit(opts InitOptions) {
	fmt.Printf("Cloning template %s from flake %s\n", opts.Template, opts.FlakePath)
	hasOutput, err := nix.HasFlakeOutput(opts.FlakePath, "flakeupTemplates")

	if err != nil {
		utils.Panic("Something went wrong", err)
	}

	fmt.Println("Does my flake have the proper output?", hasOutput)

	template, err := nix.GetFlakeOutput[Templates](opts.FlakePath, "flakeupTemplates")

	fmt.Println("Got flake output:\n", utils.Prettify(template))
}
