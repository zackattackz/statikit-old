package previewer

import (
	"fmt"
	"net/http"
)

type Previewer interface {
	Preview() error
}

type previewer struct {
	path string
	port string
}

func New(path, port string) Previewer {
	return &previewer{path: path, port: port}
}

func (t *previewer) Preview() error {
	http.Handle("/", http.FileServer(http.Dir(t.path)))
	fmt.Printf("Previewing contents of %s\n", t.path)
	fmt.Println("Listening on port " + t.port)
	return http.ListenAndServe(t.port, nil)
}
