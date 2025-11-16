{
  pkgs,
}:
rec {
  flakeup = pkgs.buildGo125Module (finalAttrs: {
    pname = "flakeup";
    version = "0.0.1";

    src = ../.;

    vendorHash = "sha256-IeTUlKAYgviNERu7g42d79y9O4iUpYL9bXT9SgJ4Vh0=";

    meta = {
      description = "";
      homepage = "https://github.com/mmfallacy/flakeup";
      license = pkgs.lib.licenses.mit;
      maintainers = [
        {
          name = "Michael M.";
          github = "mmfallacy";
          githubId = 31348500;
        }
      ];
    };
  });
  default = flakeup;
}
