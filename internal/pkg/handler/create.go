package handler

import (
	logger "github.com/sirupsen/logrus"
	"github.com/stakater/ProxyInjector/internal/pkg/callbacks"
	"github.com/stakater/ProxyInjector/internal/pkg/constants"
)

// ResourceCreatedHandler contains new objects
type ResourceCreatedHandler struct {
	Resource interface{}
}

// Handle processes the newly created resource
func (r ResourceCreatedHandler) Handle() error {
	if r.Resource == nil {
		logger.Errorf("Resource creation handler received nil resource")
	} else {
		logger.Info("Resource created")
		//client, err := kube.GetClient()
		annotations := callbacks.GetDeploymentAnnotations(r.Resource)
		value := annotations[constants.AuthProxyUpstreamAnnotation]
		logger.Info(value)

	}
	return nil
}
