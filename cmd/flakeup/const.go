package main

const (
	ConflictPrepend TConflictAction = iota
	ConflictAppend
	ConflictOverwrite
	ConflictIgnore
	ConflictAsk
)
