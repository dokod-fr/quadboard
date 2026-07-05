package mock

import "github.com/dokod-fr/quadboard/internal/domain"

type Provider struct{}

func New() *Provider {
	return &Provider{}
}

func (p *Provider) Resources() ([]domain.Resource, error) {
	return []domain.Resource{
		{
			ID:          "grafana",
			Name:        "Grafana",
			Description: "Dashboards & Metrics",
			Group:       "Monitoring",
			Health:      domain.HealthHealthy,
			Actions: []domain.Action{
				{
					Label: "Open",
					URL:   "https://grafana.example.com",
				},
			},
		},
	}, nil
}
