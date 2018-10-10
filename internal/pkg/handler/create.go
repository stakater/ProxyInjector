package handler

import (
	logger "github.com/sirupsen/logrus"
	"github.com/stakater/ProxyInjector/internal/pkg/callbacks"
	"github.com/stakater/ProxyInjector/internal/pkg/constants"
	//"k8s.io/client-go/util/retry"
	"encoding/json"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/types"
)

// ResourceCreatedHandler contains new objects
type ResourceCreatedHandler struct {
	//client *kubernetes.Clientset,
	Resource interface{}
}

type patchValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value struct {
		Name  string `json:"name"`
		Image string `json:"image"`
		Ports struct {
			Name          string `json:"name"`
			Protocol      string `json:"protocol"`
			ContainerPort uint32 `json:"value"`
		}
	}
}

/*Name:  "web",
Image: "nginx:1.12",
Ports: []v1.ContainerPort{{
Name:          "http",
Protocol:      v1.ProtocolTCP,
ContainerPort: 80,
}},*/
// Handle processes the newly created resource
func (r ResourceCreatedHandler) Handle() error {
	if r.Resource == nil {
		logger.Errorf("Resource creation handler received nil resource")
	} else {
		logger.Info("Resource created")
		annotations := callbacks.GetDeploymentAnnotations(r.Resource)
		value := annotations[constants.AuthProxyUpstreamAnnotation]
		logger.Info(value)

		logger.Info("Updating deployment...")

		payload := []patchValue{{
			Op:   "merge",
			Path: "/spec/template/spec/containers",
			Value: {
				Name:  "web",
				Image: "nginx:1.12",
				Ports: {
					Name:          "http",
					Protocol:      v1.ProtocolTCP,
					ContainerPort: 80,
				},
			},
		}}

		payloadBytes, _ := json.Marshal(payload)
		_, err := r.Resource.(*v1beta1.Deployment).Patch(replicasetName, types.JSONPatchType, payloadBytes)
		return err

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

		logger.Info("Updated deployment...")
	}
	return nil
}
