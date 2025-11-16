{ pkgs, selfpkgs }:
with pkgs;
mkShell {
  name = ".nixconfig base template";

  nativeBuildInputs = [
    go
    gopls

    # For injected formatting
    nixfmt-rfc-style

    selfpkgs.flakeup
  ];

}
