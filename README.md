# `flakeup`

`flakeup` is a supercharged `nix flake init -t`.

## Features

- [ ] command line tool
  - [ ] `flakeup i[nit] <template>`
    > `flakeup` uses the flake specified in the following precedence:
    > `--flake` > `$FLAKEUP_FLAKE` > `$FLAKE` > `~/.nixconfig`
  - [ ] `flakeup i[nit] --flake <FLAKE> <template>`
    > Specify the flake template source via `--flake` flag
- [ ] Reads the `flakeupTemplates` custom flake outputs
- [ ] Copies files from flake to current active directory
- [ ] When a conflict occurs, follow precedence rules (conflict: `"prepend"`, `"append"`, `"overwrite"`, `"ignore"`, `"ask"`)
- [ ] `flakeupTemplates` may specify specific arguments per template(`ARG1,ARG2,...`) with their defaults. These would substitute their values in any file of the template that contains `@@ARG1@@,@@ARG2@@,...`

## `outputs.flakeupTemplates` schema:

```nix
let
  types = (import <nixpkgs>).lib.types;

  templateModule = types.submodule {
    options = {
      root = lib.mkOption {
        type = types.path;
        description = "Root path of template";
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
      substitutes = lib.mkOption {
        type = types.attrsOf types.string;
        description = "Attrset where the key is a pattern to match and the value is the replacement string.";
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

Sample:

```nix
{
  flakeupTemplates = {
    someTemplate = {
      root = ./templates/base;
      rules = {
        "./somefile.txt" = {
          substitutes = {
            "@@SOMEARG@@" = "replacement";
          };
        };
        # flakeupTemplates can ask for globs instead of direct filenames.
        # all files recursively within this will follow onConflict and substitutes by default
        "./subpath-dir/*" = {
          onConflict = "ignore";
          substitutes = {
            "@@SOMEARG@@" = "replacement";
          };
        };

        # rules on files targeted by globs can be overriden
        "./subpath-dir/.gitignore" = {
          onConflict = "append";
        };

      };
    };
  };
}
```

which you can apply via `flakeup init someTemplate`
