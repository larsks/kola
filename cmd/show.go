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
	"html/template"
	"kola/packagemanager"
	"os"

	"github.com/spf13/cobra"

	_ "embed"
)

type (
	ShowFlags struct {
	}
)

var (
	showFlags = ShowFlags{}

	// showCmd represents the show command
	showCmd = &cobra.Command{
		Use:          "show",
		Short:        "Show details about a package",
		RunE:         runShow,
		SilenceUsage: true,
	}

	//go:embed templates/show.tpl
	showTemplate string
)

func init() {
	rootCmd.AddCommand(showCmd)
	AddFlagsFromSpec(showCmd, &showFlags, false)
}

func runShow(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("show: %w", err)
		}
	}()

	pm, err := getCachedPackageManager(rootFlags.Kubeconfig)
	if err != nil {
		return err
	}

	for _, pkgName := range args {
		pkg, err := pm.GetPackageManifest(pkgName)
		if err != nil {
			return err
		}
		if err := showPackage(pkg); err != nil {
			return err
		}
	}

	return nil
}

func showPackage(pkg *packagemanager.Package) error {
	data := struct {
		Package *packagemanager.Package
		Flags   *ShowFlags
		Verbose int
	}{pkg, &showFlags, rootFlags.Verbose}

	tmpl, err := template.New("package").Parse(showTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(os.Stdout, data); err != nil {
		return err
	}

	return nil
}
