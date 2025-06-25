package builder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mebyus/ku/goku/compiler/srcmap/origin"
)

func TestWalk(t *testing.T) {
	const dir = "testdata"
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Error(err)
		return
	}

	var dirs []string
	for _, entry := range entries {
		dirs = append(dirs, entry.Name())
	}

	for _, d := range dirs {
		t.Run(d, func(t *testing.T) {
			const entname = "entry"

			base := filepath.Join(dir, d)
			entdir := filepath.Join(base, entname)
			entries, err := os.ReadDir(entdir)
			if err != nil {
				t.Error(err)
				return
			}
			if len(entries) == 0 {
				t.Errorf("no units in \"%s\" dir", entname)
				return
			}

			var inits []QueueItem
			for _, entry := range entries {
				inits = append(inits, QueueItem{
					Path: origin.Local(filepath.Join(entname, entry.Name())),
				})
			}

			bundle, err := Walk(WalkConfig{Dir: BaseDirs{Loc: base}}, inits...)
			if err != nil {
				t.Error(err)
				return
			}
			err = CompileBundle(bundle)
			if err != nil {
				t.Error(err)
				return
			}
		})
	}
}
