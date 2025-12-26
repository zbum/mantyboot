# MantyBoot

Go 언어를 위한 Spring Boot 스타일 유틸리티 라이브러리입니다.

## 개요

MantyBoot은 Spring Boot의 편리한 기능들을 Go에서 사용할 수 있도록 만든 라이브러리입니다. 설정 관리, 데이터 접근 에러 처리, HTTP 유틸리티 등의 기능을 제공합니다.

## 모듈

| 모듈 | 설명 | 설치 |
|------|------|------|
| [configuration](#configuration) | Spring Boot의 ConfigurationProperties 스타일 YAML 설정 관리 | `go get github.com/zbum/mantyboot/configuration` |
| [data](#data) | 데이터베이스 에러 추상화 및 번역 | `go get github.com/zbum/mantyboot/data` |
| [http](#http) | HTTP 요청 처리 유틸리티 | `go get github.com/zbum/mantyboot/http` |

## 요구 사항

- Go 1.22 이상

## 설치

```shell
go get github.com/zbum/mantyboot
```

---

### Configuration

Spring Boot의 `ConfigurationProperties`와 유사한 YAML 설정 관리 모듈입니다.

#### 설치

```shell
go get github.com/zbum/mantyboot/configuration
```

#### 설정 파일 로딩 순서

다음 순서로 설정 파일을 찾아 로드합니다. 나중에 로드된 값이 이전 값을 덮어씁니다.

1. embed.FS (임베딩된 파일)
2. `./application-{profile}.yaml`
3. `./config/application-{profile}.yaml`

#### 사용 예시

디렉토리 구조:
```
.
├── application-dev.yaml
├── config
│   └── application-dev.yaml
└── example
    ├── main.go
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
    config, err := configuration.NewConfiguration[ExampleConfiguration](devfs, "dev")
    if err != nil {
        return
    }
    fmt.Println(config.GetConfiguration())
}
```

---

### Data

데이터베이스 에러를 추상화하여 처리하는 모듈입니다.

#### 설치

```shell
go get github.com/zbum/mantyboot/data
```

#### MySQL 에러 번역기

MySQL 에러 코드를 의미 있는 에러 타입으로 변환합니다.

| 에러 코드 | 에러 타입 | 설명 |
|-----------|-----------|------|
| 1062 | `DuplicateKeyError` | 중복 키 에러 |
| 1452 | `FkConstraintError` | 외래 키 제약 조건 위반 |

```go
package main

import (
    "github.com/zbum/mantyboot/data/mysql"
)

func main() {
    translator := mysql.MysqlErrorTranslator{}

    // DB 작업 수행
    err := someDBOperation()

    // 에러 변환
    dataAccessErr := translator.TranslateExceptionIfPossible(err)

    switch dataAccessErr.(type) {
    case mysql.DuplicateKeyError:
        // 중복 키 처리
    case mysql.FkConstraintError:
        // FK 제약 조건 위반 처리
    }
}
```

---

### HTTP

HTTP 요청 처리를 위한 유틸리티 모듈입니다.

#### 설치

```shell
go get github.com/zbum/mantyboot/http
```

#### RequestWrapper

HTTP 요청에서 파라미터를 쉽게 추출할 수 있습니다.

```go
package main

import (
    "net/http"
    mantyhttp "github.com/zbum/mantyboot/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    wrapper := mantyhttp.NewRequestWrapper(r)

    // Query String 또는 POST Form에서 int64 파라미터 추출
    id, err := wrapper.ParamInt64("id")
    if err != nil {
        // 에러 처리
    }

    // 다른 정수 타입도 지원
    count, _ := wrapper.ParamInt32("count")
    page, _ := wrapper.ParseInt("page")
}
```

#### Parse 함수

문자열을 다양한 정수 타입으로 변환하는 제네릭 함수입니다.

```go
import mantyhttp "github.com/zbum/mantyboot/http"

// int64로 변환
val, err := mantyhttp.Parse[int64]("12345")

// int32로 변환
val32, err := mantyhttp.Parse[int32]("123")
```

#### MIME 타입 상수

일반적으로 사용되는 Content-Type 값들을 상수로 제공합니다.

```go
import "github.com/zbum/mantyboot/http/mime"

// 사용 예시
contentType := mime.ContentTypeApplicationJson      // "application/json"
formType := mime.ContentTypeApplicationFormUrlencoded // "application/x-www-form-urlencoded"
```

지원하는 MIME 타입:
- `application/json`
- `application/xml`
- `application/x-www-form-urlencoded`
- `multipart/form-data`
- `text/html`
- `text/plain`
- 기타 다수

---

## 라이선스

Apache License 2.0

자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.
