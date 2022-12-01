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
	"errors"
	"html/template"
	"kola/client"
	"kola/packagemanager"
	"os"

	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
	"github.com/spf13/cobra"
)

type (
	ShowFlags struct {
		Description bool `short:"d" help:"Include description in output"`
	}
)

var showFlags = ShowFlags{}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show details about a package",
	Run:   runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
	AddFlagsFromSpec(showCmd, &showFlags, false)
}

func runShow(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		panic(errors.New("show requires a single package name"))
	}

	clientset, err := client.GetClient(rootFlags.Kubeconfig)
	if err != nil {
		panic(err)
	}

	pm := packagemanager.NewPackageManager(clientset)
	pkg, err := pm.GetPackageManifest(args[0])
	if err != nil {
		panic(err)
	}

	if err := showPackage(pkg); err != nil {
		panic(err)
	}
}

func showPackage(pkg *operators.PackageManifest) error {
	data := struct {
		Package *operators.PackageManifest
		Flags   *ShowFlags
	}{pkg, &showFlags}

	tmpl, err := template.New("package").Parse(`
Name: {{ .Package.Name }}
Catalog source: {{ .Package.Status.CatalogSourceDisplayName }} ({{ .Package.Status.CatalogSource }})
Publisher: {{ .Package.Status.CatalogSourcePublisher }}
Provider: {{ .Package.Status.Provider.Name }}
Channels:
{{ range .Package.Status.Channels -}}
  - {{ .Name }} ({{ .CurrentCSV }})
{{ end }}
{{ if .Flags.Description -}}
Description:
{{ (index .Package.Status.Channels 0).CurrentCSVDesc.LongDescription }}
{{ end -}}
`)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		return err
	}

	return nil
}
