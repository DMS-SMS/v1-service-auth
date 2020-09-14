package mysqlerr

import (
	"errors"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
)

func ParseFKConstraintFailErrorFrom(mysqlErr *mysql.MySQLError) (fk FKInform, ref RefInform, err error) {
	const (
		dbNameIndex = iota
		tableNameIndex
		constraintNameIndex
		attrNameIndex
		refTableNameIndex
		refAttrNameIndex
	)

	if mysqlErr == nil || mysqlErr.Number != mysqlerr.ER_NO_REFERENCED_ROW_2 {
		err = errors.New("parameter must be an FK Construct Fail Error")
		return
	}

	regex := regexp.MustCompile("`.*?`")
	matched := regex.FindAllString(mysqlErr.Message, -1)
	for i := range matched {
		matched[i] = strings.Trim(matched[i], "`")
	}

	if len(matched) != 6 {
		err = errors.New("this parameter is incorrect for the FK Construct Fail Error format")
		return
	}

	fk = FKInform{
		DBName:         matched[dbNameIndex],
		TableName:      matched[tableNameIndex],
		ConstraintName: matched[constraintNameIndex],
		AttrName:       matched[attrNameIndex],
	}
	ref = RefInform{
		TableName: matched[refTableNameIndex],
		AttrName:  matched[refAttrNameIndex],
	}
	return
}