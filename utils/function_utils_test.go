package utils

import (
	"net/http"
	"testing"
)

func sampleHandler(w *http.ResponseWriter, r http.Request) string {
	// do nothing
	return "test"
}

func TestGetFunctionName(t *testing.T) {
	type args struct {
		f any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "sampleHandler",
			args: args{
				f: sampleHandler,
			},
			want: "github.com/zbum/mantyboot/utils.sampleHandler",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFunctionName(tt.args.f); got != tt.want {
				t.Errorf("GetFunctionName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSignature(t *testing.T) {
	type args struct {
		f any
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "sampleHandler",
			args: args{
				f: sampleHandler,
			},
			want: "func (*http.ResponseWriter, http.Request) string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetSignature(tt.args.f); got != tt.want {
				t.Errorf("GetSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
