package help

import (
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/runners"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"
)

func Runner(args []string, usageFor runners.UsageForFunc) error {
	if len(args) > 3 ||
		len(args) < 3 {
		usageFor(usage.Help)()
	}
	m := args[2]
	if usage.IsMode(usage.Mode(m)) {
		usageFor(usage.Mode(m))()
	}
	usageFor(usage.Help)()
	return nil
}
