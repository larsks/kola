package cmd

import (
	"crypto/sha256"
	"fmt"
	"kola/cache"
	"kola/client"
	"kola/packagemanager"
	"log"
)

// Return a new PackageManager with an associated Cache (unless --no-cache
// was specified at runtime).
func getCachedPackageManager(kubeconfig string) (*packagemanager.PackageManager, error) {
	config, clientset, err := client.GetClient(kubeconfig)
	if err != nil {
		return nil, err
	}

	pm := packagemanager.NewPackageManager(clientset)

	if !rootFlags.NoCache {

		// Generate a hash of (Host, APIPath) to use as a cache
		// identifier. This ensures we don't accidentally use
		// cached information for the wrong remote host.
		hash := sha256.New()
		hash.Write([]byte(config.Host))
		hash.Write([]byte(config.APIPath))

		cache := cache.NewCache("kola", fmt.Sprintf("%x", hash.Sum(nil))).
			WithLifetime(rootFlags.CacheLifetime)
		if err := cache.Start(); err != nil {
			log.Printf("failed to start cache: %v", err)
		} else {
			pm = pm.WithCache(cache)
		}
	}
	return pm, nil
}
