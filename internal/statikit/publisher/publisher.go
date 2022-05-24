package publisher

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type PublishFunc func(Args) error

type Args struct {
	Path          string // Path to directory to publish
	AccountName   string // Storage account name
	ContainerName string // Container to store to
	Key           string // Storage account access key
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

	fmt.Println(client.URL())

	return nil
}
