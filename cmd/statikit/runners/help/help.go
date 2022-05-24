package help

import "github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"

func Runner(args []string) error {
	if len(args) > 3 ||
		len(args) < 3 {
		usage.PrintUsageAndExit(args[0], usage.Help, nil)
	}
	m := args[2]
	if usage.IsMode(usage.Mode(m)) {
		usage.PrintUsageAndExit(args[0], usage.Mode(m), nil) // TODO use specific flagset for mode

	}
	usage.PrintUsageAndExit(args[0], usage.Help, nil)
	return nil
}
