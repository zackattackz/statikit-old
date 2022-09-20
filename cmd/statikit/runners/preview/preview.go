package preview

import (
	"path/filepath"

	"github.com/spf13/afero"
	"codeberg.org/zackattackz/statikit/cmd/statikit/runners"
	"codeberg.org/zackattackz/statikit/cmd/statikit/usage"
)

type previewFunc func(fs afero.Fs, path string, port string) error

func Runner(preview previewFunc) runners.Runner {
	return func(fs afero.Fs, args []string, usageFor runners.UsageForFunc) error {
		if len(args) < 3 {
			usageFor(usage.Preview)()
		}
		inDir := filepath.Clean(args[2])
		return preview(fs, inDir, ":8080")
	}
}
