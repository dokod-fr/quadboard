package quadlet

import (
	"regexp"
	"sort"
	"strings"

	"github.com/dokod-fr/quadboard/internal/domain"
)

// traefikHostRegex extracts the domain from a Traefik rule: Host(`example.com`)
var traefikHostRegex = regexp.MustCompile("Host\\(`([^`]+)`\\)")

func Build(model *Model) ([]domain.Resource, error) {
	resources := make(map[string]*domain.Resource)

	// Every Pod becomes a Resource.
	for _, pod := range model.Pods {
		res := &domain.Resource{
			ID:     pod.Name,
			Name:   pod.Name,
			Health: domain.HealthUnknown,
		}
		// On utilise la description du Pod par défaut, puis on enrichit avec ses labels
		res.Description = pod.Description
		enrichResource(res, pod.Labels)

		resources[pod.Name] = res
	}

	// Standalone containers become Resources.
	// Containers belonging to a Pod enrich the Pod Resource.
	for _, container := range model.Containers {
		if container.Pod == "" {
			res := &domain.Resource{
				ID:     container.Name,
				Name:   container.Name,
				Health: domain.HealthUnknown,
			}
			res.Description = container.Description
			enrichResource(res, container.Labels)

			resources[container.Name] = res
			continue
		}

		// Pod not found: ignore for now.
		podRes, ok := resources[container.Pod]
		if !ok {
			continue
		}

		// Si la description du Pod est vide, on prend celle du conteneur
		if podRes.Description == "" && container.Description != "" {
			podRes.Description = container.Description
		}

		// On enrichit le Pod avec les labels du conteneur (ex: Traefik est souvent sur le conteneur)
		enrichResource(podRes, container.Labels)
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

// enrichResource met à jour une ressource en se basant sur les labels fournis.
// Les labels quadboard.* ont la priorité absolue.
func enrichResource(res *domain.Resource, labels map[string]string) {
	if len(labels) == 0 {
		return
	}

	// Group
	if val, ok := labels["quadboard.group"]; ok {
		res.Group = val
	}

	// Icon / Logo
	if val, ok := labels["quadboard.icon"]; ok {
		res.Icon = val
	}

	if val, ok := labels["quadboard.logo"]; ok {
		res.Logo = val
	}

	// Description (Quadboard override > Podman standard)
	if val, ok := labels["quadboard.description"]; ok {
		res.Description = val
	} else if res.Description == "" {
		if val, ok := labels["io.containers.description"]; ok {
			res.Description = val
		} else if val, ok := labels["org.opencontainers.image.description"]; ok {
			res.Description = val
		}
	}

	// URL
	if val, ok := labels["quadboard.url"]; ok {
		res.URL = val
	} else if res.URL == "" {
		// Recherche d'une règle Traefik si l'URL n'est pas déjà définie
		for key, value := range labels {
			if strings.HasPrefix(key, "traefik.http.routers.") && strings.HasSuffix(key, ".rule") {
				matches := traefikHostRegex.FindStringSubmatch(value)
				if len(matches) > 1 {
					res.URL = "https://" + matches[1]
					break
				}
			}
		}
	}
}
