package preview

import (
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/runners"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"
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
