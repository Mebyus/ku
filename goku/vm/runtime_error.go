package vm

type RuntimeError struct {
}

func (r *RuntimeError) Error() string {
	return ""
}
