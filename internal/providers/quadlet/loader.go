package quadlet

import (
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

func Load(paths ...string) (*Model, error) {

	model := &Model{}

	for _, path := range paths {
		// Check directory existence before continue
		if _, err := os.Stat(path); os.IsNotExist(err) {
			slog.Warn("Directory in configuration not found, ignored",
				slog.String("path", path),
			)
			continue
		}

		slog.Debug("Walk through directory", slog.String("path", path))

		err := filepath.WalkDir(path, func(currentPath string, d fs.DirEntry, err error) error {
			if err != nil {
				slog.Warn("Nothing found in this directory",
					slog.String("path", currentPath),
					slog.Any("error", err),
				)
				return filepath.SkipDir
			}

			return loadEntry(model, currentPath, d)
		})

		if err != nil {
			return nil, fmt.Errorf("Critical error while walking through %s: %w", path, err)
		}
	}

	return model, nil
}

func loadEntry(model *Model, path string, d fs.DirEntry) error {
	info, err := d.Info()
	if err != nil {
		slog.Debug("Impossible to raad metadata. Ignored",
			slog.String("path", path),
			slog.Any("error", err),
		)
		return nil
	}

	isDir := d.IsDir()
	if info.Mode()&os.ModeSymlink != 0 {
		resolvedPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			slog.Debug("Link seems broken, ignored", slog.String("path", path))
			return nil
		}
		resolvedInfo, err := os.Stat(resolvedPath)
		if err != nil {
			return nil
		}
		isDir = resolvedInfo.IsDir()
	}

	if isDir || isDropIn(path) || isTemplate(path) {
		return nil
	}

	ext := strings.ToLower(filepath.Ext(path))
	baseName := filepath.Base(path)

	switch ext {
	case ".pod":
		slog.Debug("Pod detected", slog.String("path", path))
		model.Pods = append(model.Pods, Pod{
			Name: strings.TrimSuffix(baseName, ext),
			Path: path,
		})

	case ".container":
		slog.Debug("Container detected", slog.String("path", path))
		model.Containers = append(model.Containers, Container{
			Name: strings.TrimSuffix(baseName, ext),
			Path: path,
		})

	case ".volume":
		slog.Debug("Volume detected", slog.String("path", path))
		model.Volumes = append(model.Volumes, Volume{
			Name: strings.TrimSuffix(baseName, ext),
			Path: path,
		})
	}

	return nil
}
