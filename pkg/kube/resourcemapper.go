package kube

import (
	ext "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	//apps "k8s.io/api/apps/v1beta2"
)

// ResourceMap are resources from where changes are going to be detected
var ResourceMap = map[string]runtime.Object{
	"deployments": &ext.Deployment{},
	//"daemonsets": &apps.DaemonSet{},
	//"statefulsets": &apps.StatefulSet{},
}
