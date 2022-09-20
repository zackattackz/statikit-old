package runners

import (
	"github.com/spf13/afero"
	"codeberg.org/zackattackz/statikit/cmd/statikit/usage"
)

type UsageForFunc func(usage.Mode) func()

type Runner func(fs afero.Fs, args []string, usageForFunc UsageForFunc) error
