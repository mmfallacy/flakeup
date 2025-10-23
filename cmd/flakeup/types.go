package main

type TConflictAction string

type TParameter struct {
	Name *string `json:"name"`
	// Nullable
	Prompt  *string `json:"prompt"`
	Default *string `json:"default"`
}

type TRule struct {
	OnConflict *TConflictAction
}

type TRules map[string]TRule

type TTemplate struct {
	Root       *string       `json:"root"`
	Parameters *[]TParameter `json:"parameters"`
	Rules      *TRules       `json:"rules"`
}

type TTemplates map[string]TTemplate
