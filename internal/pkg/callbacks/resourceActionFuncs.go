package callbacks

import (
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
