package test

import (
	"auth/model"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_Access_DeleteStudentAuth(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Commit()
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

	tests := []struct {
		StudentUUIDForArgs string
		ExpectError        error
	} {
		{ // success case
			StudentUUIDForArgs: "student-111111111111",
			ExpectError:        nil,
		}, { // no exist student uuid -> not error!
			StudentUUIDForArgs: "student-222222222222",
			ExpectError:        nil,
		},
	}

	for _, test := range tests {
		err := access.DeleteStudentAuth(test.StudentUUIDForArgs)
		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
	}

	//testForConfirmDelete := []struct {
	//	StudentUUIDForArgs string
	//	ExpectError        error
	//} {
	//	{
	//		StudentUUIDForArgs: "student-111111111111",
	//		ExpectError:        nil,
	//	},
	//} -> CheckIfStudentAuthExists 개발 후 적용
}
