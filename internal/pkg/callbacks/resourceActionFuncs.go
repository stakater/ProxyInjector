package callbacks

import (
	"k8s.io/api/apps/v1beta2"
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

func IsDeployment(resource interface{}) bool {
	if _, ok := resource.(*v1beta1.Deployment); ok {
		return true
	}
	return false
}

func IsDeamonset(resource interface{}) bool {
	if _, ok := resource.(*v1beta2.DaemonSet); ok {
		return true
	}
	return false
}

func IsStatefulset(resource interface{}) bool {
	if _, ok := resource.(*v1beta2.StatefulSet); ok {
		return true
	}
	return false
}
