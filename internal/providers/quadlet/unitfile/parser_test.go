package unitfile

import (
	"strings"
	"testing"
)

func TestParseContainer(t *testing.T) {
	input := `
[Unit]
Description=Dozzle logging

[Container]
ContainerName=dozzle
Image=docker.io/amir20/dozzle:v10

Environment=A=1
Environment=B=2
`

	file, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}

	if got := file.Unit.Description; got != "Dozzle logging" {
		t.Fatalf("got %q", got)
	}

	if got := file.Container.ContainerName; got != "dozzle" {
		t.Fatalf("got %q", got)
	}

	if got := len(file.Container.Environment); got != 2 {
		t.Fatalf("got %d", got)
	}
}
