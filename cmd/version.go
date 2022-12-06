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
	"kola/version"

	"github.com/spf13/cobra"
)

type (
	VersionFlags struct {
	}
)

var versionFlags = VersionFlags{}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:          "version",
	Short:        "Show command version",
	RunE:         runVersion,
	SilenceUsage: true,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	AddFlagsFromSpec(versionCmd, &versionFlags, false)
}

func runVersion(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("version: %w", err)
		}
	}()

	fmt.Printf("Version %s built on %s from %s\n",
		version.BuildVersion, version.BuildDate, version.BuildRef)

	return nil
}
