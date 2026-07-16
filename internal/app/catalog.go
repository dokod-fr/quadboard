package app

import (
	"sync"

	"github.com/dokod-fr/quadboard/internal/domain"
)

type Catalog struct {
	registry  *ProviderRegistry
	mu        sync.RWMutex
	resources []domain.Resource
}

func NewCatalog(registry *ProviderRegistry) *Catalog {
	return &Catalog{
		registry:  registry,
		resources: make([]domain.Resource, 0),
	}
}

// Put resource in cache to keep HTTP request fast
func (c *Catalog) Resources() []domain.Resource {
	c.mu.RLock()
	defer c.mu.RUnlock()

	copied := make([]domain.Resource, len(c.resources))
	copy(copied, c.resources)

	return copied
}

// Keep resource up to date
func (c *Catalog) Refresh() error {
	newResources, err := c.registry.Resources()
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.resources = newResources
	c.mu.Unlock()

	return nil
}
