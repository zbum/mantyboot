# manty-boot
## Configuration
* configuration module like spring-boot's ConfigurationProperties

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
	var ex1Configuration = configuration.NewConfiguration[ExampleConfiguration](devfs, "dev")
	_, err := ex1Configuration.Load()
	if err != nil {
		return
	}
	fmt.Println(ex1Configuration.GetConfiguration())
}
```