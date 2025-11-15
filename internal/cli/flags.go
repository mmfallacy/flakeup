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

// Subcommand: show
var ShowSource = Flag{
	Desc:  "Include template sources in output",
	Short: "s",
	Full:  "source",
}

var ShowDesc = Flag{
	Desc:  "Include template descriptions in output",
	Short: "d",
	Full:  "description",
}

var ShowRules = Flag{
	Desc:  "Include template rules in output",
	Short: "r",
	Full:  "rules",
}
