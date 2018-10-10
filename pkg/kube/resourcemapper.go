package kube

import (
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

// ResourceMap are resources from where changes are going to be detected
var ResourceMap = map[string]runtime.Object{
	"deployments": &v1beta1.Deployment{},
}
