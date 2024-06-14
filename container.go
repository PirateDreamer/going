package going

import "go.uber.org/dig"

var Container *dig.Container

func InitContainer() {
	Container = dig.New()
}
