package callbacks

import (
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
)

func GetDeploymentAnnotations(resource interface{}) map[string]string {
	return resource.(*v1beta1.Deployment).GetAnnotations()
}

// GetDeploymentPods returns the pods of given deployment
func GetDeploymentPods(item interface{}) []v1.Pod {
	return item.(v1beta1.Deployment).Spec.Template
}