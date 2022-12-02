package packagemanager

import (
	"context"
	"encoding/json"
	"fmt"
	"kola/cache"
	"log"

	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
	"k8s.io/client-go/kubernetes"
)

type (
	PackageManager struct {
		clientset *kubernetes.Clientset
		cache     cache.Cache
	}

	PackageManifestFilter func(pkg *operators.PackageManifest) bool
)

func NewPackageManager(clientset *kubernetes.Clientset) *PackageManager {
	return &PackageManager{
		clientset: clientset,
	}
}

func (pm *PackageManager) WithCache(cache cache.Cache) *PackageManager {
	pm.cache = cache
	return pm
}

func (pm *PackageManager) getCached(path string) ([]byte, error) {
	var data []byte
	var err error

	if pm.cache != nil {
		if data, err = pm.cache.Get(path); err != nil {
			log.Printf("cache fetch failed: %v", err)
			data = nil
		}
	}

	if data == nil {
		if data, err = pm.clientset.RESTClient().Get().AbsPath(path).DoRaw(context.TODO()); err != nil {
			return nil, err
		}

		if pm.cache != nil {
			if err = pm.cache.Put(path, data); err != nil {
				log.Printf("cache store failed: %v", err)
			}
		}
	}

	return data, nil
}

func (pm *PackageManager) GetPackageManifest(packageName string) (*operators.PackageManifest, error) {
	var pkg operators.PackageManifest

	data, err := pm.getCached(
		fmt.Sprintf("/apis/packages.operators.coreos.com/v1/namespaces/default/packagemanifests/%s", packageName))
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

	data, err := pm.getCached("/apis/packages.operators.coreos.com/v1/namespaces/default/packagemanifests")
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
