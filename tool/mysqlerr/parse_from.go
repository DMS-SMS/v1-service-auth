package mysqlerr

import (
	"errors"
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"regexp"
	"strings"
)

var (
	regexForDuplicateError = regexp.MustCompile("'.*?'")
	regexForFKConstraintFail = regexp.MustCompile("`.*?`")
)

func ParseDuplicateEntryErrorFrom(mysqlErr *mysql.MySQLError) (key, entry string, err error) {
	const (
		indexEntry = iota
		indexKey
	)

	if mysqlErr == nil || mysqlErr.Number != mysqlerr.ER_DUP_ENTRY {
		err = errors.New("parameter must be Duplicate Entry Error")
		return
	}

	matched := regexForDuplicateError.FindAllString(mysqlErr.Message, -1)
	if len(matched) != 2 {
		err = errors.New("this parameter is incorrect for the Duplicate Entry Error format")
		return
	}

	for i := range matched {
		matched[i] = strings.Trim(matched[i], "'")
	}

	key = matched[indexKey]
	entry = matched[indexEntry]
	return
}

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

	matched := regexForFKConstraintFail.FindAllString(mysqlErr.Message, -1)
	if len(matched) != 6 {
		err = errors.New("this parameter is incorrect for the FK Construct Fail Error format")
		return
	}

	for i := range matched {
		matched[i] = strings.Trim(matched[i], "`")
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