package handler

import (
	"encoding/json"
	"fmt"
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
	Resource interface{}
}

type ContainerVolumes struct {
	Name      string `json:"name"`
	MountPath string `json:"mountPath"`
}

type Container struct {
	Name         string             `json:"name"`
	Image        string             `json:"image"`
	Args         []string           `json:"args"`
	VolumeMounts []ContainerVolumes `json:"volumeMounts"`
}

type Spec2 struct {
	Containers []Container `json:"containers"`
	Volumes    []Volume    `json:"volumes"`
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

type Volume struct {
	Name      string    `json:"name"`
	ConfigMap ConfigMap `json:"configMap"`
}

type ConfigMap struct {
	Name string `json:"name"`
}

// Handle processes the newly created resource
func (r ResourceCreatedHandler) Handle() error {
	if r.Resource == nil {
		logger.Errorf("Resource creation handler received nil resource")
	} else {
		logger.Infof("Resource created: %+v", r.Resource)
		name := callbacks.GetDeploymentName(r.Resource)
		namespace := callbacks.GetDeploymentNamespace(r.Resource)
		annotations := callbacks.GetDeploymentAnnotations(r.Resource)
		value := annotations[constants.ImageNameAnnotation]

		logger.Info(value)

		if value != "" {

			logger.Infof("Updating deployment ... %s", name)

			containerArgs := []string{
				"--config=" + annotations[constants.ConfigAnnotation],
				"--upstream-url=" + annotations[constants.UpstreamUrlAnnotation],
				"--redirection-url=" + annotations[constants.RedirectionUrlAnnotation],
			}

			if annotations[constants.EnableAuthorizationAnnotation] == "false" {
				logger.Info("authproxy.stakater.com/enable-authorization-header = " + annotations[constants.EnableAuthorizationAnnotation])
				containerArgs = append(containerArgs, "--enable-authorization-header=false")
			} else {
				logger.Info("authproxy.stakater.com/enable-authorization-header != " + annotations[constants.EnableAuthorizationAnnotation])
				containerArgs = append(containerArgs,
					"--upstream-response-header-timeout="+annotations[constants.ResponseHeaderTimeoutAnnotation],
					"--upstream-timeout="+annotations[constants.TimeoutAnnotation],
					"--upstream-keepalive-timeout"+annotations[constants.KeepaliveTimeoutAnnotation],
					"--server-read-timeout"+annotations[constants.ServerReadTimeoutAnnotation],
					"--server-write-timeout"+annotations[constants.ServerWriteTimeoutAnnotation])
			}

			payload := patch{
				Spec: Spec1{
					Tmpl: Template{
						Spec: Spec2{
							Containers: []Container{{
								Name:  "proxy",
								Image: annotations[constants.ImageNameAnnotation] + ":" + annotations[constants.ImageTagAnnotation],
								Args:  containerArgs,
								VolumeMounts: []ContainerVolumes{{
									Name:      "keycloak-proxy-config",
									MountPath: "/etc/config",
								}},
							}},
							Volumes: []Volume{{
								Name: "keycloak-proxy-config",
								ConfigMap: ConfigMap{
									Name: "keycloak-proxy",
								},
							}},
						},
					},
				},
			}

			client, err := kube.GetClient()
			if err == nil {
				payloadBytes, err3 := json.Marshal(payload)

				if err3 == nil {
					deployment, err2 := client.ExtensionsV1beta1().Deployments(namespace).Patch(name, types.StrategicMergePatchType, payloadBytes)

					if err2 == nil {
						logger.Infof("Updated deployment... %s", callbacks.GetDeploymentName(deployment))
					} else {
						logger.Error(err2)
					}

					retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
						// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
						result, getErr := client.CoreV1().Services(namespace).Get(annotations[constants.SourceServiceNameAnnotation], metav1.GetOptions{})
						if getErr != nil {
							panic(fmt.Errorf("Failed to get latest version of Service: %v", getErr))
						}

						result.Spec.Ports[0].TargetPort = intstr.FromInt(80)
						_, updateErr := client.CoreV1().Services(namespace).Update(result)
						return updateErr
					})

					if retryErr == nil {
						logger.Infof("Updated service... %s", annotations[constants.SourceServiceNameAnnotation])
					} else {
						panic(fmt.Errorf("Update failed: %v", retryErr))
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
