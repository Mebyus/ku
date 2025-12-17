package stg

type Statement interface{}

type Block struct {
	Scope Scope

	Nodes []Statement
}
