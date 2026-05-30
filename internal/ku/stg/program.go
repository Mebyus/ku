package stg

import (
	"github.com/mebyus/ku/internal/ku/enums/scok"
	"github.com/mebyus/ku/internal/ku/sx"
)

type Program struct {
	// Common

	Units []*Unit
}

type Unit struct {
	Scope Scope

	Funs []*Symbol

	Errors []*Error

	Name string
}

func (u *Unit) init(global *Scope) {
	u.Scope.Init(scok.Unit, global)
}

type Error struct {
	Short string
	Pin   sx.Pin
}
