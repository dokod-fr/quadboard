package quadlet

import (
	"os"
	"path/filepath"
	"strings"
)

func isDropIn(path string) bool {
	return strings.Contains(path, ".container.d")
}

func isTemplate(path string) bool {
	return strings.Contains(filepath.Base(path), "@.")
}

func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeSymlink != 0
}
