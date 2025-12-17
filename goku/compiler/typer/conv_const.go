package typer

import (
	"fmt"

	"github.com/mebyus/ku/goku/compiler/diag"
	"github.com/mebyus/ku/goku/compiler/enums/tpk"
	"github.com/mebyus/ku/goku/compiler/typer/stg"
)

func (t *Typer) convConstSymbol(s *stg.Symbol) diag.Error {
	const debug = false

	c := t.box.Const(s.Aux)
	name := c.Name.Str

	exp, err := t.unit.Scope.EvalConstExp(c.Exp)
	if err != nil {
		return err
	}

	var typ *stg.Type
	if c.Type == nil {
		typ = exp.Type()
	} else {
		var err diag.Error
		typ, err = t.ctx.Types.Lookup(&t.unit.Scope, c.Type)
		if err != nil {
			return err
		}
		err = checkConstValueType(typ, exp)
		if err != nil {
			return err
		}
		// TODO: convert exp to the same type as declared
	}

	s.Def = stg.StaticValue{Exp: exp}

	if debug {
		fmt.Printf("const %s: %s = %s\n", name, typ, exp)
	}

	return nil
}

func checkConstValueType(want *stg.Type, exp stg.Exp) diag.Error {
	t := exp.Type()

	switch want.Kind {
	case tpk.Custom:
		c := want.Def.(*stg.Custom)
		if c.Type.Kind == tpk.Integer && t.Kind == tpk.Integer && t.Size == 0 {
			return nil
		}
	default:
		panic(fmt.Sprintf("unexpected \"%s\" (=%d) type (%T)", want.Kind, want.Kind, want.Def))
	}

	return &diag.SimpleMessageError{
		Pin:  exp.Span().Pin,
		Text: fmt.Sprintf("declared constant type %s and definition value type %s are incompatible", want, exp.Type()),
	}
}
