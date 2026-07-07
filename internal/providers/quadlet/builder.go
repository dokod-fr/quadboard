package quadlet

import (
	"sort"

	"github.com/dokod-fr/quadboard/internal/domain"
)

func Build(model *Model) ([]domain.Resource, error) {
	resources := make(map[string]*domain.Resource)

	// Every Pod becomes a Resource.
	for _, pod := range model.Pods {
		resources[pod.Name] = &domain.Resource{
			ID:          pod.Name,
			Name:        pod.Name,
			Description: pod.Description,
			Health:      domain.HealthUnknown,
		}
	}

	// Standalone containers become Resources.
	// Containers belonging to a Pod enrich the Pod Resource.
	for _, container := range model.Containers {
		if container.Pod == "" {
			resources[container.Name] = &domain.Resource{
				ID:          container.Name,
				Name:        container.Name,
				Description: container.Description,
				Health:      domain.HealthUnknown,
			}
			continue
		}

		// Pod not found: ignore for now.
		if _, ok := resources[container.Pod]; !ok {
			continue
		}

		// Nothing to enrich yet.
	}

	list := make([]domain.Resource, 0, len(resources))

	for _, resource := range resources {
		list = append(list, *resource)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}
