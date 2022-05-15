package previewer

import (
	"fmt"
	"net/http"
	"os"
)

const port = ":8080"

func Preview(path string) {
	http.Handle("/", http.FileServer(http.Dir(path)))
	fmt.Printf("Previewing contents of %s\n", path)
	fmt.Println("Listening on port " + port)
	fmt.Fprintln(os.Stderr, http.ListenAndServe(port, nil))
}
