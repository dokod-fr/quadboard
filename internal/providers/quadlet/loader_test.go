package quadlet

import (
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	model, err := Load(filepath.Join("testdata", "discovery"))
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if len(model.Pods) != 1 {
		t.Fatalf("expected 1 pod, got %d", len(model.Pods))
	}

	if len(model.Containers) != 6 {
		t.Fatalf("expected 6 containers, got %d", len(model.Containers))
	}

	if len(model.Volumes) != 5 {
		t.Fatalf("expected 5 volumes, got %d", len(model.Volumes))
	}
}

func TestLoadNames(t *testing.T) {
	model, err := Load(filepath.Join("testdata", "discovery"))
	if err != nil {
		t.Fatal(err)
	}

	assertContainsPod(t, model, "nextcloud")

	assertContainsContainer(t, model, "authelia")
	assertContainsContainer(t, model, "lldap")
	assertContainsContainer(t, model, "traefik")
	assertContainsContainer(t, model, "dozzle")
	assertContainsContainer(t, model, "nextcloud-app")
	assertContainsContainer(t, model, "nextcloud-nginx")
}

func TestIgnoredTemplateContainer(t *testing.T) {
	model, err := Load(filepath.Join("testdata", "discovery"))
	if err != nil {
		t.Fatal(err)
	}

	assertNotContainsContainer(t, model, "redis@")
}

// ==== Helper ====

func assertContainsContainer(t *testing.T, model *Model, name string) {
	t.Helper()

	for _, container := range model.Containers {
		if container.Name == name {
			return
		}
	}

	t.Fatalf("container %q not found", name)
}

func assertNotContainsContainer(t *testing.T, model *Model, name string) {
	t.Helper()

	for _, container := range model.Containers {
		if container.Name == name {
			t.Fatalf("container %q should have been ignored", name)
		}
	}
}

func assertContainsPod(t *testing.T, model *Model, name string) {
	t.Helper()

	for _, pod := range model.Pods {
		if pod.Name == name {
			return
		}
	}

	t.Fatalf("pod %q not found", name)
}
