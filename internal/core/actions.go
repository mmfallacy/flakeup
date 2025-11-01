package core

import (
	"errors"
	"os"
	"path/filepath"
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
	Path    string
	Pattern string
	Rule    Rule
	Write   bool
}

func (a ActionApply) Kind() ActionKind { return ActionKindApply }
func (a ActionApply) Process() error {
	rootpath := filepath.Join(a.Src, a.Path)
	outpath := filepath.Join(a.Dest, a.Path)
	// return nil
	return Copy(rootpath, outpath)
}

type ActionAsk struct {
	Desc    string
	Src     string
	Dest    string
	Path    string
	Pattern string
	Rule    Rule
}

func (a ActionAsk) Kind() ActionKind { return ActionKindAsk }

var ErrActionAskProcessAttempt = errors.New("this ActionAsk should have been realised into an ActionApply first before processing")

func (a ActionAsk) Process() error {
	return ErrActionAskProcessAttempt
}

type ActionMkdir struct {
	Desc string
	Dest string
	Path string
}

func (a ActionMkdir) Kind() ActionKind { return ActionKindMkdir }

func (a ActionMkdir) Process() error {
	outpath := filepath.Join(a.Dest, a.Path)
	return os.Mkdir(outpath, 0o755)
}
