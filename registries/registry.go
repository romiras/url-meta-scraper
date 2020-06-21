package registries

import "github.com/romiras/url-meta-scraper/services"

type Registry struct {
	Fetcher     services.Fetcher
	FetchHelper services.FetchHelper
}

func NewRegistry() *Registry {
	return &Registry{
		services.Fetcher{
			Client: services.NewHTTPClient(),
		},
		services.FetchHelper{},
	}
}
