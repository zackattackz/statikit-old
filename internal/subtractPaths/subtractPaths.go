package subtractpaths

import (
	"path/filepath"
	"strings"
)

func SubtractPaths(parent, child string) string {
	parentList := strings.Split(parent, string(filepath.Separator))
	childList := strings.Split(child, string(filepath.Separator))

	return filepath.Join(childList[len(parentList):]...)
}
