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

type RefInform struct {
	TableName, AttrName string
}

type FKInform struct {
	DBName, TableName, ConstraintName, AttrName string
}

func FKConstraintFailWithoutReferenceInform(fk FKInform, ref RefInform) *mysql.MySQLError {
	prefix := "Cannot add or update a child row: a foreign key constraint fails"
	suffix := fmt.Sprintf("(`%s`.`%s`, CONSTRAINT `%s` FOREIGN KEY (`%s`) REFERENCES `%s` (`%s`))",
		fk.DBName, fk.TableName, fk.ConstraintName, fk.AttrName, ref.TableName, ref.AttrName)
	return &mysql.MySQLError{
		Number:  mysqlerr.ER_NO_REFERENCED_ROW_2,
		Message: fmt.Sprintf("%s %s", prefix, suffix),
	}
}
