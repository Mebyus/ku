package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mebyus/ku/internal/ku/ast"
	"github.com/mebyus/ku/internal/ku/parser"
	"github.com/mebyus/ku/internal/ku/sx"
)

type walker struct {
	queue queue

	// all units loaded during walk
	units []*unit

	errors []*sx.Error

	// Base directory for standard library units lookup.
	std string

	// Base directory for local units lookup.
	loc string

	// maps unit path to its unique id (which equals to its index inside units list)
	m map[sx.Path]uid

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

	for _, item := range items {
		w.queue.add(item)
	}
}

func (w *walker) walk() {
	for {
		var item witem
		if !w.queue.next(&item) {
			return
		}

		var u unit
		ok := w.load(item, &u)
		if !ok {
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

// returns true if unit was loaded successfully (there may still be some errors)
func (w *walker) load(item witem, u *unit) bool {
	dir := w.resolve(item.path)
	if dir == "" {
		return false
	}

	u.path = item.path
	u.name = item.path.Name()
	u.dir = dir

	files, loadErr := w.pool.LoadDir(&sx.LoadParams{
		Dir:              dir,
		IncludeTestFiles: item.includeTestFiles,
	})
	if loadErr != nil {
		w.addError(&sx.Error{
			Pin:   item.pin,
			Short: fmt.Sprintf("load unit \"%s\" source files: %v", item.path, loadErr),
		})
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
		for _, t := range texts {
			for _, m := range t.Imports {
				imports = append(imports, isite{
					path: m.Path,
					pin:  m.ImpPin,
				})
			}
		}
	}
	u.imports = imports

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

func (w *walker) addError(e *sx.Error) {
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
type uid uint32

// intermediate object for holding various data related to loaded unit
type unit struct {
	// list of unit paths imported by this unit
	imports []isite

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
	_, ok := q.visited[item.path]
	if ok {
		// given path is already known to walker
		return
	}

	q.visited[item.path] = struct{}{}
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
