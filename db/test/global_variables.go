package test

import (
	"auth/db"
	"auth/model"
	"auth/tool/mysqlerr"
	"github.com/jinzhu/gorm"
	"strings"
	"sync"
)

var (
	manager db.AccessorManage
	dbc *gorm.DB
	waitForFinish sync.WaitGroup
)

const numberOfTestFunc = 27

// Hashed Passwords
var passwords = map[string]string{
	"testPW1": "$2a$10$POwSnghOjkriuQ4w1Bj3zeHIGA7fXv8UI/UFXEhnnO5YrcwkUDcXq",
	"testPW2": "$2a$10$XxGXTboHZxhoqzKcBVqkJOiNSy6narAvIQ/ljfTJ4m93jAt8GyX.e",
	"testPW3": "$2a$10$sfZLOR8iVyhXI0y8nXcKIuKseahKu4NLSlocUWqoBdGrpLIZzxJ2S",
}

// FK Constraint Fail Errors
var (
	// StudentAuth 테이블의 ParentUUID 속성의 FK 제약조건 위반에 대한 에러 변수
	studentAuthParentUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Auth_Test_DB"),
		TableName:      model.StudentAuthInstance.TableName(),
		ConstraintName: model.StudentAuthInstance.ParentUUIDConstraintName(),
		AttrName:       model.StudentAuthInstance.ParentUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: model.ParentAuthInstance.TableName(),
		AttrName:  model.ParentAuthInstance.UUID.KeyName(),
	})

	// StudentInform 테이블의 StudentUUID 속성의 FK 제약조건 위반에 대한 에러 변수
	studentInformStudentUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Auth_Test_DB"),
		TableName:      model.StudentInformInstance.TableName(),
		ConstraintName: model.StudentInformInstance.StudentUUIDConstraintName(),
		AttrName:       model.StudentInformInstance.StudentUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: model.StudentAuthInstance.TableName(),
		AttrName:  model.StudentAuthInstance.UUID.KeyName(),
	})

	// TeacherInform 테이블의 TeacherUUID 속성의 FK 제약조건 위반에 대한 에러 변수
	teacherInformTeacherUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Auth_Test_DB"),
		TableName:      model.TeacherInformInstance.TableName(),
		ConstraintName: model.TeacherInformInstance.TeacherUUIDConstraintName(),
		AttrName:       model.TeacherInformInstance.TeacherUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: model.TeacherAuthInstance.TableName(),
		AttrName:  model.TeacherAuthInstance.UUID.KeyName(),
	})

	// TeacherInform 테이블의 TeacherUUID 속성의 FK 제약조건 위반에 대한 에러 변수
	parentInformParentUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Auth_Test_DB"),
		TableName:      model.ParentInformInstance.TableName(),
		ConstraintName: model.ParentInformInstance.ParentUUIDConstraintName(),
		AttrName:       model.ParentInformInstance.ParentUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: model.ParentAuthInstance.TableName(),
		AttrName:  model.ParentAuthInstance.UUID.KeyName(),
	})
)