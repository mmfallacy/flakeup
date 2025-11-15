{
  description = "Flakeup base example template 2";

  outputs =
    { nixpkgs-stable, systems, ... }@inputs:
    builtins.foldl' (a: b: a // b) { } (
      builtins.map (
        system:
        let
          pkgs = nixpkgs-stable.legacyPackages.${system};
          extras = {
            pkgs-unstable = inputs.nixpkgs-unstable.legacyPackages.${system};
          };
        in
        {
        }
      ) (import systems)
    );

  inputs = {
    nixpkgs-stable.url = "github:nixos/nixpkgs/nixos-25.05";
    nixpkgs-unstable.url = "github:nixos/nixpkgs/nixos-unstable";
    systems.url = "github:nix-systems/default";
  };
}
