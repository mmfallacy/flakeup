# `flakeup`

`flakeup` is a supercharged `nix flake init -t`.

## Features

- [ ] command line tool
  - [ ] `flakeup i[nit] <template>`
    > `flakeup` uses the flake specified in the following precedence:
    > `$FLAKEUP_FLAKE` > `$FLAKE` > `~/.nixconfig`
  - [ ] `flakeup i[nit] --flake <FLAKE> <template>`
    > Specify the flake template source via `--flake` flag
- [ ] Reads the `flakeupTemplates` custom flake outputs
- [ ] Copies files from flake to current active directory
- [ ] When a conflict occurs, follow precedence rules (conflict: `"prepend"`, `"append"`, `"overwrite"`, `"ignore"`, `"ask"`)
- [ ] `flakeupTemplates` may specify specific arguments per template(`ARG1,ARG2,...`) with their defaults. These would substitute their values in any file of the template that contains `@@ARG1@@,@@ARG2@@,...`
