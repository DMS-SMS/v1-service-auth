package test

import (
	"auth/model"
	"auth/tool/mysqlerr"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_Accessor_CreateStudentAuth(t *testing.T) {
	// Tx 시작
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
	}()

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
			ExpectError: mysqlerr.DuplicateEntry(model.StudentAuthInstance.UUID.KeyName(), "student-111111111111"),
		}, { // StudentID duplicate
			UUID: "student-222222222222",
			ParentUUID: "parent-111111111111",
			StudentID: "jinhong0719",
			StudentPW: passwords["testPW1"],
			ExpectError: mysqlerr.DuplicateEntry(model.StudentAuthInstance.StudentID.KeyName(), "jinhong0719"),
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

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			err = mysqlerr.ExceptReferenceInformFrom(mysqlErr)
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectAuth, auth.ExceptGormModel(), "result model assertion (test case: %v)", test)
	}
}

func Test_Accessor_CreateParentAuth(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
	}()

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
			ExpectError: mysqlerr.DuplicateEntry(model.ParentAuthInstance.UUID.KeyName(), "parent-111111111111"),
		}, { // ParentId duplicate
			UUID: "parent-222222222222",
			ParentId: "parent1",
			ParentPw: passwords["testPW2"],
			ExpectError: mysqlerr.DuplicateEntry(model.ParentAuthInstance.ParentID.KeyName(), "parent1"),
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
}

func Test_Accessor_CreateTeacherAuth(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
	}()

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
			ExpectError: mysqlerr.DuplicateEntry(model.TeacherAuthInstance.UUID.KeyName(), "teacher-111111111111"),
		}, { // TeacherID duplicate
			UUID:        "teacher-222222222222",
			TeacherID:   "teacher1",
			TeacherPW:   passwords["testPW2"],
			ExpectError: mysqlerr.DuplicateEntry(model.TeacherAuthInstance.TeacherID.KeyName(), "teacher1"),
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
}

func Test_Accessor_CreateStudentInform(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
	}()

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
			ExpectError:   mysqlerr.DuplicateEntry(model.StudentInformInstance.StudentUUID.KeyName(), "student-111111111111"),
		}, { // student number duplicate
			StudentUUID:   "student-222222222222",
			Grade:         2,
			Class:         2,
			StudentNumber: 7,
			Name:          "빡진홍",
			PhoneNumber:   "01012341234",
			ProfileURI:    "example.com/profiles/student-222222222222",
			ExpectError:   mysqlerr.DuplicateEntry(model.StudentInformInstance.StudentNumber.KeyName(), "2207"),
		}, { // phone number duplicate
			StudentUUID:   "student-222222222222",
			Grade:         1,
			Class:         2,
			StudentNumber: 8,
			Name:          "빡진홍",
			PhoneNumber:   "01088378347",
			ProfileURI:    "example.com/profiles/student-222222222222",
			ExpectError:   mysqlerr.DuplicateEntry(model.StudentInformInstance.PhoneNumber.KeyName(), "01088378347"),
		}, { // profile uri duplicate
			StudentUUID:   "student-222222222222",
			Grade:         1,
			Class:         2,
			StudentNumber: 8,
			Name:          "빡진홍",
			PhoneNumber:   "01012341234",
			ProfileURI:    "example.com/profiles/student-111111111111",
			ExpectError:   mysqlerr.DuplicateEntry(model.StudentInformInstance.ProfileURI.KeyName(), "example.com/profiles/student-111111111111"),
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
}

func Test_Accessor_CreateTeacherInform(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
	}()

	for _, init := range []struct {
		UUID, TeacherID, TeacherPW string
	} {
		{
			UUID:      "teacher-111111111111",
			TeacherID: "jinhong07191",
			TeacherPW: passwords["testPW1"],
		}, {
			UUID:      "teacher-222222222222",
			TeacherID: "jinhong07192",
			TeacherPW: passwords["testPW2"],
		},
	} {
		_, err := access.CreateTeacherAuth(&model.TeacherAuth{
			UUID:      model.UUID(init.UUID),
			TeacherID: model.TeacherID(init.TeacherID),
			TeacherPW: model.TeacherPW(init.TeacherPW),
		})
		if err != nil {
			access.Rollback()
			log.Fatal(fmt.Sprintf("error occurs while creating teacher auth. err: %v", err))
		}
	}

	tests := []struct {
		TeacherUUID, Name, PhoneNumber string
		Grade, Class                   int64
		ExpectResult                   *model.TeacherInform
		ExpectError                    error
	} {
		{ // success test case
			TeacherUUID: "teacher-111111111111",
			Grade:       2,
			Class:       2,
			Name:        "박진홍",
			PhoneNumber: "01088378347",
			ExpectError: nil,
		}, { // teacher uuid duplicate
			TeacherUUID: "teacher-111111111111",
			Grade:       1,
			Class:       2,
			Name:        "박진홍",
			PhoneNumber: "01012341234",
			ExpectError: mysqlerr.DuplicateEntry(model.TeacherInformInstance.TeacherUUID.KeyName(), "teacher-111111111111"),
		}, { // phone number duplicate
			TeacherUUID: "teacher-222222222222",
			Grade:       1,
			Class:       2,
			Name:        "박진홍",
			PhoneNumber: "01088378347",
			ExpectError: mysqlerr.DuplicateEntry(model.TeacherInformInstance.PhoneNumber.KeyName(), "01088378347"),
		}, { // teacher uuid fk constraint fail
			TeacherUUID: "teacher-333333333333",
			Grade:       1,
			Class:       2,
			Name:        "박진홍",
			PhoneNumber: "01012341234",
			ExpectError: teacherInformTeacherUUIDFKConstraintFailError,
		},
	}

	for _, test := range tests {
		inform := &model.TeacherInform{
			TeacherUUID: model.TeacherUUID(test.TeacherUUID),
			Grade:       model.Grade(test.Grade),
			Class:       model.Class(test.Class),
			Name:        model.Name(test.Name),
			PhoneNumber: model.PhoneNumber(test.PhoneNumber),
		}

		test.ExpectResult = inform.DeepCopy()
		result, err := access.CreateTeacherInform(inform)

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			err = mysqlerr.ExceptReferenceInformFrom(mysqlErr)
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result model assertion error (test case: %v)", test)
	}
}

func Test_Accessor_CreateParentInform(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
	}()

	for _, init := range []struct{
		UUID, ParentID, ParentPW string
	}{
		{
			UUID: "parent-111111111111",
			ParentID: "jinhong07191",
			ParentPW: passwords["testPW1"],
		}, {
			UUID: "parent-222222222222",
			ParentID: "jinhong07192",
			ParentPW: passwords["testPW2"],
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

	tests := []struct {
		ParentUUID, Name, PhoneNumber string
		ExpectResult                  *model.ParentInform
		ExpectError                   error
	} {
		{ // success case
			ParentUUID:  "parent-111111111111",
			Name:        "박진홍",
			PhoneNumber: "01088378347",
			ExpectError: nil,
		}, { // parent uuid duplicate error
			ParentUUID:  "parent-111111111111",
			Name:        "박진홍",
			PhoneNumber: "01012341234",
			ExpectError: mysqlerr.DuplicateEntry(model.ParentInformInstance.ParentUUID.KeyName(), "parent-111111111111"),
		}, { // phone number duplicate error
			ParentUUID:  "parent-222222222222",
			Name:        "박진홍",
			PhoneNumber: "01088378347",
			ExpectError: mysqlerr.DuplicateEntry(model.ParentInformInstance.PhoneNumber.KeyName(), "01088378347"),
		}, { // parent uuid FK constraint fail
			ParentUUID:  "parent-333333333333",
			Name:        "박진홍",
			PhoneNumber: "01012341234",
			ExpectError: parentInformParentUUIDFKConstraintFailError,
		},
	}

	for _, test := range tests {
		inform := &model.ParentInform{
			ParentUUID:  model.ParentUUID(test.ParentUUID),
			Name:        model.Name(test.Name),
			PhoneNumber: model.PhoneNumber(test.PhoneNumber),
		}

		test.ExpectResult = inform.DeepCopy()
		result, err := access.CreateParentInform(inform)

		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			err = mysqlerr.ExceptReferenceInformFrom(mysqlErr)
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectResult, result.ExceptGormModel(), "result model assertion error (test case: %v)", test)
	}
}
