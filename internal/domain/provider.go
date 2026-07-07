package domain

type Provider interface {
	Resources() ([]Resource, error)
}
