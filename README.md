# `flakeup`

`flakeup` is a supercharged `nix flake init -t`.

## Features

- [x] command line tool

  - [x] `flakeup i[nit] <template>`
    > `flakeup` uses the flake specified in the following precedence:
    > `--flake` > `$FLAKEUP_FLAKE` > `$FLAKE` > `~/.nixconfig`
  - [x] `flakeup i[nit] --flake <FLAKE> <template>`
    > Specify the flake template source via `--flake` flag
  - [ ] `flakeup s[how]`
    > Shows the list of flakeupTemplates

- [x] Reads the `flakeupTemplates` custom flake outputs

### `init`

- [x] Copies files from flake to target directory
- [x] When a conflict occurs, follow precedence rules (conflict: `"prepend"`, `"append"`, `"overwrite"`, `"ignore"`, `"ask"`)
- [ ] `flakeupTemplates` may specify specific arguments per template(`ARG1,ARG2,...`) with their defaults. These would substitute their values in any file of the template that contains `@@ARG1@@,@@ARG2@@,...`
- [ ] when arbitrary flags like `--ARG1 somevalue` are passed, it will override the replacement string for ALL matching substitutes.
- [x] `--dry-run`: only show summary
- [x] `--no-confirm`: ask to apply template by default. when `--no-confirm` is passed, automatically apply after summary
- [ ] `--no-conflict-ask [prepend|append|overwrite|ignore]`: do not ask on conflicting files without rules/ with "ask" rule, automatically use passed response
- [ ] `--no-substitute`: do not ask for parameter values. do not do substitution entirely.
- [ ] `flakeupTemplates.<templates>.init.defaultFlags`. default flags for init
- [ ] `flakeupTemplates.<templates>.defaultFlags`.for all
- [ ] nix package
- [ ] nixpkgs package

### `show`

- [ ] without argument, show list of template names with description
- [ ] with argument <template>, show rules for <template>

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
