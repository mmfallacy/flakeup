package core

import (
	"errors"
	"fmt"
	"os"

	u "github.com/mmfallacy/flakeup/internal/utils"
)

type ActionEntry struct {
	Desc    string
	Pattern string

	Action Action
}

func (a ActionEntry) Process() error {
	return a.Action.Do()
}

func (a ActionEntry) String() string {
	return a.Action.String()
}

type Action interface {
	Do() error
	String() string
}

// Action.(Ask): Asks the user on what to do
type Ask struct {
	Src  u.Path
	Dest u.Path
}

func (a Ask) Do() error {
	return errors.New("action ask do attempt")
}

func (a Ask) Resolve(to ConflictAction) (Action, error) {
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

func (a Ask) String() string {
	return fmt.Sprintf("ask: %s", a.Dest.Shorten())
}

// Action.(Mkdir): Creates a directory
type Mkdir struct {
	Dest u.Path
}

func (a Mkdir) Do() error {
	return os.MkdirAll(a.Dest.Resolve(), 0o755)
}

func (a Mkdir) String() string {
	return fmt.Sprintf("mkdir: %s", a.Dest.Shorten())
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

func (a Exact) String() string {
	return fmt.Sprintf("copy: %s -> %s", a.Src.Shorten(), a.Dest.Shorten())
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

func (a Append) String() string {
	return fmt.Sprintf("append: %s + %s -> %s", a.Base.Shorten(), a.Suffix.Shorten(), a.Dest.Shorten())
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

func (a Prepend) String() string {
	return fmt.Sprintf("prepend: %s + %s -> %s", a.Prefix.Shorten(), a.Base.Shorten(), a.Dest.Shorten())
}

// Action.(Ignore): Noop
type Ignore struct {
	Src  u.Path
	Dest u.Path
}

func (a Ignore) Do() error {
	return nil
}

func (a Ignore) String() string {
	return fmt.Sprintf("ignore: %s", a.Dest.Shorten())
}
