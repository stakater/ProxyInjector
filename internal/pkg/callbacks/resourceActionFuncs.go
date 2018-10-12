package callbacks

import (
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

func GetDeploymentAnnotations(resource interface{}) map[string]string {
	return resource.(*v1beta1.Deployment).GetAnnotations()
}

func GetDeploymentName(resource interface{}) string {
	return resource.(*v1beta1.Deployment).Name
}

func GetDeploymentNamespace(resource interface{}) string {
	return resource.(*v1beta1.Deployment).Namespace
}

// GetDeploymentContainers returns the containers of given deployment
func GetDeploymentContainers(item interface{}) []v1.Container {
	return item.(v1beta1.Deployment).Spec.Template.Spec.Containers
}
