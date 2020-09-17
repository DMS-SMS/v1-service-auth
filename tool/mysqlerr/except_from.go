package mysqlerr

import (
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
)

func ExceptReferenceInformFrom(mysqlErr *mysql.MySQLError) *mysql.MySQLError {
	switch mysqlErr.Number {
	case mysqlerr.ER_NO_REFERENCED_ROW_2:
		fkInform, refInform, parseErr := ParseFKConstraintFailErrorFrom(mysqlErr)
		if parseErr != nil {
			return mysqlErr
		}
		mysqlErr = FKConstraintFailWithoutReferenceInform(fkInform, refInform)
	}
	return mysqlErr
}
