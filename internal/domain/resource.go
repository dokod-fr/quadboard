package domain

type Resource struct {
	ID          string
	Name        string
	Description string

	Group string
	Tags  []string

	Health HealthStatus

	Actions []Action
}
