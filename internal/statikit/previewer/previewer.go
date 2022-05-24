package previewer

import (
	"fmt"
	"net/http"
)

type Args struct {
	Path string
	Port string
}

type PreviewFunc func(Args) error

func Preview(a Args) error {
	http.Handle("/", http.FileServer(http.Dir(a.Path)))
	fmt.Printf("Previewing contents of %s\n", a.Path)
	fmt.Println("Listening on port " + a.Port)
	return http.ListenAndServe(a.Port, nil)
}
