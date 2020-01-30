package binny

type Config struct {
	Load    string
	Build   string
	Image   string
	Args    []string
	Env     []string
	Workdir string
}
