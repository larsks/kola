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
	"os"
	"strings"

	operatorsv1 "github.com/operator-framework/api/pkg/operators/v1"
	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slices"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/kubectl/pkg/scheme"
)

type (
	SubscribeFlags struct {
		Channel             string   `short:"c" help:"Set channel for subscription"`
		Approval            string   `short:"a" help:"Set install plan approval for subscription" default:"Automatic"`
		Namespace           string   `short:"n" help:"Set namespace for subscription"`
		CreateNamespace     bool     `short:"N" help:"Create a namespace"`
		CreateOperatorGroup bool     `short:"G" help:"Create an OperatorGroup"`
		TargetNamespace     []string `short:"t" help:"Set a target namespace"`
		Selector            []string `short:"l" help:"Set a namespace selector"`
	}
)

var subscribeFlags = SubscribeFlags{}

var validApprovals = []string{
	"",
	"Automatic",
	"Manual",
}

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Aliases:      []string{"sub"},
	Use:          "subscribe",
	Short:        "Generate a Subscription for a package",
	RunE:         runSubscribe,
	Args:         cobra.ExactArgs(1),
	SilenceUsage: true,
}

func (flags *SubscribeFlags) Validate() error {
	if !slices.Contains(validApprovals, flags.Approval) {
		return NewValidationError(
			"Invalid approval",
			flags.Approval,
		)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(subscribeCmd)
	AddFlagsFromSpec(subscribeCmd, &subscribeFlags, false)
}

func runSubscribe(cmd *cobra.Command, args []string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("subscribe: %w", err)
		}
	}()

	pm, err := getCachedPackageManager(rootFlags.Kubeconfig)
	if err != nil {
		return err
	}

	pkg, err := pm.GetPackageManifest(args[0])
	if err != nil {
		return err
	}

	if err := subscribePackage(pkg); err != nil {
		return err
	}

	return nil
}

func subscribePackage(pkg *operators.PackageManifest) error {
	channelName := subscribeFlags.Channel
	if channelName == "" {
		channelName = pkg.Status.DefaultChannel
	}

	var channel *operators.PackageChannel
	for _, check := range pkg.Status.Channels {
		if check.Name == channelName {
			channel = &check
			break
		}
	}

	if channel == nil {
		return fmt.Errorf("no such channel named %s for package %s",
			channelName, pkg.Name)
	}

	namespaceName := subscribeFlags.Namespace
	if namespaceName == "" {
		if suggested, ok := channel.CurrentCSVDesc.Annotations["operatorframework.io/suggested-namespace"]; ok {
			namespaceName = suggested
		}
	}

	subscription := operatorsv1alpha1.Subscription{
		TypeMeta: metav1.TypeMeta{
			APIVersion: operatorsv1alpha1.SubscriptionCRDAPIVersion,
			Kind:       operatorsv1alpha1.SubscriptionKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespaceName,
			Name:      pkg.Name,
		},
		Spec: &operatorsv1alpha1.SubscriptionSpec{
			Package:                pkg.Name,
			Channel:                channel.Name,
			InstallPlanApproval:    operatorsv1alpha1.Approval(subscribeFlags.Approval),
			CatalogSource:          pkg.Status.CatalogSource,
			CatalogSourceNamespace: pkg.Status.CatalogSourceNamespace,
		},
	}

	//nolint:errcheck
	operatorsv1alpha1.AddToScheme(scheme.Scheme)
	//nolint:errcheck
	corev1.AddToScheme(scheme.Scheme)

	serializer := json.NewSerializerWithOptions(
		json.DefaultMetaFactory, scheme.Scheme, scheme.Scheme,
		json.SerializerOptions{
			Pretty: true,
			Yaml:   true,
			Strict: true,
		})

	if err := serializer.Encode(&subscription, os.Stdout); err != nil {
		return err
	}

	if subscribeFlags.CreateNamespace {
		namespace := corev1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Namespace",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: namespaceName,
			},
		}
		os.Stdout.Write([]byte("---\n"))
		if err := serializer.Encode(&namespace, os.Stdout); err != nil {
			return err
		}
	}

	if subscribeFlags.CreateOperatorGroup {
		operatorgroup := operatorsv1.OperatorGroup{
			TypeMeta: metav1.TypeMeta{
				Kind:       "OperatorGroup",
				APIVersion: "operators.coreos.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespaceName,
				Name:      pkg.Name,
			},
		}

		if len(subscribeFlags.TargetNamespace) > 0 {
			operatorgroup.Spec.TargetNamespaces = subscribeFlags.TargetNamespace
		} else if len(subscribeFlags.Selector) > 0 {
			operatorgroup.Spec.Selector = &metav1.LabelSelector{}
			operatorgroup.Spec.Selector.MatchLabels = make(map[string]string)
			for _, selector := range subscribeFlags.Selector {
				kv := strings.Split(selector, "=")
				if len(kv) == 2 {
					operatorgroup.Spec.Selector.MatchLabels[kv[0]] = kv[1]
				}
			}
		}

		os.Stdout.Write([]byte("---\n"))
		if err := serializer.Encode(&operatorgroup,
			os.Stdout); err != nil {
			return err
		}
	}

	return nil
}
