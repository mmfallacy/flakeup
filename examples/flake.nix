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
      };
    };
  };
}
