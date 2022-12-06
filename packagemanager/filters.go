package packagemanager

import (
	"path/filepath"
	"strconv"
	"strings"

	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
	"golang.org/x/exp/slices"
)

// Return a file that matches package names against a list of
// glob patterns. Comparisons are case insensitive.
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

// Return a filter that matches package names against a list of words.
// Comparisons are case insensitive.
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

// Returns a filter that matches a specific package name. Comparisons are
// case insensitive.
func MatchPackageName(pattern string) PackageManifestFilter {
	pattern = strings.ToLower(pattern)
	return func(pkg *operators.PackageManifest) bool {
		matches, _ := filepath.Match(pattern, strings.ToLower(pkg.Name))
		return matches
	}
}

// Return a filter that matches the package CatalogSource against a substring.
// Comparisons are case insensitive.
func MatchCatalogSource(needle string) PackageManifestFilter {
	needle = strings.ToLower(needle)
	return func(pkg *operators.PackageManifest) bool {
		return strings.Contains(strings.ToLower(pkg.Status.CatalogSource), needle) ||
			strings.Contains(strings.ToLower(pkg.Status.CatalogSourceDisplayName), needle)
	}
}

// Return a filter that matches the package description against a substring.
// Comparisons are case insensitive.
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

// Return a filter that matches the package InstallMode against a substring.
// Comparisons are case insensitive.
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

// Return a filter that matches packages if they contain any of the given
// keywords.
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

// Return a filter that matches packages based on the value of the
// certified attribute.
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
