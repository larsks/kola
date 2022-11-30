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

func (pm *PackageManager) GetPackageManifest(packageName string) (*operators.PackageManifest, error) {
	var pkg operators.PackageManifest

	data, err := pm.clientset.RESTClient().Get().AbsPath(
		fmt.Sprintf("/apis/packages.operators.coreos.com/v1/namespaces/default/packagemanifests/%s", packageName)).DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (pm *PackageManager) ListPackageManifests(filters ...PackageManifestFilter) ([]operators.PackageManifest, error) {

	pkgs := &operators.PackageManifestList{}
	selected := []operators.PackageManifest{}

	data, err := pm.clientset.RESTClient().Get().AbsPath(
		"/apis/packages.operators.coreos.com/v1/namespaces/default/packagemanifests").DoRaw(context.TODO())
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, pkgs); err != nil {
		return nil, err
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
	} else {
		selected = pkgs.Items
	}

	return selected, nil
}
