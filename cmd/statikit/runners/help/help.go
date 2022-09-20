package help

import (
	"github.com/spf13/afero"
	"codeberg.org/zackattackz/statikit/cmd/statikit/runners"
	"codeberg.org/zackattackz/statikit/cmd/statikit/usage"
)

func Runner(_ afero.Fs, args []string, usageFor runners.UsageForFunc) error {
	if len(args) > 3 ||
		len(args) < 3 {
		usageFor(usage.Help)()
	}
	m := args[2]
	if usage.Mode(m).IsValid() {
		usageFor(usage.Mode(m))()
	}
	usageFor(usage.Help)()
	return nil
}
