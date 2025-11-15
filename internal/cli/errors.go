package cli

import "errors"

var ErrCliUnexpected = errors.New("unexpected error while running cli")

var ErrCliInitMissingFlakeupOutput = errors.New("flakeup ouptut is missing")
