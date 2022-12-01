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
	"kola/client"
	"kola/packagemanager"
	"log"

	"github.com/spf13/cobra"
)

type (
	ListFlags struct {
		CatalogSource string   `short:"c" help:"packagemanager.Match string in package catalog source"`
		Description   string   `short:"d" help:"packagemanager.Match string in package description"`
		InstallMode   string   `short:"m" help:"packagemanager.Match package supported install mode"`
		Keyword       []string `short:"w" help:"packagemanager.Match package keyword"`
		Certified     bool     `short:"C" help:"packagemanager.Match only certified packages"`
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
	clientset, err := client.GetClient(rootFlags.Kubeconfig)
	if err != nil {
		panic(err)
	}

	pm := packagemanager.NewPackageManager(clientset)

	var filters []packagemanager.PackageManifestFilter

	if len(args) > 0 {
		filters = append(filters, packagemanager.MatchPackageNames(args...))
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

	packages, err := pm.ListPackageManifests(filters...)
	if err != nil {
		panic(err)
	}

	log.Printf("found %d packages", len(packages))

	for _, pkg := range packages {
		if rootFlags.Verbose > 0 {
			fmt.Printf("%s (%s)\n", pkg.Name, pkg.Status.Channels[0].CurrentCSVDesc.DisplayName)
		} else {
			fmt.Printf("%s\n", pkg.Name)
		}
	}
}

func init() {
	rootCmd.AddCommand(listCmd)
	AddFlagsFromSpec(listCmd, &listFlags, false)
}
