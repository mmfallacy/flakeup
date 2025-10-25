package cli

// ConflictAction enum
const (
	ConflictPrepend   ConflictAction = "prepend"
	ConflictAppend    ConflictAction = "append"
	ConflictOverwrite ConflictAction = "overwrite"
	ConflictIgnore    ConflictAction = "ignore"
	ConflictAsk       ConflictAction = "ask"
)
