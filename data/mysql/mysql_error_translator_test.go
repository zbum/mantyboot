package exception

import (
	"github.com/go-sql-driver/mysql"
	"github.com/zbum/mantyboot/data/support"
	"reflect"
	"testing"
)

func TestMysqlErrorTranslator_TranslateExceptionIfPossible(t1 *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want support.DataAccessError
	}{
		{
			name: "duplicate",
			args: args{
				err: &mysql.MySQLError{
					Number:  1062,
					Message: "duplicate error",
				},
			},
			want: duplicateKeyError,
		},

		{
			name: "fk constraint",
			args: args{
				err: &mysql.MySQLError{
					Number:  1452,
					Message: "fk constraint error",
				},
			},
			want: fkConstraintError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := MysqlErrorTranslator{}
			if got := t.TranslateExceptionIfPossible(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("TranslateExceptionIfPossible() = %v, want %v", got, tt.want)
			}
		})
	}
}
