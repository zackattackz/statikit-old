package previewer

import (
	"fmt"
	"net/http"
)

const port = ":8080"

func Preview(path string) error {
	http.Handle("/", http.FileServer(http.Dir(path)))
	fmt.Printf("Previewing contents of %s\n", path)
	fmt.Println("Listening on port " + port)
	return http.ListenAndServe(port, nil)
}
