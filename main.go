package main

// https://pkg.go.dev/github.com/RedHatInsights/clowder/controllers/cloud.redhat.com/providers/metrics/subscriptions#SubscriptionSpec
// https://pkg.go.dev/github.com/operator-framework/operator-lifecycle-manager/pkg/package-server/apis/operators#PackageManifestList

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	flag "github.com/spf13/pflag"
)

type (
	Options struct {
		Kubeconfig    string `short:"k" help:"Path to kubernetes client configuration"`
		CatalogSource string `short:"c" help:"Match string in package catalog source"`
		Description   string `short:"d" help:"Match string in package description"`
		InstallMode   string `short:"m" help:"Match package supported install mode"`
		Keyword       string `short:"w" help:"Match package keyword"`
		Name          string `short:"n" help:"Match package name"`
		Certified     bool   `short:"C" help:"Match only certified packages"`
	}
)

func (options *Options) ValidateInstallMode() error {
	fmt.Printf("validateInstallMode\n")
	return nil
}

func buildFlagsFromStruct(name string, options interface{}) *flag.FlagSet {
	flagset := flag.NewFlagSet(name, flag.ExitOnError)

	t := reflect.TypeOf(options)
	e := t.Elem()
	v := reflect.ValueOf(options)

	for i := 0; i < e.NumField(); i++ {
		field := e.Field(i)
		fmt.Printf("field %d: %s\n", i, field.Name)

		longOpt := field.Tag.Get("long")
		if longOpt == "" {
			longOpt = strings.ToLower(string(field.Name[0])) + string(field.Name[1:])
		}

		shortOpt := field.Tag.Get("short")
		helpText := field.Tag.Get("short")
		validator := v.MethodByName(fmt.Sprintf("Validate%s", field.Name))
		if validator.IsValid() {
			println("found validator")
		}

		fmt.Printf("%s long --%s short -%s help %s\n", field.Name, longOpt, shortOpt, helpText)
		switch p := v.Elem().Field(i).Interface().(type) {
		case string:
			ptr := v.Elem().Field(i).Addr().Interface().(*string)
			flagset.StringVarP(ptr, longOpt, shortOpt, "", helpText)
		case bool:
			ptr := v.Elem().Field(i).Addr().Interface().(*bool)
			flagset.BoolVarP(ptr, longOpt, shortOpt, false, helpText)
		default:
			fmt.Printf("wtf: %v\n", p)
		}
	}

	return flagset
}

func main() {
	options := Options{}
	flagset := buildFlagsFromStruct(os.Args[0], &options)
	if err := flagset.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	fmt.Printf("options: %+v\n", options)
}
