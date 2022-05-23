package publisher

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type Interface interface {
	Upload() error
}

type Args struct {
	Path          string // Path to directory to publish
	AccountName   string // Storage account name
	ContainerName string // Container to store to
	Key           string // Storage account access key
}

type uploader struct {
	Args
}

func NewUploader(a Args) Interface {
	return &uploader{a}
}

func (u *uploader) Upload() error {
	cred, err := azblob.NewSharedKeyCredential(u.AccountName, u.Key)
	if err != nil {
		return err
	}

	client, err := azblob.NewContainerClientWithSharedKey(
		fmt.Sprintf(
			"https://%s.blob.core.windows.net/%s",
			u.AccountName,
			u.ContainerName,
		),
		cred,
		nil,
	)
	if err != nil {
		return err
	}

	fmt.Println(client.URL())

	return nil
}
