## Features

### Core Functionality

- [x] **Command Line Interface:** `flakeup` provides a robust CLI for managing templates.
  - [x] `flakeup i[nit] <template>`: `flakeup` uses the flake specified in the following precedence: `--flake` > `$FLAKEUP_FLAKE` > `$FLAKE` > `~/.nixconfig`.
  - [x] `flakeup i[init] --flake <FLAKE> <template>`: Specify the flake template source via `--flake` flag.
  - [x] `flakeup s[how]`: Shows the list of flakeupTemplates.
- [x] **Reads Custom Flake Outputs:** Reads the `flakeupTemplates` custom flake outputs.

### `init` Subcommand

- [x] Copies files from flake to target directory.
- [x] When a conflict occurs, follow precedence rules (conflict: `"prepend"`, `"append"`, `"overwrite"`, `"ignore"`, `"ask"`).
- [x] `--dry-run`: Only show summary.
- [x] `--no-confirm`: Ask to apply template by default. When `--no-confirm` is passed, automatically apply after summary.
- [x] `-d` `--conflict-default [prepend|append|overwrite|ignore]`: Do not ask on conflicting files without rules/ with "ask" rule, automatically use passed response.
- [x] `flakeupTemplates.defaultFlags.<subcommand>`: Default flags for each subcommand.

### `show` Subcommand

- [x] Without argument, show list of template names with description.
- [x] With argument `<template>`, show rules for `<template>`.

## Planned Enhancements

- [ ] `flakeupTemplates` may specify specific arguments per template(`ARG1,ARG2,...`) with their defaults. These would substitute their values in any file of the template that contains `@@ARG1@@,@@ARG2@@,...`.
- [ ] When arbitrary flags like `--ARG1 somevalue` are passed, it will override the replacement string for ALL matching substitutes.
- [ ] `--no-substitute`: Do not ask for parameter values. Do not do substitution entirely.
- [ ] Nix package.
- [ ] Nixpkgs package.
