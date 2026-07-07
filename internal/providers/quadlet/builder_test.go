package quadlet

import (
	"testing"

	"github.com/dokod-fr/quadboard/internal/domain"
)

func TestBuild(t *testing.T) {
	model, err := Load("testdata")
	if err != nil {
		t.Fatal(err)
	}

	if err := Parse(model); err != nil {
		t.Fatal(err)
	}

	resources, err := Build(model)
	if err != nil {
		t.Fatal(err)
	}

	if len(resources) != 3 {
		t.Fatalf("expected 3 resources, got %d", len(resources))
	}

	assertContainsResource(t, resources, "dozzle")
	assertContainsResource(t, resources, "nextcloud")
	assertContainsResource(t, resources, "proxy")

	assertNotContainsResource(t, resources, "nextcloud-app")
	assertNotContainsResource(t, resources, "nextcloud-nginx")
	assertNotContainsResource(t, resources, "authelia")
	assertNotContainsResource(t, resources, "lldap")
	assertNotContainsResource(t, resources, "traefik")
}

func assertContainsResource(t *testing.T, resources []domain.Resource, name string) {
	t.Helper()

	for _, resource := range resources {
		if resource.Name == name {
			return
		}
	}

	t.Fatalf("resource %q not found", name)
}

func assertNotContainsResource(t *testing.T, resources []domain.Resource, name string) {
	t.Helper()

	for _, resource := range resources {
		if resource.Name == name {
			t.Fatalf("unexpected resource %q", name)
		}
	}
}
