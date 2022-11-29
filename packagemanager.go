package main

import (
	"context"
	"encoding/json"
	"fmt"

	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
	"k8s.io/client-go/kubernetes"
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
