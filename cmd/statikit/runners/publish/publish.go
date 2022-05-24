package publish

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/runners"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/publisher"
	"github.com/zackattackz/azure_static_site_kit/pkg/secret"
	"golang.org/x/term"
)

func Runner(publish publisher.PublishFunc) runners.Runner {
	return func(args []string, usageFor runners.UsageForFunc) error {
		if len(args) < 3 {
			usageFor(usage.Publish)()
		}

		inDir := filepath.Clean(args[2])
		aes, err := os.ReadFile(filepath.Join(inDir, initializer.KeyFileName))
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

		a := publisher.Args{Path: inDir, Key: string(key)}
		return publish(a)
	}
}
