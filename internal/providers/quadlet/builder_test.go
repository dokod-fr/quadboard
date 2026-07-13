package quadlet

import (
	"fmt"
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

	if len(resources) != 4 {
		t.Fatalf("expected 4 resources, got %d", len(resources))
	}

	fmt.Println("Resources built:")
	for _, resource := range resources {
		fmt.Printf("- %s\n", resource.Name)
	}
	assertContainsResource(t, resources, "dozzle")
	assertContainsResource(t, resources, "nextcloud")
	assertContainsResource(t, resources, "lldap")
	assertContainsResource(t, resources, "traefik")

	assertNotContainsResource(t, resources, "proxy")
	assertNotContainsResource(t, resources, "nextcloud-app")
	assertNotContainsResource(t, resources, "nextcloud-nginx")
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
