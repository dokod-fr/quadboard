package app_test

import (
	"testing"

	"github.com/dokod-fr/quadboard/internal/app"
	"github.com/dokod-fr/quadboard/internal/domain"
)

type MockProvider struct {
	mockResources []domain.Resource
}

func (m *MockProvider) Resources() ([]domain.Resource, error) {
	return m.mockResources, nil
}

func TestCatalog_RefreshAndResources(t *testing.T) {
	mockResource := domain.Resource{ /* configure un mock minimaliste ici */ }
	provider := &MockProvider{mockResources: []domain.Resource{mockResource}}

	registry := app.NewProviderRegistry(provider)
	catalog := app.NewCatalog(registry)

	// Clean up the catalog
	if len(catalog.Resources()) != 0 {
		t.Fatal("Catalog should be empty at initial state")
	}

	err := catalog.Refresh()
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}

	res := catalog.Resources()
	if len(res) != 1 {
		t.Fatalf("Expected 1 resource got %d", len(res))
	}
}
