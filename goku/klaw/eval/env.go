package eval

type Env struct {
	Exe     bool
	TestExe bool

	m map[string]string
}

func NewEnv() *Env {
	return &Env{m: make(map[string]string)}
}
