package quadlet

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dokod-fr/quadboard/internal/config"
	"github.com/dokod-fr/quadboard/internal/domain"
)

var traefikHostRegex = regexp.MustCompile("Host\\(`([^`]+)`\\)")

// --- Logo cache SimpleIcons ---
var (
	simpleIconsCache = make(map[string]bool)
	cacheMutex       sync.Mutex
)

/*
 * ==== Utility functions =======
 */

// shouldIgnore determine if container or pod should be ignored
func shouldIgnore(name string, labels map[string]string) bool {
	// Ignore Quadlet templates  mean container which name least with '@'
	if strings.Contains(name, "@") {
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

	generics := map[string]bool{
		"server": true, "app": true, "web": true, "api": true,
		"latest": true, "cli": true, "core": true, "main": true,
	}
	if generics[appName] && len(parts) > 1 {
		appName = parts[len(parts)-2]
	}
	appName = strings.ToLower(appName)

	cacheMutex.Lock()
	isValid, exists := simpleIconsCache[appName]
	cacheMutex.Unlock()

	if !exists {
		url := fmt.Sprintf("https://cdn.simpleicons.org/%s", appName)
		client := http.Client{Timeout: 3 * time.Second}
		resp, err := client.Head(url)

		if err == nil {
			resp.Body.Close()
			isValid = resp.StatusCode == http.StatusOK
		} else {
			isValid = false
		}

		cacheMutex.Lock()
		simpleIconsCache[appName] = isValid
		cacheMutex.Unlock()
	}

	if isValid {
		return fmt.Sprintf("https://cdn.simpleicons.org/%s", appName)
	}

	return ""
}

/* ======= End utility functions ======== */

func Build(model *Model, cfg *config.QuadletConfig) ([]domain.Resource, error) {
	resources := make(map[string]*domain.Resource)

	// Every Pod becomes a Resource.
	for _, pod := range model.Pods {
		if shouldIgnore(pod.Name, pod.Labels) {
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
		if shouldIgnore(container.Name, container.Labels) {
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
