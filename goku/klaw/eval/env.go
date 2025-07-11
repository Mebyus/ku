package eval

type Env map[string]string

func (e Env) isTestExe() bool {
	return e["build.target.kind"] == "test"
}

func (e Env) isExe() bool {
	return e["build.target.kind"] == "exe"
}
