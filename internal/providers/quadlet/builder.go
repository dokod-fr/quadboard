package quadlet

import (
	"regexp"
	"sort"
	"strings"

	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/domain"
)

// traefikHostRegex extracts the domain from a Traefik rule: Host(`example.com`)
var traefikHostRegex = regexp.MustCompile("Host\\(`([^`]+)`\\)")

/*
 * ==== Utility functions
 */

// shouldIgnore determine if container or pod should be ignored
func shouldIgnore(name string, labels map[string]string, showItSelf bool) bool {
	// Ignore Quadlet templates  mean container which name least with '@'
	if strings.Contains(name, "@") {
		return true
	}

	// Ignore QuadBoard container/pod if asking
	isQuadboard := strings.EqualFold(name, "quadboard") || strings.Contains(strings.ToLower(name), "quadboard")
	if isQuadboard && !showItSelf {
		return true
	}

	return false
}

func Build(model *Model, cfg *config.QuadletConfig) ([]domain.Resource, error) {
	resources := make(map[string]*domain.Resource)

	// Every Pod becomes a Resource.
	for _, pod := range model.Pods {
		if shouldIgnore(pod.Name, pod.Labels, cfg.ShowItSelf) {
			continue
		}

		res := &domain.Resource{
			ID:     pod.Name,
			Name:   pod.Name,
			Health: domain.HealthUnknown,
		}
		// Description comes from pod first
		res.Description = pod.Description
		enrichResource(res, pod.Labels)

		resources[pod.Name] = res
	}

	// Standalone containers become Resources.
	// Containers belonging to a Pod enrich the Pod Resource.
	for _, container := range model.Containers {
		if shouldIgnore(container.Name, container.Labels, cfg.ShowItSelf) {
			continue
		}

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

		// If not found in Pod, it would be in container
		if podRes.Description == "" && container.Description != "" {
			podRes.Description = container.Description
		}

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

// Use Quadboard labels (first class citizen)
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
		// Use Traefik URL if presents
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
