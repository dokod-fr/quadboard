package mock

import "github.com/dokod-fr/quadboard/internal/domain"

func Grafana() domain.Resource {
	return domain.Resource{
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
	}
}
