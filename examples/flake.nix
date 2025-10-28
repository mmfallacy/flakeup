{
  description = "Flakeup example flake with flakeupTemplates";

  inputs = { };

  outputs = inputs: {
    flakeupTemplates = {
      template = {
        root = ./template;
        rules = {
          "nix/*" = {
            onConflict = "ask";
          };
          "*" = {
            onConflict = "overwrite";
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
}
