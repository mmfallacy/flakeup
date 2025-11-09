package core

import (
	"errors"
	"os"

	u "github.com/mmfallacy/flakeup/internal/utils"
)

type ActionEntry struct {
	Desc    string
	Pattern string

	Action _Action
}

func (a ActionEntry) Process() error {
	return a.Action.Do()
}

type _Action interface {
	Do() error
}

// Action.(Ask): Asks the user on what to do
type Ask struct {
	Src  u.Path
	Dest u.Path
}

func (a Ask) Do() error {
	return ErrActionAskProcessAttempt
}

func (a Ask) Resolve(to ConflictAction) (_Action, error) {
	switch to {
	default:
		return nil, errors.New("action ask unhandled conflict action")
	case ConflictAppend:
		return &Append{
			Base:   a.Dest,
			Suffix: a.Src,
			Dest:   a.Dest,
		}, nil
	case ConflictPrepend:
		return &Prepend{
			Base:   a.Dest,
			Prefix: a.Src,
			Dest:   a.Dest,
		}, nil
	case ConflictOverwrite:
		return &Exact{
			Src:  a.Src,
			Dest: a.Dest,
		}, nil
	case ConflictIgnore:
		return &Ignore{
			Src:  a.Src,
			Dest: a.Dest,
		}, nil
	}
}

// Action.(Mkdir): Creates a directory
type Mkdir struct {
	Dest u.Path
}

func (a Mkdir) Do() error {
	return os.Mkdir(a.Dest.Resolve(), 0o755)
}

// Action.(Exact): Copies a file from Src to Dest assuming no conflicts
// We will hoist Overwrite behavior here as we initially write to a tempdir
type Exact struct {
	Src  u.Path
	Dest u.Path
}

func (a Exact) Do() error {
	return MergeInto(a.Src.Resolve(), nil, a.Dest.Resolve())
}

// TODO: Should I merge Append and Prepend to combined type?
// Action.(Append): Copies template from Src to Dest, appending to existing file
type Append struct {
	Base   u.Path
	Suffix u.Path
	Dest   u.Path
}

func (a Append) Do() error {
	s := a.Suffix.Resolve()
	return MergeInto(a.Base.Resolve(), &s, a.Dest.Resolve())
}

// Action.(Prepend): Copies template from Src to Dest, prepending to existing file
type Prepend struct {
	Base   u.Path
	Prefix u.Path
	Dest   u.Path
}

func (a Prepend) Do() error {
	p := a.Prefix.Resolve()
	return MergeInto(a.Base.Resolve(), &p, a.Dest.Resolve())
}

// Action.(Ignore): Noop
type Ignore struct {
	Src  u.Path
	Dest u.Path
}

func (a Ignore) Do() error {
	return nil
}
