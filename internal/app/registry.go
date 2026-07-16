package app

import "github.com/dokod-fr/quadboard/internal/domain"

type ProviderRegistry struct {
	providers []domain.Provider
}

func NewProviderRegistry(providers ...domain.Provider) *ProviderRegistry {
	return &ProviderRegistry{
		providers: providers,
	}
}

func (s *ProviderRegistry) Resources() ([]domain.Resource, error) {
	var resources []domain.Resource

	for _, provider := range s.providers {
		r, err := provider.Resources()
		if err != nil {
			return nil, err
		}

		resources = append(resources, r...)
	}

	return resources, nil
}
