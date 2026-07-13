package quadlet

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func Load(paths ...string) (*Model, error) {
	model := &Model{}

	for _, path := range paths {
		if err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			return loadEntry(model, path, d)
		}); err != nil {
			return nil, err
		}
	}

	return model, nil
}

func loadEntry(model *Model, path string, d fs.DirEntry) error {
	if d.IsDir() ||
		isDropIn(path) ||
		isTemplate(path) ||
		isSymlink(path) {
		return nil
	}

	switch kind(path) {

	case PodKind:
		model.Pods = append(model.Pods, Pod{
			Name: strings.TrimSuffix(filepath.Base(path), ".pod"),
			Path: path,
		})

	case ContainerKind:
		model.Containers = append(model.Containers, Container{
			Name: strings.TrimSuffix(filepath.Base(path), ".container"),
			Path: path,
		})

	case VolumeKind:
		model.Volumes = append(model.Volumes, Volume{
			Name: strings.TrimSuffix(filepath.Base(path), ".volume"),
			Path: path,
		})
	}

	return nil
}
