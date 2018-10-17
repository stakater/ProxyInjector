package handler

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	logger "github.com/sirupsen/logrus"
	"github.com/stakater/ProxyInjector/internal/pkg/callbacks"
	"github.com/stakater/ProxyInjector/internal/pkg/config"
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
	Config   config.Config
	//Config   string
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
		name := callbacks.GetDeploymentName(r.Resource)
		namespace := callbacks.GetDeploymentNamespace(r.Resource)
		annotations := callbacks.GetDeploymentAnnotations(r.Resource)

		if annotations[constants.EnabledAnnotation] == "true" {

			logger.Infof("Updating deployment ... %s", name)

			// TODO Handle config fields through map, rather than struct
			containerArgs := []string{
				"--client-id=" + r.Config.ClientId,
				"--client-secret=" + r.Config.ClientSecret,
				"--discovery-url=" + r.Config.DiscoveryUrl,
				"--enable-default-deny=" + r.Config.EnableDefaultDeny,
				"--listen=" + r.Config.Listen,
				"--secure-cookie=" + r.Config.SecureCookie,
				"--verbose=" + r.Config.Verbose,
				"--enable-logging=" + r.Config.EnableLogging,
				"--config=" + annotations[constants.ConfigAnnotation],
				"--upstream-url=" + annotations[constants.UpstreamUrlAnnotation],
				"--redirection-url=" + annotations[constants.RedirectionUrlAnnotation],
			}

			for _, origin := range r.Config.CorsOrigins {
				containerArgs = append(containerArgs, "--cors-origins="+origin)
			}
			for _, method := range r.Config.CorsMethods {
				containerArgs = append(containerArgs, "--cors-methods="+method)
			}
			for _, resource := range r.Config.Resources {
				containerArgs = append(containerArgs, "--resources=\"uri="+resource.URI+"\"")
			}

			// TODO Handle annotations dynamically instead of being hardcoded
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
								/*								VolumeMounts: []ContainerVolumes{{
																Name:      "keycloak-proxy-config",
																MountPath: "/etc/config",
															}},*/
							}},
							/*							Volumes: []Volume{{
														Name: "keycloak-proxy-config",
														ConfigMap: ConfigMap{
															Name: "keycloak-proxy",
														},
													}},*/
						},
					},
				},
			}

			client, err := kube.GetClient()
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
