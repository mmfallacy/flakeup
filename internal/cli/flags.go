package cli

type Flag struct {
	Desc  string
	Short string
	Full  string
}

// Subcommand: init
var DryRun = Flag{
	Desc:  "Show changes only, do not apply.",
	Short: "",
	Full:  "dry-run",
}

var NoConfirm = Flag{
	Desc:  "Apply template changes automatically",
	Short: "c",
	Full:  "no-confirm",
}

var ConflictDefault = Flag{
	Desc:  "On asks, set this as default conflict resolution",
	Short: "d",
	Full:  "conflict-default",
}
