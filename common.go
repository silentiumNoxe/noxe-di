package di

import (
	"fmt"
)

var container = make([]*dependency, 0)

type dependency struct {
	name     string
	instance any
}

func Define(name string, value any) {
	for _, x := range container {
		if x.name == name {
			panic(fmt.Errorf("Dependency with name %s already exist ", name))
		}
	}

	container = append(container, &dependency{name, value})
}

func Get[T any](qualifier ...string) T {
	var pretenders = make([]*dependency, 0)

	for _, x := range container {
		if _, ok := x.instance.(T); ok {
			pretenders = append(pretenders, x)
		}
	}

	if len(pretenders) > 1 {
		if len(qualifier) > 0 {
			for _, x := range pretenders {
				if x.name == qualifier[0] {
					return x.instance.(T)
				}
			}
		}

		panic(fmt.Errorf("found %d pretenders for injection, %v", len(pretenders), pretenders))
	}

	if len(pretenders) == 0 {
		panic(fmt.Errorf("not found dependency (qualifier: %s)", qualifier[0]))
	}

	return pretenders[0].instance.(T)
}
