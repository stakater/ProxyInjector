package handler

import (
	"encoding/json"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/stakater/ProxyInjector/internal/pkg/callbacks"
	"github.com/stakater/ProxyInjector/internal/pkg/constants"
	"github.com/stakater/ProxyInjector/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/util/retry"
)

// ResourceCreatedHandler contains new objects
type ResourceCreatedHandler struct {
	Resource interface{} `json:"resource"`
}

type Container struct {
	Name  string   `json:"name"`
	Image string   `json:"image"`
	Args  []string `json:"args"`
}

type Spec2 struct {
	Containers []Container `json:"containers"`
}

type Template struct {
	Spec Spec2 `json:"spec"`
}

type Spec1 struct {
	Tmpl Template `json:"template"`
}

type patch struct {
	Spec Spec1 `json:"spec"`
}

// Handle processes the newly created resource
func (r ResourceCreatedHandler) Handle(conf []string) error {
	if r.Resource == nil {
		logger.Errorf("Resource creation handler received nil resource")
	} else {
		name := callbacks.GetDeploymentName(r.Resource)
		namespace := callbacks.GetDeploymentNamespace(r.Resource)
		annotations := callbacks.GetDeploymentAnnotations(r.Resource)

		if annotations[constants.EnabledAnnotation] == "true" {

			client, err := kube.GetClient()

			logger.Infof("Updating deployment ... %s", name)

			containerArgs := conf

			for _, arg := range constants.KeycloakArgs {
				if annotations[constants.AnnotationPrefix+arg] != "" {
					containerArgs = append(containerArgs, "--"+arg+"="+annotations[constants.AnnotationPrefix+arg])
				}
			}

			payload := patch{
				Spec: Spec1{
					Tmpl: Template{
						Spec: Spec2{
							Containers: []Container{{
								Name:  "proxy",
								Image: annotations[constants.ImageNameAnnotation] + ":" + annotations[constants.ImageTagAnnotation],
								Args:  containerArgs,
							}},
						},
					},
				},
			}

			if err == nil {
				payloadBytes, err3 := json.Marshal(payload)

				if err3 == nil {

					var err2 error
					logger.Info("checking resource type and updating...")
					if callbacks.IsDeployment(r.Resource) {
						logger.Info("resource is a deployment")
						_, err2 = client.ExtensionsV1beta1().Deployments(namespace).Patch(name, types.StrategicMergePatchType, payloadBytes)
					} else if callbacks.IsDeamonset(r.Resource) {
						logger.Info("resource is a daemonset")
						_, err2 = client.AppsV1beta2().DaemonSets(namespace).Patch(name, types.StrategicMergePatchType, payloadBytes)
					} else if callbacks.IsStatefulset(r.Resource) {
						logger.Info("resource is a statefulset")
						_, err2 = client.AppsV1beta2().StatefulSets(namespace).Patch(name, types.StrategicMergePatchType, payloadBytes)
					} else {
						return errors.New("unexpected resource type")
					}

					if err2 == nil {
						logger.Infof("Updated deployment... %s", name)
						client.ExtensionsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
					} else {
						logger.Error(err2)
					}

					retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
						// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
						result, getErr := client.CoreV1().Services(namespace).Get(annotations[constants.SourceServiceNameAnnotation], metav1.GetOptions{})
						if getErr != nil {
							logger.Errorf("Failed to get latest version of Service: %v", getErr)
						}

						result.Spec.Ports[0].TargetPort = intstr.FromInt(80)
						_, updateErr := client.CoreV1().Services(namespace).Update(result)
						return updateErr
					})

					if retryErr == nil {
						logger.Infof("Updated service... %s", annotations[constants.SourceServiceNameAnnotation])
					} else {
						logger.Errorf("Update failed: %v", retryErr)
					}

				} else {
					logger.Error(err3)
				}
			} else {
				logger.Error(err)
			}

		}
	}
	return nil
}
