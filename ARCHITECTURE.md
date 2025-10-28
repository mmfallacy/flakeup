# `flakeup` - Nix Flake Project Initializer

`flakeup` is a CLI tool that simplifies creating Nix flake projects from templates. It enhances `nix flake init -t` with advanced features for template management, conflict resolution, and dynamic content generation.

## Key Features

- **Flexible Template Sources**: Get templates from flags, environment variables, or a default config.
- **Custom Flake Outputs**: Uses `flakeupTemplates` defined in Nix flakes for rich template definitions.
- **Smart File Copying**: Copies template files to your project directory.
- **Conflict Resolution**: Handles file conflicts with `prepend`, `append`, `overwrite`, `ignore`, or `ask` options.
- **Dynamic Content**: Substitutes template parameters (e.g., `@@ARG@@`) with user input or default values, and allows command-line overrides.

## How `flakeup` Works (Code Context)

`flakeup` is a Go application that leverages the `flaggy` library for command-line argument parsing. Its core functionality revolves around the `init` subcommand, which orchestrates the template initialization process.

1.  **CLI Setup (`cmd/flakeup/main.go`):**

    - The `main.go` file defines the CLI structure using `flaggy`, setting up the `flakeup` command and its `init` subcommand.
    - It handles global flags like `--flake` and parses positional arguments such as `template` name and `outdir`.
    - Environment variables (`FLAKEUP_FLAKE`, `FLAKE`) and a default path (`~/.nixconfig`) are checked to determine the flake source if the `--flake` flag is not explicitly set.

2.  **Global Options & Command Handling (`internal/cli/init.go`, `internal/cli/errors.go`):**

    - Global options, including the resolved `FlakePath`, are encapsulated in a `GlobalOptions` struct.
    - The `HandleInit` function (within `internal/cli/init.go`) is the entry point for the `init` subcommand's logic. It receives `InitOptions` which include `GlobalOptions`, the `Template` name, and the `OutDir`.
    - Error handling is centralized in `internal/cli/errors.go` to provide consistent and user-friendly error messages.

3.  **Core Logic (`internal/core/`):**

    - The `internal/core` package is intended to house the main business logic, independent of the CLI. This includes:
      - **`template.go`**: Responsible for reading and parsing the `flakeupTemplates` output from the specified Nix flake. It understands the schema for template definitions, including parameters and conflict rules.
      - **`copy.go`**: Handles the actual file copying from the template's `root` path to the target `outdir`. It applies the defined `onConflict` rules during this process.
      - **`types.go`**: Defines the core data structures used across the application, such as `Template`, `Parameter`, and `Rule`.

4.  **Nix Interaction (`internal/nix/nix.go`):**

    - The `internal/nix` package is dedicated to interacting with the Nix ecosystem. This would involve functions to:
      - Evaluate Nix flakes to extract the `flakeupTemplates` output.
      - Potentially execute Nix commands to fetch or build flake inputs.

5.  **Utilities (`internal/utils/`):**
    - The `internal/utils` package provides general utility functions, such as `print_helpers.go` for formatted output and `utils.go` for other common helper functions.

## Quick Start

1.  **Clone & Setup:**
    ```bash
    git clone https://github.com/mmfallacy/flakeup.git
    cd flakeup
    nix develop
    ```
2.  **Initialize a Project:**

    ```bash
    # Basic usage
    flakeup init <template-name> [output-directory]

    # With custom flake source and parameter override
    flakeup init my-go-app --flake github:user/repo --PROJECT_NAME "MyService"
    ```

## Template Schema (`outputs.flakeupTemplates`)

Templates are defined in your Nix flake's `flakeupTemplates` output, following this structure:

```nix
{
  flakeupTemplates = {
    "template-name" = {
      root = ./path/to/template; # Path to template files
      parameters = [
        {
          name = "ARG1";
          prompt = "Enter value for ARG1:";
          default = "default-value";
        }
      ];
      rules = {
        "**/*.nix" = {
          onConflict = "overwrite";
        }; # Conflict resolution rules
      };
    };
  };
}
```

This schema allows defining template roots, customizable parameters, and conflict resolution rules for different file patterns.
