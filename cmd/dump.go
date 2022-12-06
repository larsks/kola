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
	"fmt"
	"os"

	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/kubectl/pkg/scheme"
)

type (
	DumpFlags struct {
	}
)

var dumpFlags = DumpFlags{}

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump details about a package",
	RunE:  runDump,
}

func init() {
	rootCmd.AddCommand(dumpCmd)
	AddFlagsFromSpec(dumpCmd, &dumpFlags, false)
}

func runDump(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("dump: %w", err)
		}
	}()

	if len(args) != 1 {
		return errors.New("dump requires a single package name")
	}

	pm, err := getCachedPackageManager(rootFlags.Kubeconfig)
	if err != nil {
		return err
	}
	pkg, err := pm.GetPackageManifest(args[0])
	if err != nil {
		return err
	}

	if err := dumpPackage(pkg); err != nil {
		return err
	}

	return nil
}

func dumpPackage(pkg *operators.PackageManifest) error {
	//nolint:errcheck
	{
		operatorsv1alpha1.AddToScheme(scheme.Scheme)
		operators.AddToScheme(scheme.Scheme)
	}
	s := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme.Scheme,
		scheme.Scheme)

	err := s.Encode(pkg, os.Stdout)
	return err
}
