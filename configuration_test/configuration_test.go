package configuration_test

import (
	"embed"
	"github.com/zbum/mantyboot/configuration"
	"reflect"
	"testing"
)

type TestConfiguration struct {
	A string `yaml:"a"`
	B string `yaml:"b"`
}

//go:embed embed/application-dev.yaml
var sampleFs embed.FS

func TestConfiguration_Load1(t *testing.T) {
	type testCase[T any] struct {
		name    string
		c       configuration.Configuration[T]
		want    *T
		wantErr bool
	}
	tests := []testCase[TestConfiguration]{
		{
			name:    "root",
			c:       *configuration.NewConfiguration[TestConfiguration](sampleFs, "dev"),
			want:    &TestConfiguration{A: "aValue", B: "bNewValue"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() got = %v, want %v", got, tt.want)
			}
		})
	}
}
