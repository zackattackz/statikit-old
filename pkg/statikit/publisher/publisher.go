package publisher

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/secret"
	"golang.org/x/term"
)

type PublisherArgs struct {
	Path        string // Path to directory to publish
	AccountName string // Storage account name
	Key         string // Storage account access key
}

func Publish(a PublisherArgs) error {
	const keyFileName = "key.aes256"

	aes, err := os.ReadFile(filepath.Join(a.Path, keyFileName))
	if err != nil {
		return err
	}

	fmt.Printf("Enter key file password: ")
	pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()

	key, err := secret.Decrypt(string(pwd), aes)
	if err != nil {
		return err
	}

	fmt.Println(string(key))

	// cred, err := azblob.NewSharedKeyCredential(a.AccountName, string(key))
	// if err != nil {
	// 	return err
	// }

	// client, err := azblob.NewContainerClient()

	return nil
}
