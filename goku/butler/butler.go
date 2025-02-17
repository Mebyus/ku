package butler

import (
	"fmt"
	"io"
	"os"
	"strings"
)

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
	Lackeys []*Butler

	// This function will be used to handle remaining arguments if no suitable
	// subcommand is found. This field is optional, default help will be used if
	// left empty.
	Exec func(r *Butler, args []string) error

	// Template for parameters (flags) parsing.
	Params ParamBox

	// concatenated list of command names that leads to this command
	prefix string
}

func Run(r *Butler, out io.Writer, args []string) error {
	r.init()
	return r.run(out, args)
}

func (r *Butler) init() {

}

func (r *Butler) run(out io.Writer, args []string) error {
	if len(r.Lackeys) == 0 {
		if r.Exec == nil {
			panic(fmt.Sprintf("command \"%s\" has no subcommands or executor", r.prefix))
		}
		unbound, err := Parse(r.Params, args)
		if err != nil {
			return err
		}
		return r.Exec(r, unbound)
	}

	if len(args) == 0 && r.Exec == nil {
		return r.help(out)
	}

	name := args[0]
	if name == "help" {
		args = args[1:]
		if len(args) == 0 {
			return r.help(out)
		}
		subName := args[0]
		return r.subHelp(out, subName)
	}
	for _, s := range r.Lackeys {
		if name == s.Name {
			return s.run(out, args[1:])
		}
	}
	return r.unknown(name)
}

func (r *Butler) subHelp(out io.Writer, name string) error {
	if name == "" {
		return r.help(out)
	}

	for _, s := range r.Lackeys {
		if name == s.Name {
			return s.help(out)
		}
	}
	return fmt.Errorf("unknown help topic \"%s\"", name)
}

func (r *Butler) execHelp(out io.Writer) error {
	r.write(r.Short)
	r.nl()
	r.nl()

	r.write("Usage:")
	r.nl()
	r.nl()
	r.indent()
	r.write(r.Usage)
	r.nl()
	r.nl()

	if r.Full != "" {
		r.write(r.Full)
		r.nl()
		r.nl()
	}

	if r.Params == nil {
		return nil
	}

	r.write("Available options:")
	r.nl()
	r.nl()
	params := r.Params.Params()
	for _, param := range params {
		r.indent()
		r.write("--")
		r.write(param.Name)
		r.write(" (")
		r.write(param.Kind.String())
		r.write(")")
		r.write("    ")
		r.write(param.Desc)
		r.nl()
	}
	r.nl()
	return nil
}

func (r *Butler) help(out io.Writer) error {
	if r.Exec != nil {
		r.execHelp(out)
		return nil
	}

	r.write(r.Short)
	r.nl()
	r.nl()

	r.write("Usage:")
	r.nl()
	r.nl()
	r.indent()
	r.write(r.Usage)
	r.nl()
	r.nl()

	r.write("Available commands:")
	r.nl()
	r.nl()

	subNamesColumnWidth := r.maxSubNameWidth() + 4
	for _, s := range r.Lackeys {
		r.indent()
		r.write(s.Name)
		r.write(strings.Repeat(" ", subNamesColumnWidth-len(s.Name)))
		r.write(s.Short)
		r.nl()
	}

	r.nl()
	r.write(fmt.Sprintf(`Use "%s help <command>" for more information about a specific command.`, r.Name))
	r.nl()
	return nil
}

func (r *Butler) maxSubNameWidth() int {
	m := 0
	for _, s := range r.Lackeys {
		if m < len(s.Name) {
			m = len(s.Name)
		}
	}
	return m
}

func (r *Butler) nl() {
	r.write("\n")
}

func (r *Butler) indent() {
	r.write("\t")
}

func (r *Butler) write(s string) {
	io.WriteString(os.Stderr, s)
}

func (r *Butler) unknown(name string) error {
	return fmt.Errorf("%s does not recognize '%s' command", r.Name, name)
}
