package cmd

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/camelcase"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/exp/slices"
)

func longOptFromFieldName(s string) string {
	parts := camelcase.Split(s)
	return strings.ToLower(strings.Join(parts, "-"))
}

func AddFlagsFromSpec(command *cobra.Command, spec interface{}, persistent bool) {
	specType := reflect.TypeOf(spec)
	specElem := specType.Elem()
	specValue := reflect.ValueOf(spec)

	for i := 0; i < specElem.NumField(); i++ {
		field := specElem.Field(i)

		target := field.Tag.Get("target")
		if target == "" {
			target = field.Name
		}

		longOpt := field.Tag.Get("long")
		if longOpt == "" {
			longOpt = longOptFromFieldName(field.Name)
		}

		shortOpt := field.Tag.Get("short")
		helpText := field.Tag.Get("help")
		envvar := field.Tag.Get("envvar")
		defval := field.Tag.Get("default")
		hide := field.Tag.Get("hide")
		subtype := field.Tag.Get("subtype")

		if defval == "" && envvar != "" {
			defval = os.Getenv(envvar)
		}

		var flagset *pflag.FlagSet
		if persistent {
			flagset = command.PersistentFlags()
		} else {
			flagset = command.Flags()
		}

		switch p := specValue.Elem().Field(i).Interface().(type) {
		case string:
			ptr := specValue.Elem().FieldByName(target).Addr().Interface().(*string)
			flagset.StringVarP(ptr, longOpt, shortOpt, defval, helpText)
		case []string:
			ptr := specValue.Elem().FieldByName(target).Addr().Interface().(*[]string)
			flagset.StringSliceVarP(ptr, longOpt, shortOpt, []string{}, helpText)
		case int:
			ptr := specValue.Elem().FieldByName(target).Addr().Interface().(*int)
			switch subtype {
			case "counter":
				flagset.CountVarP(ptr, longOpt, shortOpt, helpText)
			default:
				flagset.IntVarP(ptr, longOpt, shortOpt, stringToInt(defval), helpText)
			}
		case bool:
			ptr := specValue.Elem().FieldByName(target).Addr().Interface().(*bool)
			flagset.BoolVarP(ptr, longOpt, shortOpt, stringToBool(defval), helpText)
		case time.Duration:
			ptr := specValue.Elem().FieldByName(target).Addr().Interface().(*time.Duration)
			flagset.DurationVarP(ptr, longOpt, shortOpt, stringToDuration(defval), helpText)
		default:
			fmt.Printf("unsupported: %v\n", p)
		}

		if stringToBool(hide) {
			//nolint:errcheck
			flagset.MarkHidden(longOpt)
		}
	}
}

func stringToBool(s string) bool {
	return slices.Contains([]string{"1", "true"}, strings.ToLower(s))
}

func stringToInt(s string) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return 0
}

func stringToDuration(s string) time.Duration {
	if val, err := time.ParseDuration(s); err == nil {
		return val
	}
	return 0
}
