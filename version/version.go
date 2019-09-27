package version

import (
	"fmt"
	"strings"
)

type Info struct {
	Name   string
	Tag    string
	Commit string
	Branch string
}

// Version show version info
func (i *Info) Version() string {
	parts := []string{i.Name, i.Tag, i.Branch, i.Commit}

	for k, v := range parts {
		if len(v) == 0 {
			parts[k] = "unknown"
		}
	}

	gitInfo := fmt.Sprintf("(git: %s %s)", parts[2], parts[3])
	parts = append(parts[:2], gitInfo)

	return strings.Join(parts, " ")
}
