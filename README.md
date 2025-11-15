# `flakeup`

`flakeup` is a supercharged `nix flake init -t`, designed to provide a more flexible and powerful way to initialize Nix flakes from templates.

## Features

`flakeup` provides a robust command-line interface for managing Nix flake templates. Key features include:

- **Template Initialization:** Initialize new projects from specified templates with flexible source precedence (flags, environment variables, or `~/.nixconfig`).
- **Conflict Resolution:** Configurable handling of file conflicts during template application (prepend, append, overwrite, ignore, or ask).
- **Dry Run & No Confirmation Modes:** Safely preview changes or automate template application.
- **Default Conflict Actions:** Set default behaviors for conflict resolution.
- **Default Flags:** Define subcommand-specific default flags within `flakeupTemplates`.
- **Template Discovery:** List available templates and view detailed rules for specific templates.
- **Custom Flake Outputs:** Utilizes `flakeupTemplates` custom flake outputs for template discovery.

For a detailed list of features and planned enhancements, please refer to [TODO.md](./TODO.md).

## `outputs.flakeupTemplates` schema:

```nix
let
  types = (import <nixpkgs>).lib.types;

  parameterType = types.submodule {
    options = {
      name = lib.mkOption {
        type = types.string;
        description = "Parameter name";
      };
      prompt = lib.mkOption {
        type = types.nullOr types.string;
        default = null;
        description = "Prompt to use to ask user for replacement value";
      };
      default = lib.mkOption {
        type = types.nullOr types.string;
        default = null;
        description = "Default value for parameter";
      };

    };
  };

  templateModule = types.submodule {
    options = {
      root = lib.mkOption {
        type = types.path;
        description = "Root path of template";
      };

      parameters = lib.mkOption {
        type = types.listOf parameterType;
        description = "List of valid parameters flakeup will process.";
      };

      rules = lib.mkOption {
        type = types.attrsOf ruleModule;
        description = "Attrset where the key is a glob pattern and the values are rules";
      };
    };
  };

  ruleModule = types.submodule {
    options = {
      onConflict = lib.mkOption {
        type = types.enum [
          "prepend"
          "append"
          "overwrite"
          "ignore"
          "ask"
        ];
        description = "Determines how flakeup handles file conflicts";
      };
    };
  };

in
{
  flakeupTemplates = lib.mkOption {
    type = types.attrsOf templateModule;
    description = "Attrset where the key is a template name and the values are template options";
  };
}
```
