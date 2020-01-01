package callbacks

import (
	apps "k8s.io/api/apps/v1"
	ext "k8s.io/api/apps/v1"
)

func GetDeploymentAnnotations(resource interface{}) map[string]string {
	return resource.(*ext.Deployment).GetAnnotations()
}

func GetDeploymentName(resource interface{}) string {
	return resource.(*ext.Deployment).Name
}

func GetDeploymentNamespace(resource interface{}) string {
	return resource.(*ext.Deployment).Namespace
}

func GetDaemonsetAnnotations(resource interface{}) map[string]string {
	return resource.(*apps.DaemonSet).GetAnnotations()
}

func GetDaemonsetName(resource interface{}) string {
	return resource.(*apps.DaemonSet).Name
}

func GetDaemonsetNamespace(resource interface{}) string {
	return resource.(*apps.DaemonSet).Namespace
}

func GetStatefulsetAnnotations(resource interface{}) map[string]string {
	return resource.(*apps.StatefulSet).GetAnnotations()
}

func GetStatefulsetName(resource interface{}) string {
	return resource.(*apps.StatefulSet).Name
}

func GetStatefulsetNamespace(resource interface{}) string {
	return resource.(*apps.StatefulSet).Namespace
}
