package domain

type HealthStatus string

const (
	HealthUnknown     HealthStatus = "unknown"
	HealthHealthy     HealthStatus = "healthy"
	HealthStarting    HealthStatus = "starting"
	HealthUnreachable HealthStatus = "unreachable"
)
