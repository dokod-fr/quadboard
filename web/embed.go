package web

import (
	"embed"
	"io/fs"
)

//go:embed assets
var Assets embed.FS

func FS() fs.FS {
	sub, _ := fs.Sub(Assets, "assets")
	return sub
}
