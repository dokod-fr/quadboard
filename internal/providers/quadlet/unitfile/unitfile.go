package unitfile

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Section represents a section like [Container] or [Unit].
// It maps keys to a slice of values to support repeated keys (e.g., Label=).
type Section map[string][]string

// UnitFile represents a parsed systemd/quadlet file.
type UnitFile map[string]Section

// Parse reads an io.Reader and returns a generic UnitFile.
func Parse(r io.Reader) (UnitFile, error) {
	file := make(UnitFile)
	var currentSection Section

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.Trim(line, "[]")
			currentSection = make(Section)
			file[sectionName] = currentSection
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok || currentSection == nil {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		currentSection[key] = append(currentSection[key], value)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return file, nil
}

// ParseFile opens a file from the filesystem and parses it.
func ParseFile(path string) (UnitFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

func (u UnitFile) Section(name string) Section {
	if s, ok := u[name]; ok {
		return s
	}
	return make(Section)
}

// First is a helper method on Section to get the first value of a key.
func (s Section) First(key string) string {
	if values, ok := s[key]; ok && len(values) > 0 {
		return values[0]
	}
	return ""
}
