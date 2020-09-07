package access

import (
	"auth/adapter"
	"auth/db"
	"auth/model"
	"auth/tool/mysqlerr"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/consul/api"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

var (
	manager db.AccessorManage
	dbc *gorm.DB
)

var passwords = map[string]string{
	"testPW1": "$2a$10$POwSnghOjkriuQ4w1Bj3zeHIGA7fXv8UI/UFXEhnnO5YrcwkUDcXq",
	"testPW2": "$2a$10$XxGXTboHZxhoqzKcBVqkJOiNSy6narAvIQ/ljfTJ4m93jAt8GyX.e",
	"testPW3": "$2a$10$sfZLOR8iVyhXI0y8nXcKIuKseahKu4NLSlocUWqoBdGrpLIZzxJ2S",
}

var (
	studentAuthModel = new(model.StudentAuth)
	teacherAuthModel = new(model.TeacherAuth)
	parentAuthModel = new(model.ParentAuth)
	studentInformModel = new(model.StudentInform)
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
)

func init() {
	cli, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}

	dbc, _, err = adapter.ConnectDBWithConsul(cli, "db/auth/local_test")
	if err != nil {
		log.Fatal(err)
	}
	db.Migrate(dbc)

	manager, err = db.NewAccessorManage(DefaultReflectType(), dbc)
	if err != nil {
		log.Fatal(err)
	}
}

func Test_default_CreateStudentAuth(t *testing.T) {
	// Tx 시작
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}

	// StudentAuth.ParentUUID에 설정할 값을 위한 학부모 계정 생성
	inits := []struct {
		UUID, ParentID, ParentPW string
	} {
		{
			UUID:     "parent-111111111111",
			ParentID: "jinhong07191",
			ParentPW: passwords["testPW1"],
		}, {
			UUID:     "parent-222222222222",
			ParentID: "jinhong07192",
			ParentPW: passwords["testPW2"],
		},
	}

	for _, init := range inits {
		auth := &model.ParentAuth{
			UUID:     model.UUID(init.UUID),
			ParentID: model.ParentID(init.ParentID),
			ParentPW: model.ParentPW(init.ParentPW),
		}
		if _, err := access.CreateParentAuth(auth); err != nil {
			access.Rollback()
			log.Fatal(fmt.Sprintf("error occurs while creating parent auth, err: %v", err))
		}
	}

	tests := []struct {
		UUID, ParentUUID, StudentID, StudentPW string
		ExpectAuth *model.StudentAuth
		ExpectError error
	} {
		{ // success case
			UUID: "student-111111111111",
			ParentUUID: "parent-111111111111",
			StudentID: "jinhong0719",
			StudentPW: passwords["testPW1"],
			ExpectError: nil,
		}, { // UUID duplicate
			UUID: "student-111111111111",
			ParentUUID: "parent-111111111111",
			StudentID: "jinhong07191",
			StudentPW: passwords["testPW1"],
			ExpectError: mysqlerr.DuplicateEntry(studentAuthModel.UUID.KeyName(), "student-111111111111"),
		}, { // StudentID duplicate
			UUID: "student-222222222222",
			ParentUUID: "parent-111111111111",
			StudentID: "jinhong0719",
			StudentPW: passwords["testPW1"],
			ExpectError: mysqlerr.DuplicateEntry(studentAuthModel.StudentID.KeyName(), "jinhong0719"),
		}, { // ParentUUID(foreign key) reference constraint X
			UUID: "student-222222222222",
			ParentUUID: "parent-123412341234", // not exist parent uuid
			StudentID: "jinhong07192",
			StudentPW: passwords["testPW1"],
			ExpectError: studentAuthParentUUIDFKConstraintFailError,
		},
	}

	for _, test := range tests {
		auth := &model.StudentAuth{
			UUID:       model.UUID(test.UUID),
			StudentID:  model.StudentID(test.StudentID),
			StudentPW:  model.StudentPW(test.StudentPW),
			ParentUUID: model.ParentUUID(test.ParentUUID),
		}

		test.ExpectAuth = auth.DeepCopy()
		auth, err := access.CreateStudentAuth(auth)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectAuth, auth.ExceptGormModel(), "result model assertion (test case: %v)", test)
	}

	access.Rollback()
}

func Test_default_CreateParentAuth(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct{
		UUID, ParentId, ParentPw string
		ExpectAuth *model.ParentAuth
		ExpectError error
	} {
		{ // success case
			UUID: "parent-111111111111",
			ParentId: "parent1",
			ParentPw: passwords["testPW1"],
			ExpectError: nil,
		}, { // UUID duplicate
			UUID: "parent-111111111111",
			ParentId: "parent2",
			ParentPw: passwords["testPW2"],
			ExpectError: mysqlerr.DuplicateEntry(parentAuthModel.UUID.KeyName(), "parent-111111111111"),
		}, { // ParentId duplicate
			UUID: "parent-222222222222",
			ParentId: "parent1",
			ParentPw: passwords["testPW2"],
			ExpectError: mysqlerr.DuplicateEntry(parentAuthModel.ParentID.KeyName(), "parent1"),
		},
	}

	for _, test := range tests {
		auth := &model.ParentAuth{
			UUID:     model.UUID(test.UUID),
			ParentID: model.ParentID(test.ParentId),
			ParentPW: model.ParentPW(test.ParentPw),
		}

		test.ExpectAuth = auth.DeepCopy()
		auth, err := access.CreateParentAuth(auth)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectAuth, auth.ExceptGormModel(), "result model assertion (test case: %v)", test)
	}

	access.Rollback()
}

func Test_default_CreateTeacherAuth(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}

	tests := []struct{
		UUID, TeacherID, TeacherPW string
		ExpectAuth *model.TeacherAuth
		ExpectError error
	} {
		{ // success case
			UUID:        "teacher-111111111111",
			TeacherID:   "teacher1",
			TeacherPW:   passwords["testPW1"],
			ExpectError: nil,
		}, { // UUID duplicate
			UUID:        "teacher-111111111111",
			TeacherID:   "teacher2",
			TeacherPW:   passwords["testPW2"],
			ExpectError: mysqlerr.DuplicateEntry(teacherAuthModel.UUID.KeyName(), "teacher-111111111111"),
		}, { // TeacherID duplicate
			UUID:        "teacher-222222222222",
			TeacherID:   "teacher1",
			TeacherPW:   passwords["testPW2"],
			ExpectError: mysqlerr.DuplicateEntry(teacherAuthModel.TeacherID.KeyName(), "teacher1"),
		},
	}

	for _, test := range tests {
		auth := &model.TeacherAuth{
			UUID:      model.UUID(test.UUID),
			TeacherID: model.TeacherID(test.TeacherID),
			TeacherPW: model.TeacherPW(test.TeacherPW),
		}

		test.ExpectAuth = auth.DeepCopy()
		auth, err := access.CreateTeacherAuth(auth)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectAuth, auth.ExceptGormModel(), "result model assertion (test case: %v)", test)
	}

	access.Rollback()
}

func Test_default_CreateStudentInform(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}

	// 학생 계정 생성을 위한 부모님 계정 생성
	for _, init := range []struct {
		UUID, ParentID, ParentPW string
	} {
		{
			UUID:     "parent-111111111111",
			ParentID: "jinhong0719",
			ParentPW: passwords["testPW1"],
		},
	} {
		_, err := access.CreateParentAuth(&model.ParentAuth{
			UUID:     model.UUID(init.UUID),
			ParentID: model.ParentID(init.ParentID),
			ParentPW: model.ParentPW(init.ParentPW),
		})
		if err != nil {
			access.Rollback()
			log.Fatal(fmt.Sprintf("error occurs while creating parent auth, err: %v", err))
		}
	}

	// 학생 정보 생성을 위한 학생 계정 생성
	for _, init := range []struct {
		UUID, StudentID, StudentPW, ParentUUID string
	} {
		{
			UUID:       "student-111111111111",
			StudentID:  "jinhong07191",
			StudentPW: passwords["testPW1"],
			ParentUUID: "parent-111111111111",
		}, {
			UUID:       "student-222222222222",
			StudentID:  "jinhong07192",
			StudentPW: passwords["testPW1"],
			ParentUUID: "parent-111111111111",
		},
	} {
		_, err := access.CreateStudentAuth(&model.StudentAuth{
			UUID:       model.UUID(init.UUID),
			StudentID:  model.StudentID(init.StudentID),
			StudentPW:  model.StudentPW(init.StudentPW),
			ParentUUID: model.ParentUUID(init.ParentUUID),
		})
		if err != nil {
			access.Rollback()
			log.Fatal(fmt.Sprintf("error occurs while creating student auth, err: %v", err))
		}
	}

	tests := []struct {
		StudentUUID, Name, PhoneNumber, ProfileURI string
		Grade, Class, StudentNumber                int64
		ExpectResult                               *model.StudentInform
		ExpectError                                error
	} {
		{ // success case
			StudentUUID:   "student-111111111111",
			Grade:         2,
			Class:         2,
			StudentNumber: 7,
			Name:          "박진홍",
			PhoneNumber:   "01088378347",
			ProfileURI:    "example.com/profiles/student-111111111111",
			ExpectError:   nil,
		}, { // student uuid duplicate
			StudentUUID:   "student-111111111111",
			Grade:         1,
			Class:         2,
			StudentNumber: 8,
			Name:          "빡진홍",
			PhoneNumber:   "01012341234",
			ProfileURI:    "example.com/profiles/student-111111111112",
			ExpectError:   mysqlerr.DuplicateEntry(studentInformModel.StudentUUID.KeyName(), "student-111111111111"),
		}, { // student number duplicate
			StudentUUID:   "student-222222222222",
			Grade:         2,
			Class:         2,
			StudentNumber: 7,
			Name:          "빡진홍",
			PhoneNumber:   "01012341234",
			ProfileURI:    "example.com/profiles/student-222222222222",
			ExpectError:   mysqlerr.DuplicateEntry(studentInformModel.StudentNumber.KeyName(), "2-2-07"),
		}, { // phone number duplicate
			StudentUUID:   "student-222222222222",
			Grade:         1,
			Class:         2,
			StudentNumber: 8,
			Name:          "빡진홍",
			PhoneNumber:   "01088378347",
			ProfileURI:    "example.com/profiles/student-222222222222",
			ExpectError:   mysqlerr.DuplicateEntry(studentInformModel.PhoneNumber.KeyName(), "01088378347"),
		}, { // profile uri duplicate
			StudentUUID:   "student-222222222222",
			Grade:         1,
			Class:         2,
			StudentNumber: 8,
			Name:          "빡진홍",
			PhoneNumber:   "01012341234",
			ProfileURI:    "example.com/profiles/student-111111111111",
			ExpectError:   mysqlerr.DuplicateEntry(studentInformModel.ProfileURI.KeyName(), "example.com/profiles/student-111111111111"),
		}, { // student uuid FK constraint fail
			StudentUUID:   "student-333333333333",
			Grade:         1,
			Class:         2,
			StudentNumber: 8,
			Name:          "빡진홍",
			PhoneNumber:   "01012341234",
			ProfileURI:    "example.com/profiles/student-333333333333",
			ExpectError:   studentInformStudentUUIDFKConstraintFailError,
		},
	}

	for _, test := range tests {
		inform := &model.StudentInform{
			StudentUUID:   model.StudentUUID(test.StudentUUID),
			Grade:         model.Grade(test.Grade),
			Class:         model.Class(test.Class),
			StudentNumber: model.StudentNumber(test.StudentNumber),
			Name:          model.Name(test.Name),
			PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
			ProfileURI:    model.ProfileURI(test.ProfileURI),
		}

		test.ExpectResult = inform.DeepCopy()
		result, err := access.CreateStudentInform(inform)

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			err = mysqlerr.ExceptReferenceInformFrom(mysqlErr)
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result model assertion error (test case: %v)", test)
	}

	access.Commit()
}

func TestDBClose(t *testing.T) {
	_ = dbc.Close()
}

