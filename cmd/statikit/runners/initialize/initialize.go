package initialize

import (
	"os"

	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/runners"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"
)

type initStatikitProjectFunc func(fs afero.Fs, path string) error

func Runner(init initStatikitProjectFunc) runners.Runner {
	return func(fs afero.Fs, args []string, usageFor runners.UsageForFunc) error {
		if len(args) < 3 {
			usageFor(usage.Init)()
		}
		outDir := args[2]
		err := init(fs, outDir)
		if err != nil {
			// Delete outDir if it was made
			os.RemoveAll(outDir)
		}
		return err
	}
}
