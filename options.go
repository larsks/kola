package main

import (
	"fmt"
	"reflect"
	"strings"

	flag "github.com/spf13/pflag"
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
	ValidationError ApplicationErrorType = ApplicationErrorType{
		Message: "Invalid option value",
		Parent:  ApplicationError,
	}
)

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
		//validator := v.MethodByName(fmt.Sprintf("Validate%s", field.Name))

		switch p := v.Elem().Field(i).Interface().(type) {
		case string:
			ptr := v.Elem().FieldByName(target).Addr().Interface().(*string)
			flagset.StringVarP(ptr, longOpt, shortOpt, "", helpText)
		case bool:
			ptr := v.Elem().FieldByName(target).Addr().Interface().(*bool)
			flagset.BoolVarP(ptr, longOpt, shortOpt, false, helpText)
		default:
			fmt.Printf("wtf: %v\n", p)
		}
	}

	return flagset
}
