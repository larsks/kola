package packagemanager

import (
	"path/filepath"
	"strconv"
	"strings"

	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
	"golang.org/x/exp/slices"
)

func MatchPackageGlobs(patterns ...string) PackageManifestFilter {
	for i := range patterns {
		patterns[i] = strings.ToLower(patterns[i])
	}

	return func(pkg *operators.PackageManifest) bool {
		for _, pattern := range patterns {
			if matches, _ := filepath.Match(pattern, strings.ToLower(pkg.Name)); matches {
				return true
			}
		}

		return false
	}
}

func MatchPackageSubstrings(patterns ...string) PackageManifestFilter {
	for i := range patterns {
		patterns[i] = strings.ToLower(patterns[i])
	}

	return func(pkg *operators.PackageManifest) bool {
		for _, pattern := range patterns {
			if strings.Contains(strings.ToLower(pkg.Name), pattern) {
				return true
			}
		}

		return false
	}
}

func MatchPackageName(pattern string) PackageManifestFilter {
	pattern = strings.ToLower(pattern)
	return func(pkg *operators.PackageManifest) bool {
		matches, _ := filepath.Match(pattern, strings.ToLower(pkg.Name))
		return matches
	}
}

func MatchCatalogSource(needle string) PackageManifestFilter {
	needle = strings.ToLower(needle)
	return func(pkg *operators.PackageManifest) bool {
		return strings.Contains(strings.ToLower(pkg.Status.CatalogSource), needle) ||
			strings.Contains(strings.ToLower(pkg.Status.CatalogSourceDisplayName), needle)
	}
}

func MatchDescription(needle string) PackageManifestFilter {
	needle = strings.ToLower(needle)
	return func(pkg *operators.PackageManifest) bool {
		for _, channel := range pkg.Status.Channels {
			if strings.Contains(strings.ToLower(channel.CurrentCSVDesc.LongDescription), needle) {
				return true
			}
		}

		return false
	}
}

func MatchInstallMode(installmode string) PackageManifestFilter {
	installmode = strings.ToLower(installmode)
	return func(pkg *operators.PackageManifest) bool {
		for _, channel := range pkg.Status.Channels {
			for _, mode := range channel.CurrentCSVDesc.InstallModes {
				if strings.ToLower(string(mode.Type)) == installmode && mode.Supported {
					return true
				}
			}
		}

		return false
	}
}

func MatchKeywords(keywords []string) PackageManifestFilter {
	return func(pkg *operators.PackageManifest) bool {
		for _, channel := range pkg.Status.Channels {
			for _, keyword := range keywords {
				if slices.Contains(channel.CurrentCSVDesc.Keywords, keyword) {
					return true
				}
			}
		}

		return false
	}
}

func MatchCertified(certified bool) PackageManifestFilter {
	certifiedString := strconv.FormatBool(certified)
	return func(pkg *operators.PackageManifest) bool {
		for _, channel := range pkg.Status.Channels {
			if channel.CurrentCSVDesc.Annotations["certified"] == certifiedString {
				return true
			}
		}

		return false
	}
}
