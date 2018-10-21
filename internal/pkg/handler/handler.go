package handler

import "github.com/stakater/ProxyInjector/internal/pkg/config"

// ResourceHandler handles the creation and update of resources
type ResourceHandler interface {
	Handle(conf config.Config) error
}
