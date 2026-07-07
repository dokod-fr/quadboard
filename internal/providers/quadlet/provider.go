package quadlet

import "github.com/dokod-fr/quadboard/internal/domain"

type Provider struct {
	paths []string
}

func New(paths ...string) *Provider {
	return &Provider{
		paths: paths,
	}
}

func (p *Provider) Resources() ([]domain.Resource, error) {
	model, err := Load(p.paths...)
	if err != nil {
		return nil, err
	}

	if err := Parse(model); err != nil {
		return nil, err
	}

	return Build(model)
}
