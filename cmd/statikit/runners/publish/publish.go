package publish

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/runners"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/config"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/publisher"
	"github.com/zackattackz/azure_static_site_kit/pkg/secret"
)

type publishFunc func(publisher.Args) error

func Runner(publish publishFunc) runners.Runner {
	return func(fs afero.Fs, args []string, usageFor runners.UsageForFunc) error {
		if len(args) < 3 {
			usageFor(usage.Publish)()
		}

		inDir := filepath.Clean(args[2])
		aes, err := os.ReadFile(filepath.Join(inDir, initializer.StatikitDirName, initializer.KeyFileName))
		if err != nil {
			return err
		}

		// fmt.Printf("Enter key file password: ")
		// pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
		// if err != nil {
		// 	return err
		// }
		// fmt.Println()

		pwd := []byte("aaa")

		key, err := secret.Decrypt(string(pwd), aes)
		if err != nil {
			return err
		}

		cfgParser, err := config.NewParser(fs, inDir)
		if err != nil {
			return err
		}

		cfg := config.T{}
		if err = cfgParser.Parse(&cfg); err != nil {
			return err
		}

		a := publisher.Args{Path: inDir, Key: string(key), AccountName: cfg.Az.AccountName, ContainerName: cfg.Az.ContainerName, Fs: fs, Ignore: cfg.Ignore}
		a.Ignore = append(a.Ignore, "_statikit")
		return publish(a)
	}
}
