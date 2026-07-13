package app

import "github.com/dokod-fr/quadboard/internal/domain"

type Discovery struct {
	providers []domain.Provider
}

func NewDiscovery(providers ...domain.Provider) *Discovery {
	return &Discovery{
		providers: providers,
	}
}

func (s *Discovery) Resources() ([]domain.Resource, error) {
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
