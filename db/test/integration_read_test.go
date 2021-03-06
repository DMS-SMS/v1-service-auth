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
		waitForFinish.Done()
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
}

func Test_Accessor_GetTeacherAuthWithID(t *testing.T) {
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
			TeacherID: "jinhong0719",
			TeacherPW: passwords["testPW1"],
		},
	} {
		_, err := access.CreateTeacherAuth(&model.TeacherAuth{
			UUID:     model.UUID(init.UUID),
			TeacherID: model.TeacherID(init.TeacherID),
			TeacherPW: model.TeacherPW(init.TeacherPW),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating teacher auth, err: %v", err))
		}
	}

	tests := []struct {
		TeacherID                                    string
		ExpectUUID, ExpectTeacherID, ExpectTeacherPW string
		ExpectError                                  error
	} {
		{ // success case
			TeacherID:       "jinhong0719",
			ExpectUUID:      "teacher-111111111111",
			ExpectTeacherID: "jinhong0719",
			ExpectTeacherPW: passwords["testPW1"],
			ExpectError:     nil,
		}, { // no exist student id
			TeacherID:   "noExistStudentID",
			ExpectError: gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		expectResult := &model.TeacherAuth{
			UUID:       model.UUID(test.ExpectUUID),
			TeacherID:  model.TeacherID(test.ExpectTeacherID),
			TeacherPW:  model.TeacherPW(test.ExpectTeacherPW),
		}

		result, err := access.GetTeacherAuthWithID(test.TeacherID)

		assert.Equalf(t, test.ExpectError, err, "error assertion fail (test case: %v)", test)
		assert.Equalf(t, expectResult, result.ExceptGormModel(), "result model assertion fail (test case: %v)", test)
	}
}

func Test_Accessor_GetParentAuthWithID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
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

	tests := []struct {
		ParentID                                   string
		ExpectUUID, ExpectParentID, ExpectParentPW string
		ExpectError                                error
	} {
		{ // success case
			ParentID:       "jinhong0719",
			ExpectUUID:     "parent-111111111111",
			ExpectParentID: "jinhong0719",
			ExpectParentPW: passwords["testPW1"],
			ExpectError:    nil,
		}, { // no exist student id
			ParentID:    "noExistStudentID",
			ExpectError: gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		expectResult := &model.ParentAuth{
			UUID:     model.UUID(test.ExpectUUID),
			ParentID: model.ParentID(test.ExpectParentID),
			ParentPW: model.ParentPW(test.ExpectParentPW),
		}

		result, err := access.GetParentAuthWithID(test.ParentID)

		assert.Equalf(t, test.ExpectError, err, "error assertion fail (test case: %v)", test)
		assert.Equalf(t, expectResult, result.ExceptGormModel(), "result model assertion fail (test case: %v)", test)
	}
}

func Test_Accessor_GetStudentAuthWithUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
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
		StudentUUID                         string
		ExpectUUID, ExpectStudentID       string
		ExpectStudentPW, ExpectParentUUID string
		ExpectError                       error
	} {
		{ // success case
			StudentUUID:      "student-111111111111",
			ExpectUUID:       "student-111111111111",
			ExpectStudentID:  "jinhong0719",
			ExpectStudentPW:  passwords["testPW1"],
			ExpectParentUUID: "parent-111111111111",
			ExpectError:      nil,
		}, { // no exist student id
			StudentUUID: "noExistStudentUUID",
			ExpectError: gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		expectResult := &model.StudentAuth{
			UUID:       model.UUID(test.ExpectUUID),
			StudentID:  model.StudentID(test.ExpectStudentID),
			StudentPW:  model.StudentPW(test.ExpectStudentPW),
			ParentUUID: model.ParentUUID(test.ExpectParentUUID),
		}

		result, err := access.GetStudentAuthWithUUID(test.StudentUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion fail (test case: %v)", test)
		assert.Equalf(t, expectResult, result.ExceptGormModel(), "result model assertion fail (test case: %v)", test)
	}
}

func Test_Accessor_GetTeacherAuthWithUUID(t *testing.T) {
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
			TeacherID: "jinhong0719",
			TeacherPW: passwords["testPW1"],
		},
	} {
		_, err := access.CreateTeacherAuth(&model.TeacherAuth{
			UUID:      model.UUID(init.UUID),
			TeacherID: model.TeacherID(init.TeacherID),
			TeacherPW: model.TeacherPW(init.TeacherPW),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating teacher auth, err: %v", err))
		}
	}

	tests := []struct {
		TeacherUUID, ExpectUUID          string
		ExpectTeacherID, ExpectTeacherPW string
		ExpectError                      error
	} {
		{ // success case
			TeacherUUID:     "teacher-111111111111",
			ExpectUUID:      "teacher-111111111111",
			ExpectTeacherID: "jinhong0719",
			ExpectTeacherPW: passwords["testPW1"],
			ExpectError:     nil,
		}, { // no exist student id
			TeacherUUID: "noExistTeacherUUID",
			ExpectError: gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		expectResult := &model.TeacherAuth{
			UUID:       model.UUID(test.ExpectUUID),
			TeacherID:  model.TeacherID(test.ExpectTeacherID),
			TeacherPW:  model.TeacherPW(test.ExpectTeacherPW),
		}

		result, err := access.GetTeacherAuthWithUUID(test.TeacherUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion fail (test case: %v)", test)
		assert.Equalf(t, expectResult, result.ExceptGormModel(), "result model assertion fail (test case: %v)", test)
	}
}

func Test_Accessor_GetParentAuthWithUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
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

	tests := []struct {
		ParentUUID, ExpectUUID           string
		ExpectParentID, ExpectParentPW string
		ExpectError                    error
	} {
		{ // success case
			ParentUUID:     "parent-111111111111",
			ExpectUUID:     "parent-111111111111",
			ExpectParentID: "jinhong0719",
			ExpectParentPW: passwords["testPW1"],
			ExpectError:    nil,
		}, { // no exist student id
			ParentUUID:  "noExistParentUUID",
			ExpectError: gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		expectResult := &model.ParentAuth{
			UUID:     model.UUID(test.ExpectUUID),
			ParentID: model.ParentID(test.ExpectParentID),
			ParentPW: model.ParentPW(test.ExpectParentPW),
		}

		result, err := access.GetParentAuthWithUUID(test.ParentUUID)

		assert.Equalf(t, test.ExpectError, err, "error assertion fail (test case: %v)", test)
		assert.Equalf(t, expectResult, result.ExceptGormModel(), "result model assertion fail (test case: %v)", test)
	}
}

func Test_Accessor_GetStudentUUIDsWithInform(t *testing.T) {
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
		StudentUUID                   string
		Grade, Class, StudentNumber   int64
		Name, PhoneNumber, ProfileURI string
		ExpectUUIDArr                 []string
		ExpectError                   error
	} {
		{
			StudentUUID:   "student-111111111111",
			ExpectUUIDArr: []string{"student-111111111111"},
			ExpectError:   nil,
		}, {
			Grade:         2,
			ExpectUUIDArr: []string{"student-111111111111", "student-222222222222", "student-333333333333"},
			ExpectError:   nil,
		}, {
			PhoneNumber:   "01088378347",
			ExpectUUIDArr: ([]string)(nil),
			ExpectError:   gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		uuidArr, err := access.GetStudentUUIDsWithInform(&model.StudentInform{
			StudentUUID:   model.StudentUUID(test.StudentUUID),
			Grade:         model.Grade(test.Grade),
			Class:         model.Class(test.Class),
			StudentNumber: model.StudentNumber(test.StudentNumber),
			Name:          model.Name(test.Name),
			PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
			ProfileURI:    model.ProfileURI(test.ProfileURI),
		})

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectUUIDArr, uuidArr, "uuid array result assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetTeacherUUIDsWithInform(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = access.Rollback()
		waitForFinish.Done()
	}()

	// 선생님 계정 생성
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
		}, {
			UUID:      "teacher-333333333333",
			TeacherID: "jinhong07193",
			TeacherPW: passwords["testPW1"],
		},
	} {
		_, err := access.CreateTeacherAuth(&model.TeacherAuth{
			UUID:      model.UUID(init.UUID),
			TeacherID: model.TeacherID(init.TeacherID),
			TeacherPW: model.TeacherPW(init.TeacherPW),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating teacher auth, err: %v", err))
		}
	}

	// 학생 정보 생성
	for _, init := range []struct {
		TeacherUUID, Name, PhoneNumber string
		Grade, Class                   int64
	} {
		{
			TeacherUUID: "teacher-111111111111",
			Grade:       2,
			Class:       2,
			Name:        "박진홍",
			PhoneNumber: "01011111111",
		}, {
			TeacherUUID: "teacher-222222222222",
			Grade:       2,
			Class:       2,
			Name:        "윤석준",
			PhoneNumber: "01022222222",
		}, {
			TeacherUUID: "teacher-333333333333",
			Name:        "오준상",
			PhoneNumber: "01033333333",
		},
	} {
		_, err := access.CreateTeacherInform(&model.TeacherInform{
			TeacherUUID: model.TeacherUUID(init.TeacherUUID),
			Grade:       model.Grade(init.Grade),
			Class:       model.Class(init.Class),
			Name:        model.Name(init.Name),
			PhoneNumber: model.PhoneNumber(init.PhoneNumber),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating student inform, err: %v", err))
		}
	}

	tests := []struct {
		TeacherUUID, Name, PhoneNumber string
		Grade, Class                   int64
		ExpectUUIDArr                  []string
		ExpectError                    error
	} {
		{
			TeacherUUID:   "teacher-111111111111",
			ExpectUUIDArr: []string{"teacher-111111111111"},
			ExpectError:   nil,
		}, {
			Grade:         2,
			Class:         2,
			ExpectUUIDArr: []string{"teacher-111111111111", "teacher-222222222222"},
			ExpectError:   nil,
		}, {
			Grade:         2,
			Class:         2,
			Name:          "박진홍",
			ExpectUUIDArr: []string{"teacher-111111111111"},
			ExpectError:   nil,
		}, {
			Grade:         model.TeacherInformInstance.Grade.NullReplaceValue(),
			Class:         model.TeacherInformInstance.Class.NullReplaceValue(),
			ExpectUUIDArr: []string{"teacher-333333333333"},
			ExpectError:   nil,
		}, {
			PhoneNumber:   "01088378347",
			ExpectUUIDArr: ([]string)(nil),
			ExpectError:   gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		uuidArr, err := access.GetTeacherUUIDsWithInform(&model.TeacherInform{
			TeacherUUID: model.TeacherUUID(test.TeacherUUID),
			Grade:       model.Grade(test.Grade),
			Class:       model.Class(test.Class),
			Name:        model.Name(test.Name),
			PhoneNumber: model.PhoneNumber(test.PhoneNumber),
		})

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectUUIDArr, uuidArr, "uuid array result assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetParentUUIDsWithInform(t *testing.T) {
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

	// 학생 정보 생성
	for _, init := range []struct {
		ParentUUID, Name, PhoneNumber string
	} {
		{
			ParentUUID: "parent-111111111111",
			Name: "박진홍",
			PhoneNumber: "01011111111",
		}, {
			ParentUUID: "parent-222222222222",
			Name: "박진홍",
			PhoneNumber: "01022222222",
		}, {
			ParentUUID: "parent-333333333333",
			Name: "박진홍",
			PhoneNumber: "01033333333",
		},
	} {
		_, err := access.CreateParentInform(&model.ParentInform{
			ParentUUID:  model.ParentUUID(init.ParentUUID),
			Name:        model.Name(init.Name),
			PhoneNumber: model.PhoneNumber(init.PhoneNumber),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating student inform, err: %v", err))
		}
	}

	tests := []struct {
		ParentUUID, Name, PhoneNumber string
		ExpectUUIDArr                 []string
		ExpectError                   error
	} {
		{
			ParentUUID:   "parent-111111111111",
			ExpectUUIDArr: []string{"parent-111111111111"},
			ExpectError:   nil,
		}, {
			Name:          "박진홍",
			ExpectUUIDArr: []string{"parent-111111111111", "parent-222222222222", "parent-333333333333"},
			ExpectError:   nil,
		}, {
			ParentUUID:    "parent-333333333333",
			Name:          "박진홍",
			PhoneNumber:   "01033333333",
			ExpectUUIDArr: []string{"parent-333333333333"},
			ExpectError:   nil,
		}, {
			PhoneNumber:   "01088378347",
			ExpectUUIDArr: ([]string)(nil),
			ExpectError:   gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		uuidArr, err := access.GetParentUUIDsWithInform(&model.ParentInform{
			ParentUUID:  model.ParentUUID(test.ParentUUID),
			Name:        model.Name(test.Name),
			PhoneNumber: model.PhoneNumber(test.PhoneNumber),
		})

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectUUIDArr, uuidArr, "uuid array result assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetStudentInformWithUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
	}()

	// 학부모 계정 생성
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
			PhoneNumber: "01088378347",
			ProfileURI: "example.com/profiles/student-111111111111",
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
		{
			StudentUUIDForArgs: "student-111111111111",
			StudentUUID:        "student-111111111111",
			Grade:              2,
			Class:              2,
			StudentNumber:      7,
			Name:               "박진홍",
			PhoneNumber:        "01088378347",
			ProfileURI:         "example.com/profiles/student-111111111111",
			ExpectError:        nil,
		},
		{
			StudentUUIDForArgs: "student-222222222222",
			ExpectError:        gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		expectResult := &model.StudentInform{
			StudentUUID:   model.StudentUUID(test.StudentUUID),
			Grade:         model.Grade(test.Grade),
			Class:         model.Class(test.Class),
			StudentNumber: model.StudentNumber(test.StudentNumber),
			Name:          model.Name(test.Name),
			PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
			ProfileURI:    model.ProfileURI(test.ProfileURI),
		}
		result, err := access.GetStudentInformWithUUID(test.StudentUUIDForArgs)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, expectResult, result.ExceptGormModel(), "result inform model assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetStudentInformsWithUUIDs(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
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
			ParentPW: passwords["testPW1"],
		}, {
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
			StudentPW:  passwords["testPW1"],
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
			StudentUUID:   "student-111111111111",
			Grade:         2,
			Class:         2,
			StudentNumber: 7,
			Name:          "박진홍",
			PhoneNumber:   "01011111111",
			ProfileURI:    "example.com/profiles/student-111111111111",
		}, {
			StudentUUID:   "student-222222222222",
			Grade:         1,
			Class:         2,
			StudentNumber: 8,
			Name:          "진홍박",
			PhoneNumber:   "01022222222",
			ProfileURI:    "example.com/profiles/student-222222222222",
		}, {
			StudentUUID:   "student-333333333333",
			Grade:         3,
			Class:         2,
			StudentNumber: 8,
			Name:          "박진홍",
			PhoneNumber:   "01033333333",
			ProfileURI:    "example.com/profiles/student-333333333333",
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

	_ = access.DeleteStudentAuth("student-333333333333")
	_ = access.DeleteStudentInform("student-333333333333")

	tests := []struct {
		StudentUUIDs         []string
		ExpectInforms []*model.StudentInform
		ExpectError          error
	} {
		{
			StudentUUIDs: []string{"student-111111111111", "student-222222222222"},
			ExpectInforms: []*model.StudentInform{
				{
					StudentUUID:   "student-111111111111",
					Grade:         2,
					Class:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01011111111",
					ProfileURI:    "example.com/profiles/student-111111111111",
				}, {
					StudentUUID:   "student-222222222222",
					Grade:         1,
					Class:         2,
					StudentNumber: 8,
					Name:          "진홍박",
					PhoneNumber:   "01022222222",
					ProfileURI:    "example.com/profiles/student-222222222222",
				},
			},
			ExpectError: nil,
		}, {
			StudentUUIDs: []string{"student-222222222222", "student-333333333333", "student-111111111111", "student-222222222222"},
			ExpectInforms: []*model.StudentInform{
				{
					StudentUUID:   "student-222222222222",
					Grade:         1,
					Class:         2,
					StudentNumber: 8,
					Name:          "진홍박",
					PhoneNumber:   "01022222222",
					ProfileURI:    "example.com/profiles/student-222222222222",
				}, { }, {
					StudentUUID:   "student-111111111111",
					Grade:         2,
					Class:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01011111111",
					ProfileURI:    "example.com/profiles/student-111111111111",
				}, {
					StudentUUID:   "student-222222222222",
					Grade:         1,
					Class:         2,
					StudentNumber: 8,
					Name:          "진홍박",
					PhoneNumber:   "01022222222",
					ProfileURI:    "example.com/profiles/student-222222222222",
				},
			},
			ExpectError: gorm.ErrRecordNotFound,
		}, {
			StudentUUIDs:  []string{"student-444444444444"},
			ExpectInforms: []*model.StudentInform{{}},
			ExpectError:   gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		resultInforms, err := access.GetStudentInformsWithUUIDs(test.StudentUUIDs)
		var exceptedInforms []*model.StudentInform
		for _, inform := range resultInforms {
			exceptedInforms = append(exceptedInforms, inform.ExceptGormModel())
		}

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, test.ExpectInforms, exceptedInforms, "result informs assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetTeacherInformWithUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
	}()

	// 선생님 계정 생성
	for _, init := range []struct {
		UUID, TeacherID, TeacherPW string
	} {
		{
			UUID:      "teacher-111111111111",
			TeacherID: "jinhong0719",
			TeacherPW: passwords["testPW1"],
		},
	} {
		_, err := access.CreateTeacherAuth(&model.TeacherAuth{
			UUID:      model.UUID(init.UUID),
			TeacherID: model.TeacherID(init.TeacherID),
			TeacherPW: model.TeacherPW(init.TeacherPW),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating teacher auth, err: %v", err))
		}
	}

	// 선생님 정보 생성
	for _, init := range []struct {
		TeacherUUID       string
		Name, PhoneNumber string
		Grade, Class      int64
	} {
		{
			TeacherUUID: "teacher-111111111111",
			Grade:       2,
			Class:       2,
			Name:        "박진홍",
			PhoneNumber: "01088378347",
		},
	} {
		_, err := access.CreateTeacherInform(&model.TeacherInform{
			TeacherUUID: model.TeacherUUID(init.TeacherUUID),
			Grade:       model.Grade(init.Grade),
			Class:       model.Class(init.Class),
			Name:        model.Name(init.Name),
			PhoneNumber: model.PhoneNumber(init.PhoneNumber),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating teacher inform, err: %v", err))
		}
	}

	tests := []struct {
		TeacherUUIDForArgs, TeacherUUID string
		Grade, Class                    int64
		Name, PhoneNumber               string
		ExpectError                     error
	} {
		{
			TeacherUUIDForArgs: "teacher-111111111111",
			TeacherUUID:        "teacher-111111111111",
			Grade:              2,
			Class:              2,
			Name:               "박진홍",
			PhoneNumber:        "01088378347",
			ExpectError:        nil,
		},
		{
			TeacherUUIDForArgs: "student-222222222222",
			ExpectError:        gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		expectResult := &model.TeacherInform{
			TeacherUUID:   model.TeacherUUID(test.TeacherUUID),
			Grade:         model.Grade(test.Grade),
			Class:         model.Class(test.Class),
			Name:          model.Name(test.Name),
			PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
		}
		result, err := access.GetTeacherInformWithUUID(test.TeacherUUIDForArgs)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, expectResult, result.ExceptGormModel(), "result inform model assertion error (test case: %v)", test)
	}
}

func Test_Accessor_GetParentInformWithUUID(t *testing.T) {
	access, err := manager.BeginTx()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		access.Rollback()
		waitForFinish.Done()
	}()

	// 학부모 계정 생성
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

	// 학부모 정보 생성
	for _, init := range []struct {
		ParentUUID, Name, PhoneNumber string
	} {
		{
			ParentUUID:  "parent-111111111111",
			Name:        "박진홍",
			PhoneNumber: "01088378347",
		},
	} {
		_, err := access.CreateParentInform(&model.ParentInform{
			ParentUUID:  model.ParentUUID(init.ParentUUID),
			Name:        model.Name(init.Name),
			PhoneNumber: model.PhoneNumber(init.PhoneNumber),
		})
		if err != nil {
			log.Fatal(fmt.Sprintf("error occurs while creating parent inform, err: %v", err))
		}
	}

	tests := []struct {
		ParentUUIDForArgs, ParentUUID string
		Name, PhoneNumber             string
		ExpectError                   error
	} {
		{
			ParentUUIDForArgs: "parent-111111111111",
			ParentUUID:        "parent-111111111111",
			Name:              "박진홍",
			PhoneNumber:       "01088378347",
			ExpectError:       nil,
		},
		{
			ParentUUIDForArgs: "student-222222222222",
			ExpectError:       gorm.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		expectResult := &model.ParentInform{
			ParentUUID:  model.ParentUUID(test.ParentUUID),
			Name:        model.Name(test.Name),
			PhoneNumber: model.PhoneNumber(test.PhoneNumber),
		}
		result, err := access.GetParentInformWithUUID(test.ParentUUIDForArgs)

		assert.Equalf(t, test.ExpectError, err, "error assertion error (test case: %v)", test)
		assert.Equalf(t, expectResult, result.ExceptGormModel(), "result inform model assertion error (test case: %v)", test)
	}
}
