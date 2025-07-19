package main

import (
	"embed"
	"fmt"
	"github.com/zbum/mantyboot/configuration"
)

type ExampleConfiguration struct {
	ServerName string `yaml:"server-name"`
	ServerPort string `yaml:"server-port"`
}

//go:embed embed/application-dev.yaml
var devfs embed.FS

func main() {
	ex1Configuration, err := configuration.NewConfiguration[ExampleConfiguration](devfs, "dev")
	if err != nil {
		return
	}
	fmt.Println(ex1Configuration.GetConfiguration())
}
