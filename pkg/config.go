package binny

import "net/http"

type Config struct {
	Load       string
	FileSystem http.FileSystem
	Build      string
	Image      string
	Args       []string
	Env        []string
	Workdir    string
}
