package test

import (
	"auth/db/access/errors"
	"auth/model"
	"auth/tool/mysqlerr"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_Access_ModifyStudentInform(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = access.Rollback()
		waitForFinish.Done()
	}()

	// 학부모 계정 생성
	for _, init := range []struct {
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
		},{
			UUID:     "parent-333333333333",
			ParentID: "jinhong07193",
			ParentPW: passwords["testPW1"],
		},
	} {
		_, err := access.CreateParentAuth(&model.ParentAuth{
			UUID:     model.UUID(init.UUID),
			ParentID: model.ParentID(init.ParentID),
			ParentPW: model.ParentPW(init.ParentPW),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating parent auth, err: %v", err))
		}
	}

	// 학생 계정 생성
	for _, init := range []struct {
		UUID, StudentID, StudentPW, ParentUUID string
	} {
		{
			UUID:       "student-111111111111",
			StudentID:  "jinhong07191",
			StudentPW:  passwords["testPW1"],
			ParentUUID: "parent-111111111111",
		}, {
			UUID:       "student-222222222222",
			StudentID:  "jinhong07192",
			StudentPW:  passwords["testPW2"],
			ParentUUID: "parent-222222222222",
		}, {
			UUID:       "student-333333333333",
			StudentID:  "jinhong07193",
			StudentPW:  passwords["testPW1"],
			ParentUUID: "parent-333333333333",
		},
	} {
		_, err := access.CreateStudentAuth(&model.StudentAuth{
			UUID:       model.UUID(init.UUID),
			StudentID:  model.StudentID(init.StudentID),
			StudentPW:  model.StudentPW(init.StudentPW),
			ParentUUID: model.ParentUUID(init.ParentUUID),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating student auth, err: %v", err))
		}
	}

	// 학생 정보 생성
	for _, init := range []struct {
		StudentUUID, Name           string
		PhoneNumber, ProfileURI     string
		Grade, Class, StudentNumber int64
	} {
		{
			StudentUUID: "student-111111111111",
			Grade: 2,
			Class: 2,
			StudentNumber: 7,
			Name: "박진홍",
			PhoneNumber: "01011111111",
			ProfileURI: "example.com/profiles/student-111111111111",
		}, {
			StudentUUID: "student-222222222222",
			Grade: 2,
			Class: 2,
			StudentNumber: 12,
			Name: "오준상",
			PhoneNumber: "01022222222",
			ProfileURI: "example.com/profiles/student-222222222222",
		}, {
			StudentUUID: "student-333333333333",
			Grade: 2,
			Class: 2,
			StudentNumber: 14,
			Name: "윤석준",
			PhoneNumber: "01033333333",
			ProfileURI: "example.com/profiles/student-333333333333",
		},
	} {
		_, err := access.CreateStudentInform(&model.StudentInform{
			StudentUUID:   model.StudentUUID(init.StudentUUID),
			Grade:         model.Grade(init.Grade),
			Class:         model.Class(init.Class),
			StudentNumber: model.StudentNumber(init.StudentNumber),
			Name:          model.Name(init.Name),
			PhoneNumber:   model.PhoneNumber(init.PhoneNumber),
			ProfileURI:    model.ProfileURI(init.ProfileURI),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating student inform, err: %v", err))
		}
	}

	tests := []struct {
		StudentUUIDForArgs, StudentUUID string
		Grade, Class, StudentNumber     int64
		Name, PhoneNumber, ProfileURI   string
		ExpectError                     error
	} {
		{ // success case 1 (about int64 field)
			StudentUUIDForArgs: "student-111111111111",
			Grade:              3,
			Class:              2,
			StudentNumber:      8,
			ExpectError:        nil,
		}, { // success case 2 (about string field)
			StudentUUIDForArgs: "student-222222222222",
			Name:               "오줌상",
			PhoneNumber:        "01044444444",
			ProfileURI:         "example.com/profiles/student/student-222222222222",
			ExpectError:        nil,
		}, { // student number duplicate error
			StudentUUIDForArgs: "student-333333333333",
			Grade:              2,
			Class:              2,
			StudentNumber:      14,
			ExpectError:        mysqlerr.DuplicateEntry(model.StudentInformInstance.StudentNumber.KeyName(), "2214"),
		}, { // student number duplicate error
			StudentUUIDForArgs: "student-333333333333",
			PhoneNumber:        "01011111111",
			ExpectError:        mysqlerr.DuplicateEntry(model.StudentInformInstance.PhoneNumber.KeyName(), "01011111111"),
		}, { // student uuid cannot be changed error
			StudentUUIDForArgs: "student-333333333333",
			StudentUUID:        "student-444444444444",
			ExpectError:        errors.StudentUUIDCannotBeChanged,
		}, { // no exist student uuid -> nil error return!
			StudentUUIDForArgs: "student-4444444444444444",
			StudentNumber:      1,
			ExpectError:        nil,
		},
		// 도메인 밖의 값이라면 어떻계?
	}

	for _, test := range tests {
		revisionInform := &model.StudentInform{
			StudentUUID:   model.StudentUUID(test.StudentUUID),
			Grade:         model.Grade(test.Grade),
			Class:         model.Class(test.Class),
			StudentNumber: model.StudentNumber(test.StudentNumber),
			Name:          model.Name(test.Name),
			PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
			ProfileURI:    model.ProfileURI(test.ProfileURI),
		}
		err := access.ModifyStudentInform(test.StudentUUIDForArgs, revisionInform)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
	}

	testsForConfirmModify := []struct {
		StudentUUIDArgs, StudentUUID  string
		Grade, Class, StudentNumber   int64
		Name, PhoneNumber, ProfileURI string
		ExpectError                   error
	} {
		{
			StudentUUIDArgs: "student-111111111111",
			StudentUUID:     "student-111111111111",
			Grade:           3,
			Class:           2,
			StudentNumber:   8,
			Name:            "박진홍",
			PhoneNumber:     "01011111111",
			ProfileURI:      "example.com/profiles/student-111111111111",
			ExpectError:     nil,
		}, {
			StudentUUIDArgs: "student-222222222222",
			StudentUUID:     "student-222222222222",
			Grade:           2,
			Class:           2,
			StudentNumber:   12,
			Name:            "오줌상",
			PhoneNumber:     "01044444444",
			ProfileURI:      "example.com/profiles/student/student-222222222222",
			ExpectError:     nil,
		}, {
			StudentUUIDArgs: "student-333333333333",
			StudentUUID:     "student-333333333333",
			Grade:           2,
			Class:           2,
			StudentNumber:   14,
			Name:            "윤석준",
			PhoneNumber:     "01033333333",
			ProfileURI:      "example.com/profiles/student-333333333333",
			ExpectError:     nil,
		},
	}

	for _, test := range testsForConfirmModify {
		expectResult := &model.StudentInform{
			StudentUUID:   model.StudentUUID(test.StudentUUID),
			Grade:         model.Grade(test.Grade),
			Class:         model.Class(test.Class),
			StudentNumber: model.StudentNumber(test.StudentNumber),
			Name:          model.Name(test.Name),
			PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
			ProfileURI:    model.ProfileURI(test.ProfileURI),
		}
		resultInform, err := access.GetStudentInformWithUUID(test.StudentUUIDArgs)


		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, expectResult, resultInform.ExceptGormModel(), "result inform model assertion error (test case: %v)", test)
	}
}