package quadlet

import (
	"strings"

	"github.com/dokod-fr/quadboard/internal/providers/quadlet/unitfile"
)

func Parse(model *Model) error {
	for i := range model.Pods {
		if err := parsePod(&model.Pods[i]); err != nil {
			return err
		}
	}

	for i := range model.Containers {
		if err := parseContainer(&model.Containers[i]); err != nil {
			return err
		}
	}

	return nil
}

func parsePod(p *Pod) error {
	file, err := unitfile.ParseFile(p.Path)
	if err != nil {
		return err
	}

	if name := file.Pod.PodName; name != "" {
		p.Name = name
	}

	if desc := file.Unit.Description; desc != "" {
		p.Description = desc
	}

	return nil
}

func parseContainer(c *Container) error {
	file, err := unitfile.ParseFile(c.Path)
	if err != nil {
		return err
	}

	if name := file.Container.ContainerName; name != "" {
		c.Name = name
	}

	if desc := file.Unit.Description; desc != "" {
		c.Description = desc
	}

	if pod := file.Container.Pod; pod != "" {
		c.Pod = strings.TrimSuffix(pod, ".pod")
	}

	return nil
}
