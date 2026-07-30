package view

import "embed"

//go:embed index.html assets
var files embed.FS

// FS retourne le système de fichiers embarqué contenant les templates et les assets.
func FS() embed.FS {
	return files
}
