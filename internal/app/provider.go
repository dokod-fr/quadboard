package app

import "github.com/dokod-fr/quadboard/internal/domain"

type Provider interface {
	Resources() ([]domain.Resource, error)
}
