package main

import (
	"fmt"
	"github.com/golangit/dic/container"
	"github.com/golangit/dic/reference"
)

func LoggerNew(verbose bool) string {
	return fmt.Sprintf("the verbosity is: %b.", verbose)
}

type TestStruct struct {
	Nerd         string
	DatabaseName string
	Logger       string
}

func main() {

	fmt.Println("Hello, High quality dev")

	cnt := container.New()
	fmt.Println("Registering parameters")
	cnt.Register("dbName", "logger")
	cnt.Register("logger.verbosity", false)

	fmt.Println("Registering services")
	cnt.Register("logger", LoggerNew, reference.New("logger.verbosity"))
	cnt.Register("my_struct", &TestStruct{}, "liuggio", reference.New("dbName"), reference.New("logger"))

	test := cnt.Get("my_struct").(TestStruct)

	fmt.Println("Getting Nerd name: `%s`, the Logger.string is `%s`", test.Nerd, test.Logger)
}
