package test

import (
	"auth/model"
	"auth/tool/mysqlerr"
	"strings"
)

var (
	StudentAuthParentUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Auth_Test_DB"),
		TableName:      model.StudentAuthInstance.TableName(),
		ConstraintName: model.StudentAuthInstance.ParentUUIDConstraintName(),
		AttrName:       model.StudentAuthInstance.ParentUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: model.ParentAuthInstance.TableName(),
		AttrName:  model.ParentAuthInstance.UUID.KeyName(),
	})
)
