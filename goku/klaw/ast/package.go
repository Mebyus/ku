package ast

type LinkEntry String

type MainEntry String

type UnitEntry String

type Module struct {
	Name String

	Links []LinkEntry
	Units []UnitEntry

	Main *MainEntry
}

type Set struct {
	Name Name
	Exp  Exp
}

type Package struct {
	Sets    []Set
	Modules []Module
}
