# manty-boot

## Install
```shell
go get github.com/zbum/mantyboot/configuration
```

## Mux
* The first goal of MantyBoot's Mux is compatibility to standard library.

### Usage

```go
package main

import (
	"log"
	"net/http"

	"github.com/zbum/mantyboot/http/mux"
	"github.com/zbum/mantyboot/http/mux/middleware"
)

func main() {

	mux := mux.NewMantyMux()

	mux.AddMiddleware(middleware.AccessLogger(log.Default()))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test"))
	})

	mux.HandleFunc("GET /manty", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("manty"))
	})

	log.Panic(http.ListenAndServe(":8080", mux))
}

```


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
    ex1Configuration, err := configuration.NewConfiguration[ExampleConfiguration](devfs, "dev")
    if err != nil {
        return
    }
    fmt.Println(ex1Configuration.GetConfiguration())
}

```