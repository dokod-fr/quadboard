package quadlet

import (
	"fmt"
	"log/slog"
	"regexp"
	"sort"
	"strings"

	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/domain"
)

// traefikHostRegex extracts the domain from a Traefik rule: Host(`example.com`)
var traefikHostRegex = regexp.MustCompile("Host\\(`([^`]+)`\\)")

/*
 * ==== Utility functions =======
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

func guessLogoFromImage(image string) string {
	if image == "" {
		return ""
	}

	baseImage := strings.Split(image, ":")[0]
	parts := strings.Split(baseImage, "/")
	appName := parts[len(parts)-1]
	appName = strings.ToLower(appName)

	url := fmt.Sprintf("https://cdn.simpleicons.org/%s", appName)
	slog.Debug("Simple Icons Url", slog.String("url", url))
	return url
}

/* ======= End utility functions ======== */

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

		enrichResource(res, pod.Labels, "")

		resources[pod.Name] = res
	}

	// Standalone containers become Resources.
	// Containers belonging to a Pod enrich the Pod Resource.
	for _, container := range model.Containers {
		if shouldIgnore(container.Name, container.Labels, cfg.ShowItSelf) {
			continue
		}

		// Case: container standalone
		if container.Pod == "" {
			res := &domain.Resource{
				ID:     container.Name,
				Name:   container.Name,
				Health: domain.HealthUnknown,
			}
			res.Description = container.Description
			enrichResource(res, container.Labels, container.Image)

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

		img := ""
		if podRes.Logo == "" && podRes.Icon == "" {
			img = container.Image
		}

		enrichResource(podRes, container.Labels, img)
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
func enrichResource(res *domain.Resource, labels map[string]string, image string) {
	// 1. Group
	if val, ok := labels["quadboard.group"]; ok {
		res.Group = val
	}

	// 2. Icon
	if val, ok := labels["quadboard.icon"]; ok {
		res.Icon = val
	}

	// 3. Logo (Quadboard > Auto-guess via SimpleIcons)
	if val, ok := labels["quadboard.logo"]; ok {
		res.Logo = val
	} else if res.Icon == "" && res.Logo == "" {
		if guessedLogo := guessLogoFromImage(image); guessedLogo != "" {
			res.Logo = guessedLogo
		}
	}

	// 4. Description
	if val, ok := labels["quadboard.description"]; ok {
		res.Description = val
	} else if res.Description == "" {
		if val, ok := labels["io.containers.description"]; ok {
			res.Description = val
		} else if val, ok := labels["org.opencontainers.image.description"]; ok {
			res.Description = val
		}
	}

	// 5. URL
	if val, ok := labels["quadboard.url"]; ok {
		res.URL = val
	} else if res.URL == "" {
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
