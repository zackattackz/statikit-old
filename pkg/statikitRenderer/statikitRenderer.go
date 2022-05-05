package statikitRenderer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Args struct {
	InDir  string
	OutDir string
}

func render(c chan<- []error, f *os.File, fPath string) {
	fmt.Println(fPath)
	c <- nil
}

func recurseRender(c chan<- []error, in []os.DirEntry, parentInPath string, parentOutPath string) {
	var errs []error
	var dirsToSearch []os.DirEntry
	renderDone := make(chan []error)

	childCount := 0

	for _, e := range in {
		if e.IsDir() {
			dirsToSearch = append(dirsToSearch, e)
		} else {
			// Render e
			fInPath := parentInPath + filepath.FromSlash("/") + e.Name()
			fmt.Printf("Rendering file: %v\n", e.Name())
			if f, err := os.Open(fInPath); err != nil {
				errs = append(errs, err)
			} else {
				fOutPath := parentOutPath + filepath.FromSlash("/") + e.Name()
				go render(renderDone, f, fOutPath)
				childCount += 1
			}
		}
	}

	for _, d := range dirsToSearch {
		dInPath := parentInPath + filepath.FromSlash("/") + d.Name()
		if e, err := os.ReadDir(dInPath); err != nil {
			errs = append(errs, err)
		} else {
			dOutPath := parentOutPath + filepath.FromSlash("/") + d.Name()
			if err := os.Mkdir(dOutPath, os.ModeDir); err != nil {
				errs = append(errs, err)
			} else {
				go recurseRender(renderDone, e, dInPath, dOutPath)
				childCount += 1
			}
		}
	}

	for i := 0; i < childCount; i += 1 {
		childErrs := <-renderDone
		errs = append(errs, childErrs...)
	}

	c <- errs
}

func Render(a Args) []error {
	var errs []error

	in, err := os.ReadDir(a.InDir)
	if err != nil {
		errs = append(errs, err)
		return errs
	}

	if err := os.RemoveAll(a.OutDir); err != nil {
		errs = append(errs, err)
		return errs
	}

	if err := os.Mkdir(a.OutDir, os.ModeDir); err != nil {
		if !errors.Is(err, os.ErrExist) {
			errs = append(errs, err)
			return errs
		}
	}

	renderDone := make(chan []error)

	go recurseRender(renderDone, in, a.InDir, a.OutDir)

	errs = append(errs, <-renderDone...)

	return errs
}
