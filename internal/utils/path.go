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
	return p.ShortenTo(4, 4)
}

func (path Path) ShortenTo(p, s int) string {
	parts := strings.Split(filepath.ToSlash(path.Resolve()), "/")
	for i, part := range parts[:len(parts)-1] {
		if len(part) > p+s+3 {
			parts[i] = part[:p+1] + "..." + part[len(part)-s:]
		}
	}
	return strings.Join(parts, "/")
}
