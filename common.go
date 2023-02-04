package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var container = make([]*dependency, 0)

type dependency struct {
	name     string
	instance any
	factory  func() any
	zero     interface{}
}

func (d dependency) getInstance() any {
	if d.instance == nil {
		d.instance = d.factory()
	}

	return d.instance
}

func Define(factory func() any) error {
	var name = runtime.FuncForPC(reflect.ValueOf(factory).Pointer()).Name()
	sl := strings.Split(name, ".")
	name = strings.Replace(sl[len(sl)-1], "New", "", 1)

	for _, x := range container {
		if x.name == name {
			return fmt.Errorf("bean with name %s already defined", name)
		}
	}

	tt := reflect.TypeOf(factory).Out(0)

	container = append(
		container,
		&dependency{name: name, factory: factory, zero: reflect.Zero(tt).Interface()},
	)

	return nil
}

func Get[T any](qualifier ...string) (T, error) {
	var defaultVal T
	var pretenders = make([]*dependency, 0)

	for _, x := range container {
		if _, ok := x.zero.(T); ok {
			pretenders = append(pretenders, x)
		}
	}

	if len(pretenders) > 1 {
		if len(qualifier) > 0 {
			for _, x := range pretenders {
				if x.name == qualifier[0] {
					return x.getInstance().(T), nil
				}
			}
		}

		return defaultVal, fmt.Errorf("found %d pretenders for injection, %v", len(pretenders), pretenders)
	}

	if len(pretenders) == 0 {
		return defaultVal, fmt.Errorf("not found bean")
	}

	return pretenders[0].getInstance().(T), nil
}
