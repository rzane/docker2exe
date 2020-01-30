package binny

type Config struct {
	Tarball string
	Image   string
	Args    []string
	Env     []string
	Workdir string
}
