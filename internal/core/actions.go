package core

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
