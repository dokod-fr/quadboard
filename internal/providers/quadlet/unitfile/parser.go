package unitfile

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func Parse(r io.Reader) (*File, error) {
	file := &File{}

	var section string

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.Trim(line, "[]")
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		parseKeyValue(
			file,
			section,
			strings.TrimSpace(key),
			strings.TrimSpace(value),
		)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return file, nil
}

func ParseFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}
