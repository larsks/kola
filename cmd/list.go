/*
Copyright © 2022 Lars Kellogg-Stedman <lars@oddbit.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"kola/packagemanager"
	"log"

	"github.com/spf13/cobra"
)

type (
	ListFlags struct {
		CatalogSource string   `short:"c" help:"Match string in package catalog source"`
		Description   string   `short:"d" help:"Match string in package description"`
		InstallMode   string   `short:"m" help:"Match package supported install mode"`
		Keyword       []string `short:"w" help:"Match package keyword"`
		Certified     bool     `short:"C" help:"Match only certified packages"`
		Glob          bool     `short:"g" help:"Arguments are glob patterns instead of substrings"`
	}
)

var listFlags = ListFlags{}

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available packages",
	Run:   runList,
}

func runList(cmd *cobra.Command, args []string) {
	pm, err := getCachedPackageManager(rootFlags.Kubeconfig)
	if err != nil {
		panic(err)
	}

	var filters []packagemanager.PackageManifestFilter

	if len(args) > 0 {
		if listFlags.Glob {
			filters = append(filters, packagemanager.MatchPackageGlobs(args...))
		} else {
			filters = append(filters, packagemanager.MatchPackageSubstrings(args...))
		}
	}

	if listFlags.CatalogSource != "" {
		filters = append(filters, packagemanager.MatchCatalogSource(listFlags.CatalogSource))
	}

	if listFlags.Description != "" {
		filters = append(filters, packagemanager.MatchDescription(listFlags.Description))
	}

	if listFlags.InstallMode != "" {
		filters = append(filters, packagemanager.MatchInstallMode(listFlags.InstallMode))
	}

	if len(listFlags.Keyword) > 0 {
		filters = append(filters, packagemanager.MatchKeywords(listFlags.Keyword))
	}

	if cmd.Flags().Lookup("certified").Changed {
		filters = append(filters, packagemanager.MatchCertified(listFlags.Certified))
	}

	packages, err := pm.ListPackageManifests(filters...)
	if err != nil {
		panic(err)
	}

	log.Printf("found %d packages", len(packages))

	for _, pkg := range packages {
		if rootFlags.Verbose > 1 {
			fmt.Printf("%s/%s %s\n", pkg.Status.CatalogSource, pkg.Name, pkg.Status.Channels[0].CurrentCSVDesc.DisplayName)
		} else if rootFlags.Verbose > 0 {
			fmt.Printf("%s/%s\n", pkg.Status.CatalogSource, pkg.Name)
		} else {
			fmt.Printf("%s\n", pkg.Name)
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	AddFlagsFromSpec(listCmd, &listFlags, false)
}
