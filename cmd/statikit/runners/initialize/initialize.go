package initialize

import (
	"os"

	"github.com/spf13/afero"
	"codeberg.org/zackattackz/statikit/cmd/statikit/runners"
	"codeberg.org/zackattackz/statikit/cmd/statikit/usage"
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
