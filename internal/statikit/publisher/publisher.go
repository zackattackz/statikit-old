package publisher

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type Publisher interface {
	Publish() error
}

type Args struct {
	Path          string // Path to directory to publish
	AccountName   string // Storage account name
	ContainerName string // Container to store to
	Key           string // Storage account access key
}

type publisher struct {
	Args
}

func New(a Args) Publisher {
	return &publisher{a}
}

func (p *publisher) Publish() error {
	cred, err := azblob.NewSharedKeyCredential(p.AccountName, p.Key)
	if err != nil {
		return err
	}

	client, err := azblob.NewContainerClientWithSharedKey(
		fmt.Sprintf(
			"https://%s.blob.core.windows.net/%s",
			p.AccountName,
			p.ContainerName,
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
