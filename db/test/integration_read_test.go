package test

import (
	"auth/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func Test_Accessor_GetStudentAuthWithID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
	}()

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
			log.Fatal(fmt.Sprintf("error occurs while creating parent auth, err: %v", err))
		}
	}

	for _, init := range []struct {
		UUID, StudentID, StudentPW, ParentUUID string
	} {
		{
			UUID:       "student-111111111111",
			StudentID:  "jinhong0719",
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
		StudentID                         string
		ExpectUUID, ExpectStudentID       string
		ExpectStudentPW, ExpectParentUUID string
		ExpectError                       error
	} {
		{ // success case
			StudentID:        "jinhong0719",
			ExpectUUID:       "student-111111111111",
			ExpectStudentID:  "jinhong0719",
			ExpectStudentPW:  passwords["testPW1"],
			ExpectParentUUID: "parent-111111111111",
			ExpectError:      nil,
		}, { // no exist student id
			StudentID:    "noExistStudentID",
			ExpectError:  gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		expectResult := &model.StudentAuth{
			UUID:       model.UUID(test.ExpectUUID),
			StudentID:  model.StudentID(test.ExpectStudentID),
			StudentPW:  model.StudentPW(test.ExpectStudentPW),
			ParentUUID: model.ParentUUID(test.ExpectParentUUID),
		}

		result, err := access.GetStudentAuthWithID(test.StudentID)

		assert.Equalf(t, test.ExpectError, err, "error assertion fail (test case: %v)", test)
		assert.Equalf(t, expectResult, result.ExceptGormModel(), "result model assertion fail (test case: %v)", test)
	}

	waitForFinish.Done()
}
