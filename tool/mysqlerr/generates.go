package mysqlerr

import (
	"fmt"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

func DuplicateEntry(key, entry string) *mysql.MySQLError {
	return &mysql.MySQLError{
		Number:  mysqlerr.ER_DUP_ENTRY,
		Message: fmt.Sprintf("Duplicate entry '%s' for key '%s'", entry, key),
	}
}

type Reference struct {
	TableName, AttrName string
}

func FKConstraintFail(dbName, tableName, constraintName, attrName string, ref Reference) *mysql.MySQLError {
	prefix := "Cannot add or update a child row: a foreign key constraint fails"
	suffix := fmt.Sprintf("(`%s`.`%s`, CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`) ON DELETE RESTRICT ON UPDATE)",
		dbName, tableName, constraintName, attrName, ref.TableName, ref.AttrName)
	return &mysql.MySQLError{
		Number:  mysqlerr.ER_NO_REFERENCED_ROW_2,
		Message: prefix + " " + suffix,
	}
}
