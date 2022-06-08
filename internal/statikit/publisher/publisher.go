package publisher

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/spf13/afero"
	sp "github.com/zackattackz/azure_static_site_kit/pkg/subtractPaths"
)

type Args struct {
	Path          string   // Path to directory to publish
	AccountName   string   // Storage account name
	ContainerName string   // Container to store to
	Key           string   // Storage account access key
	Fs            afero.Fs // Filesystem of Path
	Ignore        []string // Files to ignore
}

func Publish(a Args) error {
	cred, err := azblob.NewSharedKeyCredential(a.AccountName, a.Key)
	if err != nil {
		return err
	}

	client, err := azblob.NewContainerClientWithSharedKey(
		fmt.Sprintf(
			"https://%s.blob.core.windows.net/%s",
			a.AccountName,
			a.ContainerName,
		),
		cred,
		nil,
	)
	if err != nil {
		return err
	}

	pager := client.ListBlobsFlat(nil)
	for pager.NextPage(context.Background()) {
		resp := pager.PageResponse()
		for _, blob := range resp.Segment.BlobItems {
			fmt.Println(*blob.Name)
			blobClient, err := client.NewBlobClient(*blob.Name)
			if err != nil {
				return err
			}
			_, err = blobClient.Delete(context.Background(), nil)
			if err != nil {
				return err
			}
		}
	}

	err = afero.Walk(a.Fs, a.Path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = sp.SubtractPaths(a.Path, path)

		// If path is in ignore list, skip it
		for _, ignorePath := range a.Ignore {
			if match, _ := filepath.Match(ignorePath, path); match {
				if info.IsDir() {
					return fs.SkipDir
				} else {
					return nil
				}
			}
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		blockBlobClient, err := client.NewBlockBlobClient(path)
		if err != nil {
			return err
		}

		//fullPath := filepath.Join(a.Path, path)
		f, err := a.Fs.Open(filepath.Join(a.Path, path))
		if err != nil {
			return err
		}

		_, err = blockBlobClient.Upload(context.Background(), f, nil)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
