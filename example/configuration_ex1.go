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
	var ex1Configuration = configuration.NewConfiguration[ExampleConfiguration](devfs, "dev")
	_, err := ex1Configuration.Load()
	if err != nil {
		return
	}
	fmt.Println(ex1Configuration.GetConfiguration())
}
