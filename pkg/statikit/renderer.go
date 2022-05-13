// Based on the bounded parallel Md5All pipeline,
// from https://web.archive.org/web/20220513193256/https://go.dev/blog/pipelines/bounded.go?m=text

package statikit

import (
	"errors"
	"html/template"
	"os"
	"path/filepath"
	"sync"
)

type RendererArgs struct {
	InDir         string
	OutDir        string
	RendererCount uint
	Data          any
}

type rendererInput struct {
	fOut *os.File
	t    *template.Template
	data any
}

// renderer reads from inputs and sends result of rendering
// to c until either inputs or done is closed.
func renderer(done <-chan struct{}, inputs <-chan rendererInput, c chan error) {
	for input := range inputs {
		err := input.t.Execute(input.fOut, input.data)
		input.fOut.Close()
		select {
		case c <- err:
		case <-done:
			return
		}
	}
}

// walkFiles starts a goroutine to walk the directory tree at root and send the
// path of each regular file on the string channel.  It sends the result of the
// walk on the error channel.  If done is closed, walkFiles abandons its work.
func walkFiles(done <-chan struct{}, data any, baseIn, baseOut string) (<-chan rendererInput, <-chan error) {
	toRender := make(chan rendererInput)
	errc := make(chan error, 1)
	go func() {
		// Close the paths channel after Walk returns.
		defer close(toRender)
		// No select needed for this send, since errc is buffered.
		errc <- filepath.Walk(baseIn, func(path string, info os.FileInfo, err error) error {

			if err != nil {
				return err
			}

			fullOut := filepath.Join(baseOut, path)
			fullIn := filepath.Join(baseIn, path)

			// If it's a directory, create that directory in baseOut/prefix
			if info.IsDir() {
				return os.Mkdir(fullOut, os.ModeDir)
			}

			if filepath.Ext(path) != ".gohtml" {
				// If regular file, hard link the file from fullIn to fullOut
				return os.Link(fullIn, fullOut)
			} else {
				// If .gohtml, create a template from file at `path` and send to `toRender`

				// Open `path`
				fIn, err := os.Open(path)
				if err != nil {
					return err
				}

				// Open fullOut, except replace the ".gohtml" extension with ".html"
				fullOutNoExt := fullOut[:len(fullOut)-len(filepath.Ext(fullOut))]
				fOut, err := os.Create(fullOutNoExt + ".html")
				if err != nil {
					return err
				}

				// Now read fIn into a byte buffer `b`
				fInStat, err := fIn.Stat()
				if err != nil {
					fOut.Close()
					return err
				}
				s := fInStat.Size()
				b := make([]byte, s)
				for s > 0 {
					n, err := fIn.Read(b)
					if err != nil {
						fOut.Close()
						fIn.Close()
						return err
					}
					s -= int64(n)
				}

				fIn.Close()

				// Create a template by parsing `b`
				t, err := template.New(fInStat.Name()).Parse(string(b))
				if err != nil {
					return err
				}

				// Send `fOut`, `t`, and `data` on the `toRender` channel
				select {
				case toRender <- rendererInput{fOut: fOut, t: t, data: data}:
				case <-done:
					// (return early if done)
					fOut.Close()
					return errors.New("walk canceled")
				}

				return nil
			}
		})
	}()
	return toRender, errc
}

// Orchestrates a pipeline that walks a.InDir,
// duplicating the all directories and files into a.OutDir.
// Except for any encountered "*.gohtml" files,
// which will be rendered as html.
func Render(a RendererArgs) error {
	if a.RendererCount < 1 {
		return errors.New("a.RendererCount must be >= 1")
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
			renderer(done, paths, c)
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
