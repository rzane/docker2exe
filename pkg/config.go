package binny

import "net/http"

type Config struct {
	Build   string
	Image   string
	Args    []string
	Env     []string
	Workdir string
	Load    bool
	Open    func() (http.File, error)
}
