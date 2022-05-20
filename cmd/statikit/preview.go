package main

import "github.com/zackattackz/azure_static_site_kit/pkg/statikit/previewer"

type previewArgs struct {
	path string
}

func preview(a previewArgs) error {
	return previewer.Preview(a.path)

}
