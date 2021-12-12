package main

import (
	"github.com/tiero/wapc-test/pkg/module"
)

func main() {
	module.Handlers{
		SayHello: sayHello,
	}.Register()
}

func sayHello(name string) (string, error) {
	return "hello" + name, nil
}
