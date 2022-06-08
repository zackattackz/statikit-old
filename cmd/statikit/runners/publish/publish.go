package publish

import (
	"flag"
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

var FlagSet *flag.FlagSet

var inDir, srcDir string

// Initialize FlagSet
func init() {
	const (
		// flags
		srcFlag = "src"

		// default flag values
		defaultSrc = ""

		// flag descriptions
		descSrc = "path to source statikit project"
	)

	FlagSet = flag.NewFlagSet(string(usage.Publish), flag.ExitOnError)

	FlagSet.StringVar(&srcDir, srcFlag, defaultSrc, descSrc)
}

type publishFunc func(publisher.Args) error

func Runner(publish publishFunc) runners.Runner {
	return func(fs afero.Fs, args []string, usageFor runners.UsageForFunc) error {
		FlagSet.Usage = usageFor(usage.Publish)

		FlagSet.Parse(args[2:])

		if FlagSet.NArg() > 1 ||
			srcDir == "" {
			usageFor(usage.Publish)()
		}

		// Initialize inDir to (optionally) first non-flag arg
		inDir = FlagSet.Arg(0)
		if inDir == "" {
			inDir = "."
		}

		// Clean the in/src dirs
		inDir = filepath.Clean(inDir)
		srcDir = filepath.Clean(srcDir)

		aes, err := os.ReadFile(filepath.Join(srcDir, initializer.StatikitDirName, initializer.KeyFileName))
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

		cfgParser, err := config.NewParser(fs, srcDir)
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
