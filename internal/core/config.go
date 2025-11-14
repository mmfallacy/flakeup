package core

// JSON schema as structs
type Config struct {
	DefaultFlags *DefaultFlags `json:"defaultFlags"`
	Templates    Templates     `json:"templates"`
}
type DefaultFlags struct {
	Init []string
}

type Templates map[string]Template

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
