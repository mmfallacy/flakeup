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
{
  outputs =
    { ... }:
    {
      flakeup = {
        defaultFlags = {
          init = [ "string" ];
          show = [ ];
        };

        templates = {
          "template name" = {
            description = "string"; # Or builtins.readFile textFile;
            root = ./path;
            rules = {
              "glob" = {
                onConflict = "ask"; # or "prepend" | "append" | "overwrite" | "ignore"
              };
              "glob2" = {
                onConflict = "ask"; # or "prepend" | "append" | "overwrite" | "ignore"
              };
            };
          };
        };
      };
    };
}
```
