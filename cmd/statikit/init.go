package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/secret"
	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/config"
	"golang.org/x/term"
)

func initialize(path string) error {

	const keyFileName = "key.aes256"

	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("%s already exists", path)
	}

	if err = os.Mkdir(path, 0755); err != nil {
		return err
	}

	configFile, err := os.Create(filepath.Join(path, config.ConfigFileName+".toml"))
	if err != nil {
		return err
	}
	defer configFile.Close()
	configFile.WriteString("[Data]")

	keyFile, err := os.Create(filepath.Join(path, keyFileName))
	if err != nil {
		return err
	}
	defer keyFile.Close()

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

	aes, err := secret.Encrypt(string(pwd), key)
	if err != nil {
		return err
	}

	keyFile.Write(aes)

	return nil

}
