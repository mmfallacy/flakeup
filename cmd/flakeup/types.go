package main

type TConflictAction int

type TRule struct {
	OnConflict TConflictAction
}

type TParameter struct {
	Name string `json:"name"`
	// Nullable
	Prompt  *string `json:"prompt"`
	Default *string `json:"default"`
}

type TTemplate struct {
	Root       string       `json:"root"`
	Parameters []TParameter `json:"parameters"`
}

type TTemplates map[string]TTemplate
