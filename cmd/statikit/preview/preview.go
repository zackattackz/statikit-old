package preview

import "github.com/zackattackz/azure_static_site_kit/internal/statikit/previewer"

type Args struct {
	path string
	port string
}

func Run(p previewer.Previewer) error {
	return p.Preview()
}
