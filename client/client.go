package client

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Attempt to find a valid kubernetes client configuration:
//
// 1. If the caller has provided an explicit kubeconfig path, use that with
//    no further processing.
// 2. Otherwise, first check for the KUBECONFIG environment variable and
//    the default config file location.
// 3. Lastly, check for an in-cluster configuration.
func BuildConfigFromFlags(masterUrl, kubeconfigPath string) (*restclient.Config, error) {
	if kubeconfigPath == "" && masterUrl == "" {
		if kubeconfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			nil,
		).ClientConfig(); err == nil {
			return kubeconfig, nil
		}

		return restclient.InClusterConfig()
	}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: masterUrl}}).ClientConfig()
}

func GetClient(kubeconfig string) (*rest.Config, *kubernetes.Clientset, error) {
	config, err := BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build Kuberbetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return config, clientset, nil
}
