package utils

import "path/filepath"

type Path struct {
	Root string
	Rel  string
}

func (p Path) Resolve() string {
	return filepath.Join(p.Root, p.Rel)
}
