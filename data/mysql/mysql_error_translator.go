package exception

import (
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/zbum/mantyboot/data/support"
)

var duplicateKeyError = DuplicateKeyError{}
var fkConstraintError = FkConstraintError{}

type DuplicateKeyError struct {
}

func (d DuplicateKeyError) Error() string {
	return "[mysql]conflict"
}

type FkConstraintError struct {
}

func (d FkConstraintError) Error() string {
	return "[mysql]fk constraint"
}

type MysqlErrorTranslator struct {
}

func (t MysqlErrorTranslator) TranslateExceptionIfPossible(err error) support.DataAccessError {
	if err != nil {
		if mySQLErrorCode(err) == 1062 {
			return duplicateKeyError
		}

		if mySQLErrorCode(err) == 1452 {
			return fkConstraintError
		}
		return err
	}
	return nil
}

func mySQLErrorCode(err error) uint16 {
	var val *mysql.MySQLError
	if errors.As(err, &val) {
		return val.Number
	}
	return 0
}
