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

	p.Name = file.Section("Pod").First("PodName")
	p.Description = file.Section("Unit").First("Description")

	// La méthode Section retourne une map, on peut accéder directement au tableau de labels
	p.Labels = parseLabels(file.Section("Pod")["Label"])

	return nil
}

func parseContainer(c *Container) error {
	file, err := unitfile.ParseFile(c.Path)
	if err != nil {
		return err
	}

	c.Name = file.Section("Container").First("ContainerName")
	c.Description = file.Section("Unit").First("Description")

	if pod := file.Section("Container").First("Pod"); pod != "" {
		c.Pod = strings.TrimSuffix(pod, ".pod")
	}

	c.Labels = parseLabels(file.Section("Container")["Label"])

	return nil
}

// parseLabels convertit un tableau de chaînes "key=value" en map.
// Gère également le cas où la valeur est entourée de guillemets.
func parseLabels(rawLabels []string) map[string]string {
	labels := make(map[string]string)
	for _, l := range rawLabels {
		// On retire les guillemets englobants si présents (ex: Label="key=val")
		l = strings.Trim(l, `"`)

		parts := strings.SplitN(l, "=", 2)
		if len(parts) == 2 {
			labels[parts[0]] = parts[1]
		}
	}
	return labels
}
