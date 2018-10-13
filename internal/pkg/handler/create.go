package handler

import (
	logger "github.com/sirupsen/logrus"
	"github.com/stakater/ProxyInjector/internal/pkg/callbacks"
	"github.com/stakater/ProxyInjector/internal/pkg/constants"
	"github.com/stakater/ProxyInjector/pkg/kube"
	//"k8s.io/client-go/util/retry"
	"encoding/json"
	"k8s.io/apimachinery/pkg/types"
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
	Name         string           `json:"name"`
	Image        string           `json:"image"`
	VolumeMounts ContainerVolumes `json: "volumeMounts"`
}

/*type patch struct {
	Op    string    `json:"op"`
	Path  string    `json:"path"`
	Value Container `json: "value"`
}*/

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
func (r ResourceCreatedHandler) Handle() error {
	if r.Resource == nil {
		logger.Errorf("Resource creation handler received nil resource")
	} else {
		logger.Info("Resource created")
		name := callbacks.GetDeploymentName(r.Resource)
		namespace := callbacks.GetDeploymentNamespace(r.Resource)
		annotations := callbacks.GetDeploymentAnnotations(r.Resource)
		value := annotations[constants.AuthProxyUpstreamAnnotation]

		logger.Info(value)

		if value != "" {

			logger.Infof("Updating deployment ... %s", name)
			/*payload := patch{
				Op:   "add",
				Path: "/spec/template/spec/containers",
				Value: Container{
					Name:  "proxy",
					Image: "quay.io/gambol99/keycloak-proxy:v2.1.1",
					VolumeMounts: ContainerVolumes{
						Name:      "keycloak-proxy-config",
						MountPath: "/etc/config",
					},
				},
			}*/

			/*payload := Container{
				Name:  "proxy",
				Image: "quay.io/gambol99/keycloak-proxy:v2.1.1",
				VolumeMounts: ContainerVolumes{
					Name:      "keycloak-proxy-config",
					MountPath: "/etc/config",
				},
			}*/

			payload := patch{
				Spec: Spec1{
					Tmpl: Template{
						Spec: Spec2{
							Containers: []Container{{
								Name:  "proxy",
								Image: "quay.io/gambol99/keycloak-proxy:v2.1.1",
								VolumeMounts: ContainerVolumes{
									Name:      "keycloak-proxy-config",
									MountPath: "/etc/config",
								},
							}},
						},
					},
				},
			}

			client, err := kube.GetClient()
			if err == nil {
				payloadBytes, err3 := json.Marshal(payload)
				//deployment, err2 := client.ExtensionsV1beta1().Deployments(namespace).Patch(name, types.StrategicMergePatchType, valueBytes,"/spec/template/spec/containers")

				if err3 == nil {
					deployment, err2 := client.ExtensionsV1beta1().Deployments(namespace).Patch(name, types.StrategicMergePatchType, payloadBytes)

					if err2 == nil {
						logger.Infof("Updated deployment... %s", callbacks.GetDeploymentName(deployment))
					} else {
						logger.Error(err2)
					}
				} else {
					logger.Error(err3)
				}
			} else {
				logger.Error(err)
			}

			/*			&v1.Container{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						}
			*/

		}
	}
	return nil
}
