package mysql

import (
	"errors"
	"fmt"
	"strings"
	"github.com/go-sql-driver/mysql"
	"github.com/zbum/mantyboot/data/support"
	merrors "github.com/zbum/mantyboot/errors"
)

type DuplicateKeyError struct {
	Table   string
	Column  string
	Message string
}

func (d DuplicateKeyError) Error() string {
	if d.Table != "" && d.Column != "" {
		return fmt.Sprintf("duplicate key error on table '%s' column '%s': %s", d.Table, d.Column, d.Message)
	}
	return fmt.Sprintf("duplicate key error: %s", d.Message)
}

type FkConstraintError struct {
	Table      string
	Constraint string
	Message    string
}

func (d FkConstraintError) Error() string {
	if d.Table != "" && d.Constraint != "" {
		return fmt.Sprintf("foreign key constraint error on table '%s' constraint '%s': %s", d.Table, d.Constraint, d.Message)
	}
	return fmt.Sprintf("foreign key constraint error: %s", d.Message)
}

type ConnectionError struct {
	Message string
}

func (d ConnectionError) Error() string {
	return fmt.Sprintf("database connection error: %s", d.Message)
}

type SyntaxError struct {
	Message string
}

func (d SyntaxError) Error() string {
	return fmt.Sprintf("SQL syntax error: %s", d.Message)
}

type MysqlErrorTranslator struct {
}

func (t MysqlErrorTranslator) TranslateExceptionIfPossible(err error) support.DataAccessError {
	if err == nil {
		return nil
	}

	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return merrors.WrapDatabaseError(err, "unknown", "failed to translate MySQL error")
	}

	switch mysqlErr.Number {
	case 1062: // Duplicate entry
		return DuplicateKeyError{
			Table:   extractTableFromMessage(mysqlErr.Message),
			Column:  extractColumnFromMessage(mysqlErr.Message),
			Message: mysqlErr.Message,
		}
	case 1452: // Cannot add or update a child row: a foreign key constraint fails
		return FkConstraintError{
			Table:      extractTableFromMessage(mysqlErr.Message),
			Constraint: extractConstraintFromMessage(mysqlErr.Message),
			Message:    mysqlErr.Message,
		}
	case 2002, 2003, 2006, 2013: // Connection errors
		return ConnectionError{
			Message: mysqlErr.Message,
		}
	case 1064, 1146, 1054: // Syntax errors
		return SyntaxError{
			Message: mysqlErr.Message,
		}
	default:
		return merrors.WrapDatabaseError(err, "unknown", fmt.Sprintf("unhandled MySQL error %d", mysqlErr.Number))
	}
}

func extractTableFromMessage(message string) string {
	// Simple extraction - can be enhanced with regex
	if len(message) > 0 {
		// Look for table name in quotes
		start := strings.Index(message, "`")
		if start != -1 {
			end := strings.Index(message[start+1:], "`")
			if end != -1 {
				return message[start+1 : start+1+end]
			}
		}
	}
	return ""
}

func extractColumnFromMessage(message string) string {
	// Simple extraction - can be enhanced with regex
	if len(message) > 0 {
		// Look for column name in quotes
		start := strings.LastIndex(message, "`")
		if start != -1 {
			end := strings.Index(message[start+1:], "`")
			if end != -1 {
				return message[start+1 : start+1+end]
			}
		}
	}
	return ""
}

func extractConstraintFromMessage(message string) string {
	// Simple extraction - can be enhanced with regex
	if len(message) > 0 {
		// Look for constraint name
		start := strings.Index(message, "constraint")
		if start != -1 {
			// Extract constraint name after "constraint"
			parts := strings.Split(message[start:], " ")
			if len(parts) > 1 {
				return strings.Trim(parts[1], "`'")
			}
		}
	}
	return ""
}
