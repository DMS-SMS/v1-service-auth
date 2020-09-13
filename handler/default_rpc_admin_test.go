package handler

import (
	"auth/model"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	"net/http"
	"testing"

	//proto "auth/proto/golang/auth"
)

func Test_default_CreateNewStudent(t *testing.T) {
	const studentUUIDRegexString = "^student-\\d{12}"

	tests := []createNewStudentTest{
		{ // success case
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, nil},
				"Commit":                   {&gorm.DB{}},
			},
			ExpectedStatus:      http.StatusCreated,
			ExpectedStudentUUID: studentUUIDRegexString,
		},
		{ // not admin uuid -> forbidden
			UUID:           "NotAdminAuthUUID", // (admin-숫자 12개의 형식이여야 함)
			ExpectedMethods: map[method]returns{},
			ExpectedStatus: http.StatusForbidden,
		}, { // invalid request value -> Proxy Authorization Required
			StudentID:      "유효하지 않은 아이디", // ASCII, 4~16 사이 문자열이여야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, {
			Grade:          100, // 1~3 사이의 숫자여야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, {
			Name:           "Invalid Name", // 2~4 글자의 한글이어야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		},
	}

	for _, _ = range tests {

	}
}