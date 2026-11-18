package main

import (
	"github.com/mebyus/ku/internal/ku/builder"
	"github.com/mebyus/ku/internal/ku/builder/kub"
)

func build(path string) error {
	r := builder.Build(&builder.Config{
		Unit: path,
		Build: kub.Build{
			Target: kub.Target{
				Arch: kub.ArchAMD64,
				OS:   kub.Linux,
			},
		},
	})
	return r.Error
}
