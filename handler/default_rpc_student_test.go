package handler

import (
	test "auth/handler/for_test"
	"auth/model"
	"github.com/jinzhu/gorm"
	"net/http"
	"testing"
)

func Test_default_LoginStudentAuth(t *testing.T) {
	tests := []test.LoginStudentAuthCase{
		{ // success case
			StudentID: "jinhong0719",
			StudentPW: "testPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentAuthWithID": {&model.StudentAuth{
					UUID:       "student-111111111111", // 중복 X !!
					StudentID:  "jinhong0719",
					StudentPW:  "testPW",
					ParentUUID: "parent-111111111111",
				}, nil},
				"Commit": {},
			},
			ExpectedStatus:              http.StatusOK,
			ExpectedLoggedInStudentUUID: "student-111111111111",
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			XRequestID:      test.EmptyReplaceValueForString,
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // no exist Span-Context -> Proxy Authorization Required
			SpanContextString: test.EmptyReplaceValueForString,
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // Student ID no exists
			StudentID: "jinhong0719",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":              {},
				"GetStudentAuthWithID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"Commit":               {},
			},
			ExpectedStatus: CodeStudentIDNoExist,
		}, { // incorrect Student PW
			StudentID: "jinhong0719",
			StudentPW: "incorrectPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentAuthWithID": {&model.StudentAuth{
					UUID:       "student-111111111111", // 중복 X !!
					StudentID:  "jinhong0719",
					StudentPW:  "testPW",
					ParentUUID: "parent-111111111111",
				}, nil},
				"Commit": {},
			},
			ExpectedStatus: CodeIncorrectStudentPW,
		},
	}

	for _, _ = range tests {
		
	}
}