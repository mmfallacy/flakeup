{
  description = "Flakeup example flake with flakeupTemplates";

  inputs = { };

  outputs = inputs: {
    flakeupTemplates = {
      template = {
        root = ./template;
        rules = {
          "**/*" = {
            onConflict = "ask";
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
