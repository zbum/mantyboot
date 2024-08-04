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

func TestNewConfiguration(t *testing.T) {
	type args struct {
		embedDir embed.FS
		profile  string
	}
	type testCase[T any] struct {
		name    string
		args    args
		want    *T
		wantErr bool
	}
	tests := []testCase[TestConfiguration]{
		{
			name: "simple",
			args: args{
				embedDir: sampleFs,
				profile:  "dev",
			},
			want:    &TestConfiguration{A: "aValue", B: "bNewValue"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := configuration.NewConfiguration[TestConfiguration](tt.args.embedDir, tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfiguration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.GetConfiguration(), tt.want) {
				t.Errorf("NewConfiguration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
