package previewer

import (
	"fmt"
	"net/http"

	"github.com/spf13/afero"
)

func Preview(fs afero.Fs, path string, port string) error {
	http.Handle("/", http.FileServer(http.Dir(path)))
	fmt.Printf("Previewing contents of %s\n", path)
	fmt.Println("Listening on port " + port)
	return http.ListenAndServe(port, nil)
}
