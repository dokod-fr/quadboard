package quadlet

import (
	"path/filepath"
	"testing"
)

func TestParse(t *testing.T) {
	model, err := Load(filepath.Join("testdata", "discovery"))
	if err != nil {
		t.Fatal(err)
	}

	if err := Parse(model); err != nil {
		t.Fatal(err)
	}

	assertContainsPod(t, model, "nextcloud")
	assertContainsContainer(t, model, "nextcloud")
	assertContainsContainer(t, model, "dozzle")

	// Vérifications des métadonnées...
	container := findContainer(t, model, "nextcloud")

	if container.Pod != "nextcloud" {
		t.Fatalf("pod = %q, want %q", container.Pod, "nextcloud")
	}

	container = findContainer(t, model, "dozzle")

	if container.Pod != "" {
		t.Fatalf("expected no pod")
	}
}

func findContainer(t *testing.T, model *Model, name string) *Container {
	t.Helper()

	for i := range model.Containers {
		if model.Containers[i].Name == name {
			return &model.Containers[i]
		}
	}

	t.Fatalf("container %q not found", name)
	return nil
}
