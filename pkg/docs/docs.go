package docs

import (
	"embed"
	"io/fs"
)

//go:embed all:assets
var assets embed.FS

func Assets() fs.FS {
	fsys, err := fs.Sub(assets, "assets")
	if err != nil {
		panic(err)
	}

	return fsys
}
