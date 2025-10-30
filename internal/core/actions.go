package core

import (
	"errors"
	"fmt"

	"github.com/mmfallacy/flakeup/internal/utils"
)

type ActionKind string

var (
	ActionKindApply ActionKind = "apply"
	ActionKindAsk   ActionKind = "ask"
	ActionKindMkdir ActionKind = "mkdir"
)

type Action interface {
	Kind() ActionKind
	Process() error
}

type ActionApply struct {
	Desc    string
	Src     string
	Dest    string
	Pattern string
	Rule    Rule
	Write   bool
}

func (a ActionApply) Kind() ActionKind { return ActionKindApply }
func (a ActionApply) Process() error {
	fmt.Println("Processing action apply\n", utils.Prettify(a))
	return nil
}

type ActionAsk struct {
	Desc    string
	Src     string
	Dest    string
	Pattern string
	Rule    Rule
	Default string
}

func (a ActionAsk) Kind() ActionKind { return ActionKindAsk }

var ErrActionAskProcessAttempt = errors.New("this ActionAsk should have been realised into an ActionApply first before processing")

func (a ActionAsk) Process() error {
	return ErrActionAskProcessAttempt
}

type ActionMkdir struct {
	Desc string
	Dest string
}

func (a ActionMkdir) Kind() ActionKind { return ActionKindMkdir }

func (a ActionMkdir) Process() error {
	fmt.Println("Processing action mkdir\n", utils.Prettify(a))
	return nil
}
