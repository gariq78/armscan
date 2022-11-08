package installer

import (
	"fmt"
	"path/filepath"
)

func getOsDepProcName(name string) string {
	return name
}

func getOsDepCmd(name string) string {
	dir := filepath.Dir(name)
	if dir == "." {
		return fmt.Sprintf("./%s", getOsDepProcName(name))
	}
	return name
}
