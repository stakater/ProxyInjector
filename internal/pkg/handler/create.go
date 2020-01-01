package handler

import (
	"encoding/json"
	"errors"
	"strings"

	logger "github.com/sirupsen/logrus"
	"github.com/stakater/ProxyInjector/internal/pkg/callbacks"
	"github.com/stakater/ProxyInjector/internal/pkg/config"
	"github.com/stakater/ProxyInjector/internal/pkg/constants"
	"github.com/stakater/ProxyInjector/pkg/kube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
)

// ResourceCreatedHandler contains new objects
type ResourceCreatedHandler struct {
	Resource interface{} `json:"resource"`
}

type patch struct {
	Spec struct {
		Template struct {
			Spec struct {
				Containers []Container `json:"containers"`
			} `json:"spec"`
		} `json:"template"`
	} `json:"spec"`
}

type Container struct {
	Name  string   `json:"name"`
	Image string   `json:"image"`
	Args  []string `json:"args"`
}

// Handle processes the newly created resource
func (r ResourceCreatedHandler) Handle(conf config.Config, resourceType string) error {
	if r.Resource == nil {
		logger.Errorf("Resource creation handler received nil resource")
	} else {

		var name string
		var namespace string
		var annotations map[string]string

		if resourceType == "deployments" {
			name = callbacks.GetDeploymentName(r.Resource)
			namespace = callbacks.GetDeploymentNamespace(r.Resource)
			annotations = callbacks.GetDeploymentAnnotations(r.Resource)
		} else if resourceType == "daemonsets" {
			name = callbacks.GetDaemonsetName(r.Resource)
			namespace = callbacks.GetDaemonsetNamespace(r.Resource)
			annotations = callbacks.GetDaemonsetAnnotations(r.Resource)
		} else if resourceType == "statefulsets" {
			name = callbacks.GetStatefulsetName(r.Resource)
			namespace = callbacks.GetStatefulsetNamespace(r.Resource)
			annotations = callbacks.GetStatefulsetAnnotations(r.Resource)
		}
		logger.Infof("Resource creation handler checking resource %s of type %s in namespace %s", name, resourceType, namespace)

		if annotations[constants.EnabledAnnotation] == "true" {

			client, err := kube.GetKubernetesClient()

			logger.Infof("Updating resource ... %s", name)

			containerArgs := getConfigArgs(conf, annotations)

			for _, arg := range constants.KeycloakArgs {
				if ContainsKey(annotations, arg) {
					containerArgs = removeIfExists(containerArgs, arg)
					resourceStrings := strings.Split(annotations[constants.AnnotationPrefix+arg], "&")

					for _, resourceString := range resourceStrings {
						containerArgs = append(containerArgs, "--"+arg+"="+resourceString)
					}
				} else if annotations[constants.AnnotationPrefix+arg] != "" {
					containerArgs = removeIfExists(containerArgs, arg)
					containerArgs = append(containerArgs, "--"+arg+"="+annotations[constants.AnnotationPrefix+arg])
				}
			}

			if annotations[constants.ImageNameAnnotation] == "" {
				annotations[constants.ImageNameAnnotation] = conf.GatekeeperImage
			}

			if err == nil {
				payloadBytes, err3 := getPatch(containerArgs, annotations[constants.ImageNameAnnotation])

				if err3 == nil {

					var err2 error
					logger.Info("checking resource type and updating...")
					if resourceType == "deployments" {
						logger.Info("patching deployment")
						_, err2 = client.AppsV1().Deployments(namespace).Patch(name, types.StrategicMergePatchType, payloadBytes)
					} else if resourceType == "daemonsets" {
						logger.Info("patching daemonset")
						_, err2 = client.AppsV1().DaemonSets(namespace).Patch(name, types.StrategicMergePatchType, payloadBytes)
					} else if resourceType == "statefulsets" {
						logger.Info("patching statefulset")
						_, err2 = client.AppsV1().StatefulSets(namespace).Patch(name, types.StrategicMergePatchType, payloadBytes)
					} else {
						return errors.New("unexpected resource type")
					}

					if err2 == nil {
						logger.Infof("Updated resource... %s", name)
					} else {
						logger.Error(err2)
					}

					updateService(client, namespace, annotations[constants.SourceServiceNameAnnotation], annotations[constants.TargetPortAnnotation])

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

func removeIfExists(containerArgs []string, arg string) []string {
	i := 0 // output index
	for _, containerArg := range containerArgs {
		if !strings.Contains(containerArg, arg) {
			containerArgs[i] = containerArg
			i++
		}
	}
	return containerArgs[:i]
}

// ContainsKey tells whether a key exist in map[string]string.
func ContainsKey(list map[string]string, word string) bool {
	for key := range list {
		if strings.Contains(key, word) {
			return true
		}
	}
	return false
}

func getPatch(containerArgs []string, image string) ([]byte, error) {

	payload := &patch{}
	payload.Spec.Template.Spec.Containers = []Container{{
		Name:  "proxy",
		Image: image,
		Args:  containerArgs,
	}}

	return json.Marshal(payload)
}

func updateService(client *kubernetes.Clientset, namespace string, service string, port string) {

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := client.CoreV1().Services(namespace).Get(service, metav1.GetOptions{})
		if getErr != nil {
			logger.Errorf("Failed to get latest version of Service: %v", getErr)
		}

		if port == "" {
			result.Spec.Ports[0].TargetPort = intstr.FromInt(80)
		} else {
			result.Spec.Ports[0].TargetPort = intstr.Parse(port)
		}
		_, updateErr := client.CoreV1().Services(namespace).Update(result)
		return updateErr
	})

	if retryErr == nil {
		logger.Infof("Updated service... %s", service)
	} else {
		logger.Errorf("Update failed: %v", retryErr)
	}
}

func getConfigArgs(config config.Config, annotations map[string]string) []string {

	configArgs := []string{}

	//config from global proxy injector config is injected only if not overriden at the app level

	if config.ClientId != "" && annotations[constants.AnnotationPrefix+"client-id"] == "" {
		configArgs = append(configArgs, "--client-id="+config.ClientId)
	}
	if config.ClientSecret != "" && annotations[constants.AnnotationPrefix+"client-secret"] == "" {
		configArgs = append(configArgs, "--client-secret="+config.ClientSecret)
	}
	if config.DiscoveryUrl != "" && annotations[constants.AnnotationPrefix+"discovery-url"] == "" {
		configArgs = append(configArgs, "--discovery-url="+config.DiscoveryUrl)
	}
	/*if config.EnableDefaultDeny !="" && annotations[constants.AnnotationPrefix+"enable-default-deny"] == "" {
		configArgs = append(configArgs, "--enable-default-deny="+config.EnableDefaultDeny)
	}*/
	if config.Listen != "" && annotations[constants.AnnotationPrefix+"listen"] == "" {
		configArgs = append(configArgs, "--listen="+config.Listen)
	}
	if config.SecureCookie != "" && annotations[constants.AnnotationPrefix+"secure-cookie"] == "" {
		configArgs = append(configArgs, "--secure-cookie="+config.SecureCookie)
	}
	if config.Verbose != "" && annotations[constants.AnnotationPrefix+"verbose"] == "" {
		configArgs = append(configArgs, "--verbose="+config.Verbose)
	}
	if config.EnableLogging != "" && annotations[constants.AnnotationPrefix+"enable-logging"] == "" {
		configArgs = append(configArgs, "--enable-logging="+config.EnableLogging)
	}
	for _, origin := range config.CorsOrigins {
		configArgs = append(configArgs, "--cors-origins="+origin)
	}
	for _, method := range config.CorsMethods {
		configArgs = append(configArgs, "--cors-methods="+method)
	}
	for _, resource := range config.Resources {
		//  --resources "uri=/admin*|roles=admin,superuser|methods=POST,DELETE
		res := ""
		if resource.URI != "" {
			res = "uri=" + resource.URI
		}
		if len(resource.Methods) != 0 {
			res = res + "|methods=" + strings.Join(resource.Methods, ",")
		}
		configArgs = append(configArgs, "--resources="+res)
	}
	for _, scope := range config.Scopes {
		configArgs = append(configArgs, "--scopes="+scope)
	}

	return configArgs
}
