package test

import (
	"auth/tool/mysqlerr"
	"strings"
)

var (
	// StudentAuth 테이블의 ParentUUID 속성의 FK 제약조건 위반에 대한 에러 변수
	studentAuthParentUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Auth_Test_DB"),
		TableName:      studentAuthModel.TableName(),
		ConstraintName: studentAuthModel.ParentUUIDConstraintName(),
		AttrName:       studentAuthModel.ParentUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: parentAuthModel.TableName(),
		AttrName:  parentAuthModel.UUID.KeyName(),
	})

	// StudentInform 테이블의 StudentUUID 속성의 FK 제약조건 위반에 대한 에러 변수
	studentInformStudentUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Auth_Test_DB"),
		TableName:      studentInformModel.TableName(),
		ConstraintName: studentInformModel.StudentUUIDConstraintName(),
		AttrName:       studentInformModel.StudentUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: studentAuthModel.TableName(),
		AttrName:  studentAuthModel.UUID.KeyName(),
	})

	// TeacherInform 테이블의 TeacherUUID 속성의 FK 제약조건 위반에 대한 에러 변수
	teacherInformTeacherUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Auth_Test_DB"),
		TableName:      teacherInformModel.TableName(),
		ConstraintName: teacherInformModel.TeacherUUIDConstraintName(),
		AttrName:       teacherInformModel.TeacherUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: teacherAuthModel.TableName(),
		AttrName:  teacherAuthModel.UUID.KeyName(),
	})

	// TeacherInform 테이블의 TeacherUUID 속성의 FK 제약조건 위반에 대한 에러 변수
	parentInformParentUUIDFKConstraintFailError = mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
		DBName:         strings.ToLower("SMS_Auth_Test_DB"),
		TableName:      parentInformModel.TableName(),
		ConstraintName: parentInformModel.ParentUUIDConstraintName(),
		AttrName:       parentInformModel.ParentUUID.KeyName(),
	}, mysqlerr.RefInform{
		TableName: parentAuthModel.TableName(),
		AttrName:  parentAuthModel.UUID.KeyName(),
	})
)