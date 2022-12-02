package cmd

import (
	"crypto/sha256"
	"fmt"
	"kola/cache"
	"kola/client"
	"kola/packagemanager"
)

func getCachedPackageManager(kubeconfig string) (*packagemanager.PackageManager, error) {
	config, clientset, err := client.GetClient(kubeconfig)
	if err != nil {
		return nil, err
	}

	pm := packagemanager.NewPackageManager(clientset)

	if !rootFlags.NoCache {
		hash := sha256.New()
		hash.Write([]byte(config.Host))
		hash.Write([]byte(config.APIPath))

		cache := cache.NewCache("kola", fmt.Sprintf("%x", hash.Sum(nil))).
			WithLifetime(rootFlags.CacheLifetime)
		if err := cache.Start(); err != nil {
			return nil, err
		}
		pm = pm.WithCache(cache)
	}
	return pm, nil
}
