# manty-boot
## Configuration
* configuration module like spring-boot's ConfigurationProperties

### Install
```shell
go get github.com/zbum/mantyboot/configuration
```
### Usage
* If you have file structure like below, It loads files order by (embedfs -> ./application-dev.yaml -> ./config/application-dev.yaml) 
```
.
├── application-dev.yaml
├── config
│   └── application-dev.yaml
└── example
    ├── configuration_ex1.go
    └── embed
        └── application-dev.yaml
```

```go
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

```