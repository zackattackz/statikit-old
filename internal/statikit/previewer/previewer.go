package previewer

import (
	"fmt"
	"net/http"
)

type Interface interface {
	Preview() error
}

type T struct {
	path string
	port string
}

func New(path, port string) T {
	return T{path: path, port: port}
}

func (t *T) Preview() error {
	http.Handle("/", http.FileServer(http.Dir(t.path)))
	fmt.Printf("Previewing contents of %s\n", t.path)
	fmt.Println("Listening on port " + t.port)
	return http.ListenAndServe(t.port, nil)
}
