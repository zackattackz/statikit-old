package runners

import "github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"

type UsageForFunc func(usage.Mode) func()

type Runner func(args []string, usageForFunc UsageForFunc) error
