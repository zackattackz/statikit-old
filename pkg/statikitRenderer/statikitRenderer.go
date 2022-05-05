package statikitRenderer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Args struct {
	InDir  string
	OutDir string
}

func render(fIn, fOut *os.File) error {
	fmt.Println(fIn.Name())
	return nil
}

func initHandleEntry(baseOutPath string) fs.WalkDirFunc {
	return func(inPath string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		outPath := filepath.Join(baseOutPath, inPath)

		// Handle the current entry
		// If it's a directory, create that directory in outPath
		if entry.IsDir() {
			return os.Mkdir(outPath, os.ModeDir)
		}
		// Otherwise it is a file,
		// If the input file is a .gohtml file, render it
		if filepath.Ext(inPath) == ".gohtml" {
			fIn, err := os.Open(inPath)
			if err != nil {
				return err
			}
			fOut, err := os.Create(outPath)
			if err != nil {
				return err
			}
			return render(fIn, fOut)
		} else {
			// Otherwise, hard link the file from inPath to outPath
			return os.Link(inPath, outPath)
		}
	}
}

func Render(a Args) error {
	if err := os.RemoveAll(a.OutDir); err != nil {
		return err
	}

	handleEntry := initHandleEntry(a.OutDir)
	return filepath.WalkDir(a.InDir, handleEntry)

}
