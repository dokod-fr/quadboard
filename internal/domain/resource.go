package domain

type Resource struct {
	ID          string
	Name        string
	Description string
	URL         string
	Icon        string
	Logo        string
	Group       string
	Tags        []string
	Health      HealthStatus
	Actions     []Action
}
