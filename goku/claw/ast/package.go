package ast

type Link String

type Main String

type Unit String

type Module struct {
	Name String

	Links []Link
	Units []Unit

	Main *Main
}

type Set struct {
	Name Name
	Exp  Exp
}

type Package struct {
	Sets    []Set
	Modules []Module
}
