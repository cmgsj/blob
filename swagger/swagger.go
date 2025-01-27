package swagger

import (
	"embed"
	"io/fs"
)

var (
	//go:embed all:dist
	dist embed.FS
)

func Handler() fs.FS {
	distFS, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}

	return distFS
}
