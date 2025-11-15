{
  description = "Flakeup example flake with flakeupTemplates";

  inputs = { };

  outputs = inputs: {
    flakeupTemplates = {
      defaultFlags = {
        init = [
          "--no-confirm"
          "-d i"
        ];
      };
      templates = {
        template = {
          description = builtins.readFile ./template/README.md;
          root = ./template;
          rules = {
            "nix/*" = {
              onConflict = "ask";
            };
            ".envrc" = {
              onConflict = "ignore";
            };
          };

          parameters = [
            {
              name = "ARG1";
              prompt = "Specify Argument 1";
            }
          ];
        };
        template2 = {
          description = builtins.readFile ./template2/README.md;
          root = ./template2;
        };
      };
    };
  };
}
