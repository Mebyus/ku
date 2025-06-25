package butler

import (
	"fmt"
	"io"
	"strings"
)

const debug = false

type Butler struct {
	// Short description of command.
	Short string

	// Usage hint.
	Usage string

	// Detailed command description.
	Full string

	// Command name. This is used for dispatching CLI arguments in command tree.
	Name string

	// List of subcommands.
	Subs []*Butler

	// This function will be used to handle remaining arguments if no suitable
	// subcommand is found. This field is optional, default help will be used if
	// left empty.
	Exec func(r *Butler, args []string) error

	// Template for parameters (flags) parsing.
	Params ParamBox

	// Contains command path elements. Each element in this slice is a subcommand name
	// that leads to this command. Always starts from root command. Last element
	// always contains name of this command.
	//
	// This field is calculated after init is called.
	elems []string

	// If true then non-nil unbound arguments will cause an error before
	// Exec is called. Ignored when there are subcommands present.
	DisallowUnboundArgs bool
}

// AddSub adds subcommand to the current command.
func (r *Butler) AddSub(sub *Butler) {
	r.Subs = append(r.Subs, sub)
}

// MustGetSub same as GetSub, but panics if no subcommand is found.
func (r *Butler) MustGetSub(name string) *Butler {
	sub := r.GetSub(name)
	if sub == nil {
		panic(fmt.Sprintf("command \"%s\" does not have \"%s\" subcommand", r.Name, name))
	}
	return sub
}

func (r *Butler) path() string {
	return strings.Join(r.elems, " ")
}

// GetSub returns a subcommand by its name. Returns nil if there is no such
// subcommand.
func (r *Butler) GetSub(name string) *Butler {
	if name == "" {
		panic("empty name")
	}

	for _, sub := range r.Subs {
		if sub.Name == name {
			return sub
		}
	}
	return nil
}

func (r *Butler) init(elems []string) {
	r.elems = elems
	if r.Name == "" {
		panic(fmt.Sprintf("empty subcommand name in \"%s\"", r.path()))
	}
	if r.Name == "help" {
		panic(fmt.Sprintf("subcommand name \"help\" (builtin) cannot be used in \"%s\"", r.path()))
	}
	if strings.HasPrefix(r.Name, "-") {
		panic(fmt.Sprintf("command name \"%s\" in \"%s\" starts with \"-\"", r.Name, r.path()))
	}

	r.elems = make([]string, len(elems), len(elems)+1)
	copy(r.elems, elems)
	r.elems = append(r.elems, r.Name)

	if len(r.Subs) == 0 {
		return
	}
	if len(r.Subs) == 1 {
		r.Subs[0].init(r.elems)
		return
	}

	set := make(map[string]struct{}, len(r.Subs))
	for _, sub := range r.Subs {
		name := sub.Name
		_, ok := set[name]
		if ok {
			panic(fmt.Sprintf("command \"%s\" contains more than one subcommand with name \"%s\"", r.path(), name))
		}
		set[name] = struct{}{}

		sub.init(r.elems)
	}
}

func (r *Butler) debugTreePrintPath(out io.Writer) {
	_, _ = fmt.Fprintf(out, "%v\n", r.elems)
	for _, sub := range r.Subs {
		sub.debugTreePrintPath(out)
	}
}

func Run(r *Butler, out io.Writer, args []string) error {
	r.init(nil)
	if debug {
		r.debugTreePrintPath(out)
	}
	return r.run(out, args)
}

func (r *Butler) run(out io.Writer, args []string) error {
	if len(r.Subs) == 0 {
		if r.Exec == nil {
			panic(fmt.Sprintf("command \"%s\" has no subcommands or executor", r.path()))
		}
		unbound, err := Parse(r.Params, args)
		if err != nil {
			return fmt.Errorf("parse command \"%s\" args: %v", r.path(), err)
		}
		if r.DisallowUnboundArgs && len(unbound) != 0 {
			return fmt.Errorf("command \"%s\" does not accept unbound arguments", r.path())
		}
		return r.Exec(r, unbound)
	}

	if len(args) == 0 {
		if r.Exec != nil {
			// User override for help
			return r.Exec(r, nil)
		}
		return r.displayHelp(out)
	}

	name := args[0]
	if name == "" {
		panic(fmt.Sprintf("empty name of \"%s\" subcommand invocation", r.path()))
	}
	if name == "help" {
		args = args[1:]
		if len(args) == 0 {
			return r.displayHelp(out)
		}
		subName := args[0]
		return r.subHelp(out, subName)
	}

	sub := r.GetSub(name)
	if sub == nil {
		return r.unknown(name)
	}
	return sub.run(out, args[1:])
}

func (r *Butler) subHelp(out io.Writer, name string) error {
	if name == "" {
		return r.displayHelp(out)
	}

	sub := r.GetSub(name)
	if sub == nil {
		return fmt.Errorf("command \"%s\" does not have \"%s\" help topic", r.path(), name)

	}
	return sub.displayHelp(out)
}

func (r *Butler) displayExecHelp(out io.Writer) error {
	_, err := io.WriteString(out, r.execHelp())
	return err
}

func (r *Butler) execHelp() string {
	var buf formatBuffer

	buf.puts(r.Short)
	buf.nl()
	buf.nl()

	buf.puts("Usage:")
	buf.nl()
	buf.nl()
	buf.indent()
	buf.puts(r.path())
	buf.puts(" ")
	buf.puts(r.Usage)
	buf.nl()
	buf.nl()

	if r.Full != "" {
		buf.puts(r.Full)
		buf.nl()
		buf.nl()
	}

	if r.Params == nil {
		return buf.String()
	}
	params := r.Params.Params()
	if len(params) == 0 {
		return buf.String()
	}

	buf.puts("Available options:")
	buf.nl()
	buf.nl()
	paramNamesColumnWidth := maxParamNameWidth(params) + 4
	for _, param := range params {
		buf.indent()
		buf.puts("--")
		buf.puts(param.Name)
		if param.Alias != "" {
			buf.puts(", -")
			buf.puts(param.Alias) // TODO: align text columns
		}
		buf.puts(strings.Repeat(" ", paramNamesColumnWidth-paramNameWidth(&param)))

		buf.puts("(")
		if param.Default != nil {
			buf.puts(fmt.Sprintf("=%v, ", param.Default))
		}
		buf.puts(param.Kind.String())
		buf.puts(")  ")
		buf.puts(param.Desc)
		buf.nl()
	}
	buf.nl()

	return buf.String()
}

func (r *Butler) displayHelp(out io.Writer) error {
	if r.Exec != nil {
		return r.displayExecHelp(out)
	}

	_, err := io.WriteString(out, r.help())
	return err
}

func (r *Butler) help() string {
	var buf formatBuffer

	buf.puts(r.Short)
	buf.nl()
	buf.nl()

	buf.puts("Usage:")
	buf.nl()
	buf.nl()
	buf.indent()
	buf.puts(r.path())
	buf.space()
	buf.puts(r.Usage)
	buf.nl()
	buf.nl()

	buf.puts("Available commands:")
	buf.nl()
	buf.nl()

	subNamesColumnWidth := r.maxSubNameWidth() + 4
	for _, s := range r.Subs {
		buf.indent()
		buf.puts(s.Name)
		buf.puts(strings.Repeat(" ", subNamesColumnWidth-len(s.Name)))
		buf.puts(s.Short)
		buf.nl()
	}

	buf.nl()
	buf.puts(fmt.Sprintf(`Use "%s help <command>" for more information about a specific command.`, r.path()))
	buf.nl()

	return buf.String()
}

func (r *Butler) maxSubNameWidth() int {
	w := 0
	for _, s := range r.Subs {
		if w < len(s.Name) {
			w = len(s.Name)
		}
	}
	return w
}

func maxParamNameWidth(params []Param) int {
	w := 0
	for _, s := range params {
		k := paramNameWidth(&s)
		if w < k {
			w = k
		}
	}
	return w
}

func paramNameWidth(p *Param) int {
	w := 2 + len(p.Name) // for example "--name"
	if p.Alias != "" {
		w += 4 // for example ", -n"
	}
	return w
}

func (r *Butler) unknown(name string) error {
	return fmt.Errorf("\"%s\" does not have \"%s\" subcommand", r.path(), name)
}
