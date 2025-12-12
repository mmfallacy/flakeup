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
	Action
}

type Action interface {
	Do(sub Substitutions) error
	String() string
}

// Action.(Ask): Asks the user on what to do
type Ask struct {
	Src  u.Path
	Dest u.Path
}

func (a Ask) Do(sub Substitutions) error {
	return errors.New("action ask do attempt")
}

func (a Ask) Resolve(to ConflictAction) Action {
	switch to {
	default:
		panic("Action Ask given invalid ConflictAction")
	case ConflictAppend:
		return &Append{
			Base:   a.Dest,
			Suffix: a.Src,
			Dest:   a.Dest,
		}
	case ConflictPrepend:
		return &Prepend{
			Base:   a.Dest,
			Prefix: a.Src,
			Dest:   a.Dest,
		}
	case ConflictOverwrite:
		return &Overwrite{
			Src:  a.Src,
			Dest: a.Dest,
		}
	case ConflictIgnore:
		return &Ignore{
			Src:  a.Src,
			Dest: a.Dest,
		}
	}
}

func (a Ask) String() string {
	return fmt.Sprintf("ask: %s", a.Dest.Shorten())
}

// Action.(Mkdir): Creates a directory
type Mkdir struct {
	Dest u.Path
}

func (a Mkdir) Do(sub Substitutions) error {
	return os.MkdirAll(a.Dest.Resolve(), 0o755)
}

func (a Mkdir) String() string {
	return fmt.Sprintf("mkdir: %s", a.Dest.Shorten())
}

// Action.(Exact): Copies a file from Src to Dest assuming no conflicts
type Exact struct {
	Src  u.Path
	Dest u.Path
}

func (a Exact) Do(sub Substitutions) error {
	if sub == nil {
		return MergeInto(a.Src.Resolve(), nil, a.Dest.Resolve())
	} else {
		return MergeIntoWithSubstitutions(a.Src.Resolve(), nil, a.Dest.Resolve(), sub)
	}
}

func (a Exact) String() string {
	return fmt.Sprintf("copy: %s -> %s", a.Src.Shorten(), a.Dest.Shorten())
}

// Action.(Overwrite): Copies a file from Src to Dest assuming no conflicts
type Overwrite struct {
	Src  u.Path
	Dest u.Path
}

func (a Overwrite) Do(sub Substitutions) error {
	if sub == nil {
		return MergeInto(a.Src.Resolve(), nil, a.Dest.Resolve())
	} else {
		return MergeIntoWithSubstitutions(a.Src.Resolve(), nil, a.Dest.Resolve(), sub)
	}
}

func (a Overwrite) String() string {
	return fmt.Sprintf("overwrite: %s -> %s", a.Src.Shorten(), a.Dest.Shorten())
}

// Action.(Append): Copies template from Src to Dest, appending to existing file
type Append struct {
	Base   u.Path
	Suffix u.Path
	Dest   u.Path
}

func (a Append) Do(sub Substitutions) error {
	s := a.Suffix.Resolve()
	if sub == nil {
		return MergeInto(a.Base.Resolve(), &s, a.Dest.Resolve())
	} else {
		return MergeIntoWithSubstitutions(a.Base.Resolve(), &s, a.Dest.Resolve(), sub)
	}
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

func (a Prepend) Do(sub Substitutions) error {
	b := a.Base.Resolve()

	if sub == nil {
		return MergeInto(a.Prefix.Resolve(), &b, a.Dest.Resolve())
	} else {
		return MergeIntoWithSubstitutions(a.Prefix.Resolve(), &b, a.Dest.Resolve(), sub)
	}
}

func (a Prepend) String() string {
	return fmt.Sprintf("prepend: %s + %s -> %s", a.Prefix.Shorten(), a.Base.Shorten(), a.Dest.Shorten())
}

// Action.(Ignore): Noop
type Ignore struct {
	Src  u.Path
	Dest u.Path
}

func (a Ignore) Do(sub Substitutions) error {
	return nil
}

func (a Ignore) String() string {
	return fmt.Sprintf("ignore: %s", a.Dest.Shorten())
}
