package updater

import (
	"fmt"
	"strings"
)

func getOsDepProcName(name string) string {
	return fmt.Sprintf("%s.exe", strings.TrimSuffix(name, ".exe"))
}

func getOsDepCmd(name string) string {
	return fmt.Sprintf("%s", getOsDepProcName(name))
}
