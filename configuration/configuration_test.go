package configuration

import (
	"reflect"
	"testing"
)

type TestConfiguration struct {
	A string `yaml:"a"`
	B string `yaml:"b"`
}

func TestConfiguration_Load(t *testing.T) {
	type args struct {
		input string
	}
	type testCase[T any] struct {
		name    string
		c       Configuration[T]
		args    args
		wantErr bool
	}
	tests := []testCase[TestConfiguration]{
		{
			name: "simple",
			c:    *NewConfiguration[TestConfiguration](),
			args: args{
				input: `
a:
  atest
b:
  btest
`,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Load(tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfiguration_GetConfiguration(t *testing.T) {
	type testCase[T any] struct {
		name string
		c    Configuration[T]
		want *T
	}

	simple1 := *NewConfiguration[TestConfiguration]()
	simple1.Load(
		`
a: 
  atest
b:
  btest
`)

	tests := []testCase[TestConfiguration]{
		{
			name: "simple1",
			c:    simple1,
			want: &TestConfiguration{A: "atest", B: "btest"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.GetConfiguration(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
