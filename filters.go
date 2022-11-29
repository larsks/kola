package main

import (
	"fmt"
	"path/filepath"
	"strings"

	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
	"golang.org/x/exp/slices"
)

func MatchPackageName(pattern string) PackageManifestFilter {
	pattern = fmt.Sprintf("*%s*", strings.ToLower(pattern))
	return func(pkg *operators.PackageManifest) bool {
		matches, err := filepath.Match(pattern, strings.ToLower(pkg.Name))
		if err != nil {
			panic(err)
		}

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
func MatchKeyword(keyword string) PackageManifestFilter {
	keyword = strings.ToLower(keyword)
	return func(pkg *operators.PackageManifest) bool {
		for _, channel := range pkg.Status.Channels {
			if slices.Contains(channel.CurrentCSVDesc.Keywords, keyword) {
				return true
			}
		}

		return false
	}
}
