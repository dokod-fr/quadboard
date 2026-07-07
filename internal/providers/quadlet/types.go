package quadlet

import "path/filepath"

type Kind uint8

const (
	Unknown Kind = iota
	PodKind
	ContainerKind
	VolumeKind
)

func kind(path string) Kind {
	switch filepath.Ext(path) {
	case ".pod":
		return PodKind

	case ".container":
		return ContainerKind

	case ".volume":
		return VolumeKind

	default:
		return Unknown
	}
}
