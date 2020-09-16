package handler

import (
	test "auth/handler/for_test"
	"auth/model"
	proto "auth/proto/golang/auth"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
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

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(mockForDB)

		var req = new(proto.LoginStudentAuthRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		var resp = new(proto.LoginStudentAuthResponse)
		_ = defaultHandler.LoginStudentAuth(ctx, req, resp)

		assert.Equalf(t, testCase.ExpectedStatus, resp.Status, "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedLoggedInStudentUUID, resp.LoggedInStudentUUID, "student uuid assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	mockForDB.AssertExpectations(t)
}
