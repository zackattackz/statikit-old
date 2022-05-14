// Based on the bounded parallel Md5All pipeline,
// from https://web.archive.org/web/20220513193256/https://go.dev/blog/pipelines/bounded.go?m=text

package statikit

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Arguments to statikit.Render
type RendererArgs struct {
	InDir         string //Root input directory
	OutDir        string // Root output directory
	RendererCount uint   // # of renderer goroutines
	Data          any    // Data passed to template.Execute
}

// Combination of input/output paths
type inOutPath struct {
	in  string
	out string
}

func subtractPaths(parent, child string) string {
	parentList := strings.Split(parent, string(filepath.Separator))
	childList := strings.Split(child, string(filepath.Separator))

	return filepath.Join(childList[len(parentList):]...)
}

// Render the template at `p.in` to `p.out`, providing `data`
func render(p inOutPath, data any) error {
	fOut, err := os.Create(p.out)
	if err != nil {
		return err
	}
	defer fOut.Close()

	b, err := os.ReadFile(p.in)
	if err != nil {
		return err
	}

	t, err := template.New(p.out).Parse(string(b))
	if err != nil {
		return err
	}

	return t.Execute(fOut, data)
}

// renderer reads in/out paths from `paths` and sends result of rendering
// to `c` until either `paths` or `done` is closed.
func renderer(done <-chan struct{}, paths <-chan inOutPath, data any, c chan error) {
	for p := range paths {
		select {
		case c <- render(p, data):
		case <-done:
			return
		}
	}
}

// walkFiles starts a goroutine to walk the directory tree at root and send the
// in/out path of each "*.gohtml" file on `paths`.  It sends the result of the
// walk on the error channel.  If done is closed, walkFiles abandons its work.
// It copies directories and other regular files from in to out as it walks.
func walkFiles(done <-chan struct{}, data any, baseIn, baseOut string) (<-chan inOutPath, <-chan error) {
	paths := make(chan inOutPath)
	errc := make(chan error, 1)
	go func() {
		// Close the paths channel after Walk returns.
		defer close(paths)
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(baseIn, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			path = subtractPaths(baseIn, path)
			// Determine the full in/out paths for our file at `path`
			fullIn := filepath.Join(baseIn, path)
			fullOut := filepath.Join(baseOut, path)

			// If it's a directory, create that directory in baseOut/prefix
			if info.IsDir() {
				return os.Mkdir(fullOut, os.ModeDir)
			}

			// If not a directory or regular, error out
			if !info.Mode().IsRegular() {
				return fmt.Errorf("encountered irregular file: %s", fullIn)
			}

			// If the file is the config file, skip it
			matched, err := filepath.Match(ConfigFileName+".*", info.Name())
			if matched {
				return nil
			} else if err != nil {
				return err
			}

			// Otherwise, check if the file ends in ".gohtml"
			if filepath.Ext(fullIn) != ".gohtml" {

				// If it doesn't, copy file contents from `fullIn` to `fullOut`
				fIn, err := os.Open(fullIn)
				if err != nil {
					return err
				}
				defer fIn.Close()

				fOut, err := os.Create(fullOut)
				if err != nil {
					return err
				}
				defer fOut.Close()

				_, err = io.Copy(fOut, fIn)
				return err

			} else {
				// If it does, send in/out path to `paths`

				// Replace the ".gohtml" extension with ".html"
				fullOut = fullOut[:len(fullOut)-len(filepath.Ext(fullOut))] + ".html"

				// Send on paths or error out if `done` is closed
				select {
				case paths <- inOutPath{in: fullIn, out: fullOut}:
					return nil
				case <-done:
					return fmt.Errorf("walk canceled")
				}
			}
		})
	}()
	return paths, errc
}

// Orchestrates a pipeline that walks `a.InDir`,
// duplicating the all directories and files into `a.OutDir`.
// Except for any encountered "*.gohtml" files,
// which will be rendered as html.
func Render(a RendererArgs) error {
	if a.RendererCount < 1 {
		return fmt.Errorf("a.RendererCount must be >= 1")
	}

	// Render closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{})
	defer close(done)

	paths, errc := walkFiles(done, a.Data, a.InDir, a.OutDir)

	// Start a fixed number of goroutines to render files.
	c := make(chan error)
	var wg sync.WaitGroup
	wg.Add(int(a.RendererCount))
	for i := 0; i < int(a.RendererCount); i++ {
		go func() {
			renderer(done, paths, a.Data, c)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c)
	}()

	// Ensure all renderers don't error
	// If one does, return the error,
	// therfore closing done, which will stop other renderers
	for err := range c {
		if err != nil {
			return err
		}
	}

	// Return the result of the walk
	return <-errc
}
