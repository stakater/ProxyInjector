package cmd

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stakater/ProxyInjector/internal/pkg/config"
	"github.com/stakater/ProxyInjector/internal/pkg/controller"
	"github.com/stakater/ProxyInjector/pkg/kube"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewProxyInjectorCommand starts the proxy injector controller
func NewProxyInjectorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proxyInjector",
		Short: "An authentication proxy injector for Kubernetes pods",
		Run:   startProxyInjector,
	}
	return cmd
}

func startProxyInjector(cmd *cobra.Command, args []string) {
	logrus.Info("Starting ProxyInjector")
	currentNamespace := os.Getenv("KUBERNETES_NAMESPACE")
	if len(currentNamespace) == 0 {
		currentNamespace = v1.NamespaceAll
		logrus.Warnf("KUBERNETES_NAMESPACE is unset, will detect changes in all namespaces.")
	}

	// create the clientset
	clientset, err := kube.GetClient()
	if err != nil {
		logrus.Fatal(err)
	}

	config := config.GetControllerConfig()

	for resource := range kube.ResourceMap {
		c, err := controller.NewController(clientset, resource, config, currentNamespace)
		if err != nil {
			logrus.Fatalf("%s", err)
		}

		// Now let's start the controller
		stop := make(chan struct{})
		defer close(stop)

		go c.Run(1, stop)
	}

	// Wait forever
	select {}
}
