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
		slog.Debug("Impossible to read metadata. Ignored",
			slog.String("path", path),
			slog.Any("error", err),
		)
		return nil
	}

	isDir := d.IsDir()
	// Use path as default resolvedPath for regular files
	resolvedPath := path

	if info.Mode()&os.ModeSymlink != 0 {
		resolvedPath, err = filepath.EvalSymlinks(path)
		if err != nil {
			slog.Warn("Link seems broken, ignored", slog.String("path", path))
			return nil
		}
		resolvedInfo, err := os.Stat(resolvedPath)
		if err != nil {
			slog.Warn("Impossible to read resolved link metadata. Ignored", slog.String("path", resolvedPath))
			return nil
		}
		isDir = resolvedInfo.IsDir()
	}

	// Exclude directories, drop-ins, and templates.
	// We check both the original path and the resolved path for templates,
	// so symlinks pointing to a template are also ignored.
	if isDir || isDropIn(path) || isTemplate(path) || isTemplate(resolvedPath) {
		return nil
	}

	// We use the resolved path to get the real extension and base name,
	// as Podman will ultimately read the target file.
	ext := strings.ToLower(filepath.Ext(resolvedPath))
	baseName := filepath.Base(resolvedPath)
	name := strings.TrimSuffix(baseName, ext)

	switch ext {
	case ".pod":
		slog.Info("Pod detected", slog.String("name", name), slog.String("path", path))
		model.Pods = append(model.Pods, Pod{
			Name: name,
			Path: path,
		})

	case ".container":
		slog.Info("Container detected", slog.String("name", name), slog.String("path", path))
		model.Containers = append(model.Containers, Container{
			Name: name,
			Path: path,
		})

	case ".volume":
		slog.Info("Volume detected", slog.String("name", name), slog.String("path", path))
		model.Volumes = append(model.Volumes, Volume{
			Name: name,
			Path: path,
		})
	}

	return nil
}
