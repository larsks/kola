package main

// https://pkg.go.dev/github.com/RedHatInsights/clowder/controllers/cloud.redhat.com/providers/metrics/subscriptions#SubscriptionSpec
// https://pkg.go.dev/github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators#PackageManifestList

import (
	"errors"
	"fmt"
	"log"
	"os"
	"text/template"

	operatorsv1alpha1 "github.com/operator-framework/api/pkg/operators/v1alpha1"
	operators "github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators/v1"
	"golang.org/x/exp/slices"
	"k8s.io/client-go/kubernetes"
)

type (
	Options struct {
		Kubeconfig       string `short:"k" help:"Path to kubernetes client configuration"`
		CatalogSource    string `short:"c" help:"Match string in package catalog source"`
		Description      string `short:"d" help:"Match string in package description"`
		InstallMode      string `short:"m" help:"Match package supported install mode"`
		Keyword          string `short:"w" help:"Match package keyword"`
		PackageName      string `short:"n" long:"name" help:"Match package name"`
		Certified        bool   `short:"C" help:"Match only certified packages"`
		Show             bool   `short:"s" help:"Show details about matched packages"`
		Subscribe        bool   `short:"S" help:"Generate subscriptions for matched packages"`
		ShowDescription  bool   `short:"D" help:"Show package descriptions when using --show"`
		InstallNamespace string `help:"Namespace for subscription"`
		InstallChannel   string `help:"Select installation channel"`
		InstallApproval  string `help:"Select manual or automatic approval for updates"`
	}
)

var (
	validInstallModes = [...]operatorsv1alpha1.InstallModeType{
		operatorsv1alpha1.InstallModeTypeOwnNamespace,
		operatorsv1alpha1.InstallModeTypeSingleNamespace,
		operatorsv1alpha1.InstallModeTypeMultiNamespace,
		operatorsv1alpha1.InstallModeTypeAllNamespaces,
	}
)

func (options *Options) ValidateInstallApproval(key string) error {
	if options.InstallApproval == "" {
		return nil
	}

	if !slices.Contains([]string{"Manual", "Automatic"}, options.InstallApproval) {
		return ApplicationErrorType{
			Message: fmt.Sprintf("%s is not a valid approval method", options.InstallApproval),
			Parent:  ValidationError,
		}
	}

	return nil
}

func (options *Options) ValidateInstallMode(key string) error {
	if options.InstallMode == "" {
		return nil
	}

	for _, mode := range validInstallModes {
		if string(mode) == options.InstallMode {
			return nil
		}
	}

	return ApplicationErrorType{
		Message: fmt.Sprintf("%s is not a valid install mode", options.InstallMode),
		Parent:  ValidationError,
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			if errors.Is(r.(error), ApplicationError) {
				fmt.Printf("ERROR: %v\n", r)
			} else {
				panic(r)
			}
		}
	}()

	options := Options{}
	flagset := BuildFlagsFromStruct(os.Args[0], &options)
	if err := flagset.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	if err := ValidateOptions(&options); err != nil {
		panic(err)
	}

	config, err := BuildConfigFromFlags("", options.Kubeconfig)
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pm := PackageManager{
		clientset: clientset,
	}

	//	if options.PackageName != "" {
	//		pkg, err := pm.GetPackageManifest(options.PackageName)
	//		if err != nil {
	//			panic(err)
	//		}
	//
	//		fmt.Printf("%+v\n", pkg)
	//		return
	//	}

	var filters []PackageManifestFilter

	if options.PackageName != "" {
		filters = append(filters, MatchPackageName(options.PackageName))
	}

	if options.CatalogSource != "" {
		filters = append(filters, MatchCatalogSource(options.CatalogSource))
	}

	if options.Description != "" {
		filters = append(filters, MatchDescription(options.Description))
	}

	if options.InstallMode != "" {
		filters = append(filters, MatchInstallMode(options.InstallMode))
	}

	if options.Keyword != "" {
		filters = append(filters, MatchKeyword(options.Keyword))
	}

	res, err := pm.ListPackageManifests(filters...)
	if err != nil {
		panic(err)
	}

	log.Printf("found %d packages", len(res))
	for _, pkg := range res {
		var err error

		switch {
		case options.Show:
			err = showPackage(&pkg, &options)
		case options.Subscribe:
			fmt.Printf("subscribe\n")
		default:
			fmt.Printf("%s\n", pkg.Name)
		}

		if err != nil {
			panic(err)
		}
	}
}

func showPackage(pkg *operators.PackageManifest, options *Options) error {
	data := struct {
		Package *operators.PackageManifest
		Options *Options
	}{pkg, options}

	tmpl, err := template.New("package").Parse(`
Name: {{ .Package.Name }}
Catalog source: {{ .Package.Status.CatalogSourceDisplayName }} ({{ .Package.Status.CatalogSource }})
Publisher: {{ .Package.Status.CatalogSourcePublisher }}
Provider: {{ .Package.Status.Provider.Name }}
Channels:
{{ range .Package.Status.Channels -}}
  - {{ .Name }} ({{ .CurrentCSV }})
{{ end }}
{{ if .Options.ShowDescription -}}
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
