package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/parser"
	"github.com/mebyus/ku/internal/ku/stg/hg"
	"github.com/mebyus/ku/internal/ku/sx"
)

type walker struct {
	queue queue

	// all units loaded during walk
	//
	// due to construction algorithm unit index in this slice
	// always equals its id
	units []*unit

	// list of unit ids with no imports in them
	roots []uid

	errors []*sx.Error

	// Base directory for standard library units lookup.
	std string

	// Base directory for local units lookup.
	loc string

	// maps unit path to its unique id (which equals to its index inside units list)
	m map[sx.Path]uid

	// set for checking duplicate imports inside single unit.
	iset map[sx.Path]struct{}

	pool *sx.Pool
}

// Walk discovers (recursively) units imported by program starting from given system or unit path.
// Returns list of all discovered unit paths.
//
// This function is meant to be used for testing and sandboxing.
func Walk(path string) []sx.Path {
	var p sx.Path
	if strings.HasPrefix(path, "./src/") {
		p = sx.MakePath(sx.Loc, filepath.Clean(strings.TrimPrefix(path, "./src/")))
	} else {
		p = sx.MakePath(sx.Loc, path)
	}

	var w walker
	w.init(witem{path: p})
	w.walk()

	for _, e := range w.errors {
		sx.FormatError(w.pool, os.Stderr, e)
	}

	if len(w.units) == 0 {
		return nil
	}

	list := make([]sx.Path, 0, len(w.units))
	for _, u := range w.units {
		list = append(list, u.path)
	}
	return list
}

// initialize walker with starting items
func (w *walker) init(items ...witem) {
	w.queue.init()
	w.pool = sx.New()
	w.m = make(map[sx.Path]uid)
	w.iset = make(map[sx.Path]struct{})

	for _, item := range items {
		w.queue.add(item)
	}
}

func (w *walker) walk() {
	for {
		var item witem
		if !w.queue.next(&item) {
			w.fill()
			return
		}

		var u unit
		ok := w.load(&item, &u)
		if !ok {
			// failed to load due to error, which
			// is already reported
			continue
		}
		id := uid(len(w.units))
		u.id = id
		w.units = append(w.units, &u)
		w.m[u.path] = id

		for _, s := range u.imports {
			w.queue.add(witem{
				path: s.path,
				pin:  s.pin,
			})
		}
	}
}

// fill each unit ancestors, descendants from imports,
// select root units
func (w *walker) fill() {
	const debug = false

	for i, u := range w.units {
		u.anc = w.getlist(u.imports)
		if len(u.anc) == 0 {
			w.roots = append(w.roots, uid(i))
		}

		for _, j := range u.anc {
			w.units[j].des = append(w.units[j].des, uid(i))
		}
	}

	if debug {
		for _, u := range w.units {
			fmt.Printf("[uwalk] %s | %d -> %v\n", u.path, u.id, u.anc)
		}
	}
}

func (w *walker) getlist(imports []isite) []uid {
	if len(imports) == 0 {
		return nil
	}

	list := make([]uid, 0, len(imports))
	for _, p := range imports {
		id, ok := w.m[p.path]
		if !ok {
			panic(fmt.Sprintf("unit \"%s\" not found in index map", p.path))
		}
		list = append(list, id)
	}
	slices.Sort(list)
	return list
}

// returns nil if there was an error loading unit files
// specified in walk item
func (w *walker) loadfiles(item *witem, u *unit) []*sx.Text {
	if item.file != "" {
		// TODO: better algorithm for generating unit name
		// from single file name
		name := "_file_" + strings.TrimSuffix(filepath.Base(item.file), ".ku")
		u.name = name

		file, err := w.pool.Load(item.file)
		if err != nil {
			w.report(&sx.Error{
				Pin:   item.pin,
				Short: fmt.Sprintf("load unit \"%s\" source file: %v", name, err),
			})
			return nil
		}
		return []*sx.Text{file}
	}

	dir := w.resolve(item.path)
	if dir == "" {
		return nil
	}

	u.path = item.path
	u.name = item.path.Name()
	u.dir = dir

	files, err := w.pool.LoadDir(&sx.LoadParams{
		Dir:              dir,
		IncludeTestFiles: item.includeTestFiles,
	})
	if err != nil {
		w.report(&sx.Error{
			Pin:   item.pin,
			Short: fmt.Sprintf("load unit \"%s\" source files: %v", item.path, err),
		})
		return nil
	}

	return files
}

// returns true if unit was loaded successfully (there may still be some errors)
func (w *walker) load(item *witem, u *unit) bool {
	files := w.loadfiles(item, u)
	if len(files) == 0 {
		return false
	}

	n := 0 // total number of imports in all texts
	var texts []*ast.Text
	for _, f := range files {
		t := parser.ParseText(f)
		texts = append(texts, t)
		n += len(t.Imports)
		u.errors = append(u.errors, t.Errors...)
	}
	u.texts = texts

	var imports []isite
	if n != 0 {
		clear(w.iset)
		imports = make([]isite, 0, n)
		for _, t := range texts {
			for _, p := range t.Imports {
				if p.Path == "" {
					// Imports with empty paths can come from
					// source bad syntax.
					//
					// It was already reported during parsing.
					// No need to produce error here.
					continue
				}

				if p.Path == item.path {
					w.report(&sx.Error{
						Pin:   p.ImpPin,
						Short: fmt.Sprintf("unit \"%s\" imports itself", item.path),
					})
					continue
				}

				_, ok := w.iset[p.Path]
				if ok {
					w.report(&sx.Error{
						Pin:   p.ImpPin,
						Short: fmt.Sprintf("unit \"%s\" imported multiple times", p.Path),
					})
					continue
				}
				w.iset[p.Path] = struct{}{}

				imports = append(imports, isite{
					path: p.Path,
					pin:  p.ImpPin,
				})
			}
		}
	}
	if len(imports) != 0 {
		u.imports = imports
	}

	return true
}

// translate unit path into system path of directory where
// source code for this unit is stored
func (w *walker) resolve(path sx.Path) string {
	o, s := path.Import()
	if s == "" {
		// should be already reported during parsing
		return ""
	}

	switch o {
	case 0:
		// should be already reported during parsing
	case sx.Std:
		return filepath.Join(w.std, s)
	case sx.Pkg:
		panic("stub")
	case sx.Loc:
		// TODO: use absolute path for src?
		return filepath.Join(w.loc, s)
	default:
		// should be already reported during parsing
	}
	return ""
}

func (w *walker) report(e *sx.Error) {
	w.errors = append(w.errors, e)
}

// import site represents a place in source code which imports a unit
type isite struct {
	// unit being imported
	path sx.Path

	// where it was imported
	pin sx.Pin
}

// unique (within a program) unit id
type uid = hg.UnitID

// intermediate object for holding various data related to loaded unit
type unit struct {
	// list of unit paths imported by this unit
	imports []isite

	// list of unit ancestors
	//
	// in other words this list contains unit ids imported
	// by this unit
	//
	// empty for root units
	anc []uid

	// list of unit descendants
	//
	// in other words this list contains unit ids which
	// import this unit
	des []uid

	texts []*ast.Text

	errors []*ast.Error

	// identifies this unit within a program
	path sx.Path

	// short name, which later will be used as alias
	//
	// usually we fill it based on directory name
	name string

	// system path to directory with unit files
	dir string

	// contains unit index in the list of all loaded units
	// maintained by walker
	id uid
}

// walk item carries data needed to load a new unit
type witem struct {
	// unit path we want to load
	path sx.Path

	// system path to a single source file
	//
	// if not empty use this instead of unit path to load
	// unit from a single file instead of whole directory
	//
	// TODO: maybe we should just merge unit path and raw
	// file path into one entity with distinct origins
	file string

	// Place where unit with this path is imported.
	//
	// There may be more than one import place for a specific unit inside the
	// whole program. This field tracks the first one we encounter during unit
	// walk phase.
	pin sx.Pin

	// If true, then test files should be loaded alongside unit source files.
	includeTestFiles bool
}

// queue keeps track which units were already visited and
// which unit should be visited next during unit discovery phase.
type queue struct {
	// List of paths waiting in queue.
	backlog []witem

	// Set which contains paths of already visited units.
	visited map[sx.Path]struct{}
}

func (q *queue) init() {
	q.visited = make(map[sx.Path]struct{})
}

// Add tries to add item to backlog. If a given path was
// already visited then this call will be no-op.
func (q *queue) add(item witem) {
	switch {
	case item.file != "":
		// continue execution
	case item.path != "":
		_, ok := q.visited[item.path]
		if ok {
			// given path is already known to walker
			return
		}
		q.visited[item.path] = struct{}{}
	default:
		panic("empty walk item")
	}

	q.backlog = append(q.backlog, item)
}

// Next get next item from backlog and write it to a given pointer.
// Returns true if operation was successfull and false when there are no items left.
func (q *queue) next(item *witem) bool {
	if len(q.backlog) == 0 {
		return false
	}

	last := len(q.backlog) - 1
	*item = q.backlog[last]

	// shrink slice, but keep its capacity
	q.backlog = q.backlog[:last]
	return true
}
