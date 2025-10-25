package cli

// Command-line Options
type GlobalOptions struct {
	FlakePath string
}

type InitOptions struct {
	GlobalOptions
	Template string
}

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
