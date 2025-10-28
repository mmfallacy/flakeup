package core

// JSON schema as structs
type Templates map[string]Template

type Template struct {
	Root       *string      `json:"root"`
	Parameters *[]Parameter `json:"parameters"`
	Rules      *Rules       `json:"rules"`
}

type Parameter struct {
	Name *string `json:"name"`
	// Nullable
	Prompt  *string `json:"prompt"`
	Default *string `json:"default"`
}

type Rules map[string]Rule

type Rule struct {
	OnConflict *ConflictAction
}

type ConflictAction string

const (
	ConflictPrepend   ConflictAction = "prepend"
	ConflictAppend    ConflictAction = "append"
	ConflictOverwrite ConflictAction = "overwrite"
	ConflictIgnore    ConflictAction = "ignore"
	ConflictAsk       ConflictAction = "ask"
)

type Action interface {
	Kind() string
}

type ActionApply struct {
	Desc    string
	Src     string
	Dest    string
	Pattern string
	Rule    Rule
	Write   bool
}

func (a ActionApply) Kind() string { return "apply" }

type ActionAsk struct {
	Desc    string
	Src     string
	Dest    string
	Pattern string
	Rule    Rule
	Default string
}

func (a ActionAsk) Kind() string { return "ask" }

type ActionMkdir struct {
	Desc string
	Dest string
}

func (a ActionMkdir) Kind() string { return "mkdir" }
