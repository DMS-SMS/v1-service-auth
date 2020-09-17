package handler

import (
	test "auth/handler/for_test"
	"auth/model"
	proto "auth/proto/golang/auth"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"testing"
)

func Test_default_LoginTeacherAuth(t *testing.T) {
	hashedByte, _ := bcrypt.GenerateFromPassword([]byte("testPW"), 1)

	tests := []test.LoginTeacherAuthCase{
		{ // success case
			TeacherID: "jinhong07191",
			TeacherPW: "testPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetTeacherAuthWithID": {&model.TeacherAuth{
					UUID:      "teacher-111111111111",
					TeacherID: "jinhong07191",
					TeacherPW: model.TeacherPW(string(hashedByte)),
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus:              http.StatusOK,
			ExpectedLoggedInTeacherUUID: "teacher-111111111111",
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
			TeacherID: "jinhong07192",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":              {},
				"GetTeacherAuthWithID": {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"Rollback":             {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeTeacherIDNoExist,
		}, { // GetStudentAuthWithID unexpected error
			TeacherID: "jinhong07193",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":              {},
				"GetTeacherAuthWithID": {&model.TeacherAuth{}, errors.New("unexpected error")},
				"Rollback":             {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // incorrect Student PW
			TeacherID: "jinhong07194",
			TeacherPW: "incorrectPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetTeacherAuthWithID": {&model.TeacherAuth{
					UUID:      "teacher-111111111111", // 중복 X !!
					TeacherID: "jinhong07194",
					TeacherPW: model.TeacherPW(string(hashedByte)),
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeIncorrectTeacherPWForLogin,
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(mockForDB)

		var req = new(proto.LoginTeacherAuthRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		var resp = new(proto.LoginTeacherAuthResponse)
		_ = defaultHandler.LoginTeacherAuth(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedLoggedInTeacherUUID, resp.LoggedInTeacherUUID, "logged in uuid assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	mockForDB.AssertExpectations(t)
}
