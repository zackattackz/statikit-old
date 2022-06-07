package initialize

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/runners"
	"github.com/zackattackz/azure_static_site_kit/cmd/statikit/usage"
	"golang.org/x/term"
)

type initStatikitProjectFunc func(fs afero.Fs, path string, pwd string, key []byte) error

func Runner(init initStatikitProjectFunc) runners.Runner {
	return func(fs afero.Fs, args []string, usageFor runners.UsageForFunc) error {
		if len(args) < 3 {
			usageFor(usage.Init)()
		}

		outDir := args[2]
		fmt.Printf("Enter password to encrypt key with: ")
		pwd, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println()

		fmt.Printf("Confirm password: ")
		pwdConf, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println()

		if string(pwd) != string(pwdConf) {
			return fmt.Errorf("passwords don't match")
		}

		fmt.Printf("Enter azure storage account key: ")
		key, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		fmt.Println()

		err = init(fs, outDir, string(pwd), key)
		if err != nil {
			// Delete outDir if it was made
			os.RemoveAll(outDir)
		}
		return err
	}
}
