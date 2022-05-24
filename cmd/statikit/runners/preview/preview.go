package preview

import (
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/runners"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/previewer"
)

func Runner(preview previewer.PreviewFunc) runners.Runner {
	return func(args []string, usageFor runners.UsageForFunc) error {
		if len(args) < 3 {
			usageFor(usage.Preview)()
		}

		inDir := filepath.Clean(args[2])
		a := previewer.Args{Path: inDir, Port: ":8080"}
		return preview(a)
	}
}