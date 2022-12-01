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

	"github.com/spf13/cobra"
)

type (
	SubscribeFlags struct {
		Channel   string `short:"c" help:"Set channel for subscription"`
		Approval  string `short:"a" help:"Set install plan approval for subscription" default:"Automatic"`
		Namespace string `short:"n" help:"Set namespace for subscription"`
	}
)

var subscribeFlags = SubscribeFlags{}

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Aliases: []string{"sub"},
	Use:     "subscribe",
	Short:   "Generate a Subscription for a package",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscribe called")
	},
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
	AddFlagsFromSpec(subscribeCmd, &subscribeFlags, false)
}
