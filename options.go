package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	flag "github.com/spf13/pflag"
	"golang.org/x/exp/slices"
)

type (
	ApplicationErrorType struct {
		Message string
		Parent  error
	}
)

var (
	ApplicationError ApplicationErrorType = ApplicationErrorType{
		Message: "An unexpected error has occurred",
	}
)

func NewApplicationError(msg string, wraps error) error {
	return ApplicationErrorType{
		Message: msg,
		Parent:  wraps,
	}
}

func (err ApplicationErrorType) Error() string {
	return err.Message
}

func (err ApplicationErrorType) Unwrap() error {
	return err.Parent
}

func ValidateOptions(options interface{}) error {
	t := reflect.TypeOf(options)
	e := t.Elem()
	v := reflect.ValueOf(options)

	for i := 0; i < e.NumField(); i++ {
		field := e.Field(i)

		target := field.Tag.Get("target")
		if target != "" {
			continue
		}

		validator := v.MethodByName(fmt.Sprintf("Validate%s", field.Name))

		if validator.IsValid() {
			ret := validator.Call([]reflect.Value{reflect.ValueOf(field.Name)})
			err := ret[0].Interface()
			if err != nil {
				return err.(error)
			}
		}
	}

	return nil
}

func BuildFlagsFromStruct(name string, options interface{}) *flag.FlagSet {
	flagset := flag.NewFlagSet(name, flag.ExitOnError)

	t := reflect.TypeOf(options)
	e := t.Elem()
	v := reflect.ValueOf(options)

	for i := 0; i < e.NumField(); i++ {
		field := e.Field(i)

		target := field.Tag.Get("target")
		if target == "" {
			target = field.Name
		}

		longOpt := field.Tag.Get("long")
		if longOpt == "" {
			longOpt = strings.ToLower(string(field.Name[0])) + string(field.Name[1:])
		}

		shortOpt := field.Tag.Get("short")
		helpText := field.Tag.Get("help")
		envvar := field.Tag.Get("envvar")
		defval := field.Tag.Get("default")
		hide := field.Tag.Get("hide")

		if defval == "" && envvar != "" {
			defval = os.Getenv(envvar)
		}

		switch p := v.Elem().Field(i).Interface().(type) {
		case string:
			ptr := v.Elem().FieldByName(target).Addr().Interface().(*string)
			flagset.StringVarP(ptr, longOpt, shortOpt, defval, helpText)
		case bool:
			ptr := v.Elem().FieldByName(target).Addr().Interface().(*bool)
			flagset.BoolVarP(ptr, longOpt, shortOpt, stringToBool(defval), helpText)
		default:
			fmt.Printf("unsupported: %v\n", p)
		}

		if stringToBool(hide) {
			flagset.MarkHidden(longOpt)
		}
	}

	return flagset
}

func stringToBool(s string) bool {
	return slices.Contains([]string{"1", "true"}, strings.ToLower(s))
}
