package main

import (
	"fmt"
	"os"

	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
	"golang.org/x/term"
)

type initArgs struct {
	path string
}

func initialize(a initArgs) error {

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

	err = initializer.Create(a.path, string(pwd), key)
	if err != nil {
		// Delete outDir if it was made
		os.RemoveAll(a.path)
	}
	return err
}
