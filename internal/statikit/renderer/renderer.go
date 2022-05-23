// Based on the bounded parallel Md5All pipeline,
// from https://web.archive.org/web/20220513193256/https://go.dev/blog/pipelines/bounded.go?m=text

package renderer

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/zackattackz/azure_static_site_kit/internal/statikit/initializer"
	"github.com/zackattackz/azure_static_site_kit/internal/statikit/schema"
	sp "github.com/zackattackz/azure_static_site_kit/pkg/subtractPaths"
)

type Renderer interface {
	Render() error
}

// Arguments to statikit.Render
type Args struct {
	InDir         string     //Root input directory
	OutDir        string     // Root output directory
	RendererCount uint       // # of renderer goroutines
	SchemaMap     schema.Map // Scehma passed to template.Execute
	Ignore        []string   // Filepath globs to ignore when walking
}

// Combination of input/output paths
type inOutPath struct {
	in  string
	out string
}

// Render the template at `p.in` to `p.out`, providing `data`
func render(p inOutPath, dataMap schema.Map, baseIn string) error {
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

	path := sp.SubtractPaths(baseIn, p.in)
	pathWithoutExt := strings.TrimSuffix(path, filepath.Ext(path))
	d, ok := dataMap[pathWithoutExt]
	if !ok {
		d = dataMap[initializer.DefaultDataName]
	}

	return t.Execute(fOut, d)
}

// renderWorker reads in/out paths from `paths` and sends result of rendering
// to `c` until either `paths` or `done` is closed.
func renderWorker(done <-chan struct{}, paths <-chan inOutPath, dataMap schema.Map, baseIn string, c chan error) {
	for p := range paths {
		select {
		case c <- render(p, dataMap, baseIn):
		case <-done:
			return
		}
	}
}

// walkFiles starts a goroutine to walk the directory tree at root and send the
// in/out path of each "*.gohtml" file on `paths`.  It sends the result of the
// walk on the error channel.  If done is closed, walkFiles abandons its work.
// It copies directories and other regular files from in to out as it walks.
func walkFiles(done <-chan struct{}, baseIn, baseOut string, ignore []string) (<-chan inOutPath, <-chan error) {
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

			path = sp.SubtractPaths(baseIn, path)

			// If path is in ignore list, skip it
			for _, ignorePath := range ignore {
				if match, _ := filepath.Match(ignorePath, path); match {
					if info.IsDir() {
						return fs.SkipDir
					} else {
						return nil
					}
				}
			}

			// Determine the full in/out paths for our file at `path`
			fullIn := filepath.Join(baseIn, path)
			fullOut := filepath.Join(baseOut, path)

			// If it's a directory, create that directory in baseOut/prefix
			if info.IsDir() {
				return os.Mkdir(fullOut, 0755)
			}

			// If not a directory or regular, error out
			if !info.Mode().IsRegular() {
				return fmt.Errorf("encountered irregular file: %s", fullIn)
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
				// If it does end in ".gohtml",
				// send in/out path to `paths`

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

type renderer struct {
	Args
}

func New(a Args) Renderer {
	return &renderer{a}
}

// Orchestrates a pipeline that walks `a.InDir`,
// duplicating the all directories and files into `a.OutDir`.
// Except for any encountered "*.gohtml" files,
// which will be rendered as html.
func (t *renderer) Render() error {
	if t.RendererCount < 1 {
		return fmt.Errorf("a.RendererCount must be >= 1")
	}

	// Render closes the done channel when it returns; it may do so before
	// receiving all the values from c and errc.
	done := make(chan struct{})
	defer close(done)

	// Start the file walking goroutine
	paths, errc := walkFiles(done, t.InDir, t.OutDir, t.Ignore)

	// Start a fixed number of goroutines to render files.
	c := make(chan error)
	var wg sync.WaitGroup
	wg.Add(int(t.RendererCount))
	for i := 0; i < int(t.RendererCount); i++ {
		go func() {
			renderWorker(done, paths, t.SchemaMap, t.InDir, c)
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
