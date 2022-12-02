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
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type (
	RootFlags struct {
		Kubeconfig    string        `short:"k" help:"Path to kubernetes client configuration"`
		Verbose       int           `subtype:"counter" short:"v" help:"Increase output verbosity"`
		Debug         bool          `help:"Traceback on panic" hide:"true"`
		CacheLifetime time.Duration `default:"10m" help:"Set cache lifetime"`
		NoCache       bool          `help:"Disable local caching of results"`
	}
)

var rootCmd = &cobra.Command{
	Use:   "kola",
	Short: "Interact with OLM package manifests",
}

var rootFlags = RootFlags{}

func Execute() {
	defer func() {
		if !rootFlags.Debug {
			if r := recover(); r != nil {
				err := r.(error)
				log.Fatalf("ERROR: %v", err)
			}
		}
	}()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	AddFlagsFromSpec(rootCmd, &rootFlags, true)
}
