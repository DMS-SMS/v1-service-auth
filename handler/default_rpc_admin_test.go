package handler

import (
	"auth/model"
	"auth/tool/mysqlerr"
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
		}, { // not admin uuid -> forbidden
			UUID:            "NotAdminAuthUUID", // (admin-숫자 12개의 형식이여야 함)
			ExpectedMethods: map[method]returns{},
			ExpectedStatus:  http.StatusForbidden,
		}, { // invalid request value -> Proxy Authorization Required
			StudentID: "유효하지 않은 아이디", // ASCII, 4~16 사이 문자열이여야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Grade: 100, // 1~3 사이의 숫자여야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Name: "Invalid Name", // 2~4 글자의 한글이어야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			XRequestID:      emptyReplaceValueForString,
			ExpectedMethods: map[method]returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[method]returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // no exist Span-Context -> Proxy Authorization Required
			XRequestID: emptyReplaceValueForString,
			ExpectedMethods: map[method]returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			XRequestID:      "InvalidSpanContext",
			ExpectedMethods: map[method]returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // student id duplicate -> Conflict -101
			StudentID: "jinhong0719",
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, mysqlerr.DuplicateEntry(model.StudentAuthInstance.StudentID.KeyName(), "jinhong0719")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeStudentIDDuplicate,
		}, { // parent uuid fk constraint fail -> Conflict -102
			ParentUUID: "parent-111111111111",
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, studentAuthParentUUIDFKConstraintFailError},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeParentUUIDNoExist,
		}, { // student number duplicate -> Conflict -103
			Grade:         2,
			Class:         2,
			StudentNumber: 7,
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, mysqlerr.DuplicateEntry(model.StudentInformInstance.StudentNumber.KeyName(), "2207")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeStudentNumberDuplicate,
		}, { // phone number duplicate -> Conflict -104
			PhoneNumber: "01088378347",
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, mysqlerr.DuplicateEntry(model.StudentInformInstance.PhoneNumber.KeyName(), "01088378347")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodePhoneNumberDuplicate,
		},
	}

	for _, _ = range tests {

	}
}