package http

import (
	"github.com/zbum/mantyboot/http/header"
	"github.com/zbum/mantyboot/http/mime"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequestWrapper_ParamInt32(t *testing.T) {
	type fields struct {
		r *http.Request
	}
	type args struct {
		param string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int32
		wantErr bool
	}{
		{
			name: "basic",
			fields: fields{
				httptest.NewRequest(http.MethodGet, "/test?aa=1000000000", nil),
			},
			args: args{
				param: "aa",
			},
			want:    1000000000,
			wantErr: false,
		},

		{
			name: "post",
			fields: fields{
				makeRequest(http.MethodPost, "/test?aa=1000000000", mime.ContentTypeApplicationFormUrlencoded, makePostBody(map[string]string{"aa": "1000"})),
			},
			args: args{
				param: "aa",
			},
			want:    1000000000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := RequestWrapper{
				r: tt.fields.r,
			}
			got, err := w.ParamInt32(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParamInt32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParamInt32() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestWrapper_ParamInt64(t *testing.T) {
	type fields struct {
		r *http.Request
	}
	type args struct {
		param string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "post",
			fields: fields{
				makeRequest(http.MethodPost, "/test?aa=1000000000", mime.ContentTypeApplicationFormUrlencoded, makePostBody(map[string]string{"aa": "1000"})),
			},
			args: args{
				param: "aa",
			},
			want:    1000,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := RequestWrapper{
				r: tt.fields.r,
			}
			got, err := w.ParamInt64(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParamInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParamInt64() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestWrapper_ParseInt(t *testing.T) {
	type fields struct {
		r *http.Request
	}
	type args struct {
		param string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := RequestWrapper{
				r: tt.fields.r,
			}
			got, err := w.ParseInt(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestWrapper_ParseInt16(t *testing.T) {
	type fields struct {
		r *http.Request
	}
	type args struct {
		param string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int16
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := RequestWrapper{
				r: tt.fields.r,
			}
			got, err := w.ParseInt16(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInt16() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRequestWrapper_ParseInt8(t *testing.T) {
	type fields struct {
		r *http.Request
	}
	type args struct {
		param string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int8
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := RequestWrapper{
				r: tt.fields.r,
			}
			got, err := w.ParseInt8(tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInt8() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInt8() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func makeRequest(method string, target string, contentType string, body io.Reader) *http.Request {
	request := httptest.NewRequest(method, target, body)
	request.Header.Set(header.ContentType, contentType)
	return request
}

func makePostBody(values map[string]string) io.Reader {
	data := url.Values{}
	for k, v := range values {
		data.Set(k, v)
	}
	return strings.NewReader(data.Encode())
}
