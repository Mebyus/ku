package butler

import (
	"fmt"
	"strings"
)

type ParamKind uint32

const (
	empty ParamKind = iota

	Boolean
	String
	Integer
	List
	CommaList
)

var kindText = [...]string{
	empty: "<nil>",

	Boolean:   "boolean",
	String:    "string",
	Integer:   "integer",
	List:      "list",
	CommaList: "list.comma",
}

func (k ParamKind) String() string {
	return kindText[k]
}

func (k ParamKind) Valid() error {
	switch k {
	case empty:
		return fmt.Errorf("unspecified kind")
	case Boolean, String, Integer, List, CommaList:
		return nil
	default:
		return fmt.Errorf("unexpected kind (=%d)", k)
	}
}

type Param struct {
	Name string

	// Optional short param name.
	//
	// Must be either empty (means no alias) or contain one latin letter.
	Alias string

	Desc string

	// Default value for this param.
	//
	// If param is required then this field should be nil.
	Default any

	// stored bound param values
	val any

	Kind ParamKind
}

func (p *Param) Bind(v string) error {
	switch p.Kind {
	case empty:
		panic("unspecified kind")
	case Boolean:
		switch v {
		case "true":
			p.val = true
		case "false":
			p.val = false
		default:
			return fmt.Errorf("bad boolean value \"%s\"", v)
		}
	case String:
		p.val = v
	default:
		panic(fmt.Sprintf("unxpected kind (=%d)", p.Kind))
	}
	return nil
}

type ParamBox interface {
	Apply(*Param) error
	Params() []Param
}

func (r *parser) index() {
	params := paramsFromBox(r.box)
	if len(params) == 0 {
		return
	}

	m := make(map[string]*Param, len(params))
	r.m = m
	for i := range len(params) {
		p := params[i]
		name := strings.TrimSpace(p.Name)
		if name == "" {
			panic("empty param name")
		}
		if name != p.Name {
			panic(fmt.Sprintf("param name \"%s\" contains whitespace", p.Name))
		}
		if strings.HasPrefix(name, "-") {
			panic(fmt.Sprintf("param name \"%s\" starts with \"-\"", p.Name))
		}
		err := p.Kind.Valid()
		if err != nil {
			panic(fmt.Sprintf("param \"%s\": %s", name, err))
		}
		_, ok := m[name]
		if ok {
			panic(fmt.Sprintf("param \"%s\" is not unique", name))
		}
		m[name] = &p

		alias := strings.TrimSpace(p.Alias)
		if alias != p.Alias {
			panic("param alias contains whitespace")
		}
		if alias == name {
			panic(fmt.Sprintf("param \"%s\" has identical name and alias", name))
		}
		if alias == "" {
			// means that param has no alias
			continue
		}
		if len(alias) != 1 {
			panic(fmt.Sprintf("param \"%s\" alias \"%s\" is longer than one latin letter", name, alias))
		}

		_, ok = m[alias]
		if ok {
			panic(fmt.Sprintf("param \"%s\" alias \"%s\" is not unique", name, alias))
		}
		m[alias] = &p
	}
}

func paramsFromBox(box ParamBox) []Param {
	if box == nil {
		return nil
	}
	return box.Params()
}

// Parse parses supplied arguments according to ParamBox container.
//
// If parsing was successful returns unbound args and error otherwise.
// Param values which were found in arguments are applied to ParamBox container.
func Parse(box ParamBox, args []string) ([]string, error) {
	r := parser{
		args: args,
		box:  box,
	}
	r.index()

	unbound, err := r.bind()
	if err != nil {
		return nil, err
	}

	err = r.apply()
	if err != nil {
		return nil, err
	}

	return unbound, nil
}

type parser struct {
	args []string

	box ParamBox

	m map[ /* param name */ string]*Param

	// index of next unprocessed arg from slice
	i int
}

// Returns next element from stored list of arguments and an ok flag.
// Flag equals false if no more args are available (and thus iterator is exhausted).
// Advances iterator to next element if ok flag equals true.
func (r *parser) next() (string, bool) {
	next, ok := r.peek()
	if !ok {
		return "", false
	}

	r.advance()
	return next, true
}

// Same as next(), but does not advance arguments iterator.
func (r *parser) peek() (string, bool) {
	if r.i >= len(r.args) {
		return "", false
	}

	return r.args[r.i], true
}

func (r *parser) advance() {
	r.i += 1
}

// tail returns remaining (unprocessed) args.
func (r *parser) tail() []string {
	return r.args[r.i:]
}

func (r *parser) apply() error {
	for _, p := range r.m {
		if p.Default == nil && p.val == nil {
			return fmt.Errorf("param \"%s\" is required", p.Name)
		}
		err := r.box.Apply(p)
		if err != nil {
			return fmt.Errorf("apply param \"%s\" value %v: %w", p.Name, p.val, err)
		}
	}
	return nil
}

func (r *parser) bind() ([]string, error) {
	if len(r.args) == 0 {
		return nil, nil
	}
	if len(r.m) == 0 {
		return r.args, nil
	}

	for {
		arg, ok := r.peek()
		if !ok {
			return nil, nil
		}
		if arg == "" {
			panic("empty arg")
		}

		if arg == "--" {
			// explicit end of param arguments
			r.advance()
			return r.tail(), nil
		}

		suffix, ok := ParseParamPrefix(arg)
		if !ok {
			// end parsing because we encountered first unbound argument
			return r.tail(), nil
		}

		r.advance()

		err := r.bindSuffix(suffix)
		if err != nil {
			return nil, err
		}
	}
}

func (r *parser) bindSuffix(suffix string) error {
	j := strings.Index(suffix, "=")
	if j < 0 {
		return r.bindNextValueByParamName(suffix)
	}

	name := suffix[:j]
	value := suffix[j+1:]
	return r.bindParamAndValueFromSuffix(name, value)
}

func bindParamAndValue(p *Param, v string) error {
	err := p.Bind(v)
	if err != nil {
		return fmt.Errorf("bind param \"%s\" value \"%s\": %w", p.Name, v, err)
	}
	return nil
}

func (r *parser) bindNextValueByParamName(name string) error {
	p, ok := r.m[name]
	if !ok {
		return fmt.Errorf("unknown param \"%s\"", name)
	}
	if p.Kind == Boolean {
		p.val = true
		return nil
	}

	v, ok := r.next()
	if !ok {
		return fmt.Errorf("no value provided for param \"%s\"", name)
	}
	if v == "" {
		panic("empty arg")
	}

	return bindParamAndValue(p, v)
}

func (r *parser) bindParamAndValueFromSuffix(name string, v string) error {
	if name == "" || v == "" {
		return fmt.Errorf("invalid arg syntax \"%s=%s\"", name, v)
	}

	p, ok := r.m[name]
	if !ok {
		return fmt.Errorf("unknown param \"%s\"", name)
	}
	return bindParamAndValue(p, v)
}

// ParseParamPrefix gets param name (and possibly "=value" part) from a given cli argument.
// Returns (suffix, true) if the argument could be a param name.
// Otherwise returns ("", false).
func ParseParamPrefix(arg string) (string, bool) {
	i := 0
	for arg[i] == '-' {
		i += 1
	}

	if i == 1 || i == 2 {
		return arg[i:], true
	}
	return "", false
}

func (p *Param) Bool() bool {
	if p.val == nil {
		if p.Default == nil {
			panic(fmt.Sprintf("param \"%s\" default value cannot be nil at this point", p.Name))
		}
		return p.Default.(bool)
	}
	return p.val.(bool)
}

func (p *Param) List() []string {
	if p.val == nil {
		if p.Default == nil {
			panic(fmt.Sprintf("param \"%s\" default value cannot be nil at this point", p.Name))
		}
		return p.Default.([]string)
	}
	return p.val.([]string)
}

func (p *Param) Str() string {
	if p.val == nil {
		if p.Default == nil {
			panic(fmt.Sprintf("param \"%s\" default value cannot be nil at this point", p.Name))
		}
		return p.Default.(string)
	}
	return p.val.(string)
}
