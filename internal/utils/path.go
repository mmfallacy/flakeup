package utils

import (
	"path/filepath"
	"strings"
)

type Path struct {
	Root string
	Rel  string
}

func (p Path) Resolve() string {
	return filepath.Join(p.Root, p.Rel)
}

func (p Path) Shorten() string {
	parts := strings.Split(filepath.ToSlash(p.Resolve()), "/")
	for i, part := range parts[:len(parts)-1] {
		if len(part) > 11 {
			parts[i] = part[:4] + "..." + part[len(part)-4:]
		}
	}
	return strings.Join(parts, "/")
}
