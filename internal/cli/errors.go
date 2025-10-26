package cli

import "errors"

var ErrCliUnexpected = errors.New("unexpected error while running cli")

var ErrCliInitMissingFlakeupTemplateOutput = errors.New("flakeupTemplates ouptut is missing")
