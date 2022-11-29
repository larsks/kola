package main

// https://pkg.go.dev/github.com/RedHatInsights/clowder/controllers/cloud.redhat.com/providers/metrics/subscriptions#SubscriptionSpec
// https://pkg.go.dev/github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators#PackageManifestList

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"

	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"

	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type (
	PackageManager struct {
		clientset *kubernetes.Clientset
	}

	PackageManifestFilter func(pkg *operators.PackageManifest) bool
)

func (pm *PackageManager) GetPackageManifest(packageName string) operators.PackageManifest {
	var pkg operators.PackageManifest

	data, err := pm.clientset.RESTClient().Get().AbsPath(
		fmt.Sprintf("/apis/packages.operators.coreos.com/v1/namespaces/default/packagemanifests/%s", packageName)).DoRaw(context.TODO())
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		panic(err)
	}
	return pkg
}

func (pm *PackageManager) ListPackageManifests(filters ...PackageManifestFilter) []operators.PackageManifest {

	pkgs := &operators.PackageManifestList{}
	selected := []operators.PackageManifest{}

	data, err := pm.clientset.RESTClient().Get().AbsPath(
		"/apis/packages.operators.coreos.com/v1/namespaces/default/packagemanifests").DoRaw(context.TODO())
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, pkgs); err != nil {
		panic(err)
	}

	if len(filters) > 0 {
	PACKAGES:
		for _, pkg := range pkgs.Items {
			for _, filter := range filters {
				if !filter(&pkg) {
					continue PACKAGES
				}
			}

			selected = append(selected, pkg)
		}

		return selected
	}

	return pkgs.Items
}

func BuildConfigFromFlags(masterUrl, kubeconfigPath string) (*restclient.Config, error) {
	if kubeconfigPath == "" && masterUrl == "" {
		if kubeconfig, err := restclient.InClusterConfig(); err == nil {
			return kubeconfig, nil
		}

		return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			nil,
		).ClientConfig()
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: masterUrl}}).ClientConfig()
}

func main() {
	var kubeconfig *string
	var matchCatalogSource *string
	var matchDescription *string
	var matchInstallMode *string
	var matchKeyword *string
	var matchName *string
	var packageName *string

	matchCatalogSource = flag.String("catalogSource", "", "match substring in catalog source")
	matchDescription = flag.String("description", "", "match substring in description")
	matchName = flag.String("name", "", "match package names against glob pattern")
	matchKeyword = flag.String("keyword", "", "match keywords")
	matchInstallMode = flag.String("installmode", "", "match installmode")
	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	packageName = flag.String("packageName", "", "get single package")
	flag.Parse()

	config, err := BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pm := PackageManager{
		clientset: clientset,
	}

	if *packageName != "" {
		pkg := pm.GetPackageManifest(*packageName)
		fmt.Printf("%+v\n", pkg)
		return
	}

	var filters []PackageManifestFilter

	if *matchName != "" {
		filters = append(filters, MatchPackageName(*matchName))
	}

	if *matchCatalogSource != "" {
		filters = append(filters, MatchCatalogSource(*matchCatalogSource))
	}

	if *matchDescription != "" {
		filters = append(filters, MatchDescription(*matchDescription))
	}

	if *matchInstallMode != "" {
		filters = append(filters, MatchInstallMode(*matchInstallMode))
	}

	if *matchKeyword != "" {
		filters = append(filters, MatchKeyword(*matchKeyword))
	}

	res := pm.ListPackageManifests(filters...)
	fmt.Printf("found %d packages\n", len(res))
	for _, pkg := range res {
		fmt.Printf("%s\n", pkg.Name)
	}
}
