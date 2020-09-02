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
