package http

import (
	"fmt"
	"github.com/zbum/mantyboot/http/header"
	"github.com/zbum/mantyboot/http/mime"
	"golang.org/x/exp/constraints"
	"net/http"
	"strconv"
)

type RequestWrapper struct {
	r *http.Request
}

func NewRequestWrapper(r *http.Request) RequestWrapper {
	return RequestWrapper{r: r}
}

func Parse[T constraints.Signed](param string) (T, error) {
	var result T
	switch any(result).(type) {
	case int:
		i, err := strconv.Atoi(param)
		return T(i), err
	case int8:
		i, err := strconv.ParseInt(param, 10, 8)
		return T(i), err
	case int16:
		i, err := strconv.ParseInt(param, 10, 16)
		return T(i), err
	case int32:
		i, err := strconv.ParseInt(param, 10, 32)
		return T(i), err
	case int64:
		i, err := strconv.ParseInt(param, 10, 64)
		return T(i), err
	}

	return -1, fmt.Errorf("unsupported type for parsing: %T", result)
}

func (w RequestWrapper) ParamInt64(param string) (int64, error) {
	if w.IsPostForm() {
		err := w.r.ParseForm()
		if err != nil {
			return -1, err
		}

		postValue := w.r.PostFormValue(param)
		return Parse[int64](postValue)
	}
	idParam := w.r.URL.Query().Get(param)
	return Parse[int64](idParam)
}

func (w RequestWrapper) ParamInt32(param string) (int32, error) {
	idParam := w.r.URL.Query().Get(param)
	result, err := strconv.ParseInt(idParam, 10, 32)
	return int32(result), err
}

func (w RequestWrapper) ParseInt16(param string) (int16, error) {
	idParam := w.r.URL.Query().Get(param)
	result, err := strconv.ParseInt(idParam, 10, 16)
	return int16(result), err
}

func (w RequestWrapper) ParseInt8(param string) (int8, error) {
	idParam := w.r.URL.Query().Get(param)
	result, err := strconv.ParseInt(idParam, 10, 8)
	return int8(result), err
}

func (w RequestWrapper) ParseInt(param string) (int, error) {
	idParam := w.r.URL.Query().Get(param)
	return strconv.Atoi(idParam)
}

func (w RequestWrapper) IsPostForm() bool {
	contentType := w.r.Header.Get(header.ContentType)
	return (contentType == mime.ContentTypeApplicationFormUrlencoded) && (w.r.Method == http.MethodPost)
}
