package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tidwall/secret"
	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/config"
	"github.com/zackattackz/azure_static_site_kit/pkg/statikit/publisher"
	"golang.org/x/term"
)

func publish(path string) error {
	aes, err := os.ReadFile(filepath.Join(path, config.KeyFileName))
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

	a := publisher.Args{Path: path, Key: string(key)}
	return publisher.Publish(a)
}
