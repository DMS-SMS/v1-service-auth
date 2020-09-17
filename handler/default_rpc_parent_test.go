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

func Test_default_LoginParentAuth(t *testing.T) {
	newMock, defaultHandler := generateVarForTest()
	hashedByte, _ := bcrypt.GenerateFromPassword([]byte("testPW"), 1)

	tests := []test.LoginParentAuthCase{
		{ // success case
			ParentID: "jinhong07191",
			ParentPW: "testPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetParentAuthWithID": {&model.ParentAuth{
					UUID:     "parent-111111111111",
					ParentID: "jinhong07191",
					ParentPW: model.ParentPW(string(hashedByte)),
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus:             http.StatusOK,
			ExpectedLoggedInParentUUID: "parent-111111111111",
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
		}, { // Parent ID no exists
			ParentID: "jinhong07192",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":             {},
				"GetParentAuthWithID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeParentIDNoExist,
		}, { // GetStudentAuthWithID unexpected error
			ParentID: "jinhong07193",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":             {},
				"GetParentAuthWithID": {&model.ParentAuth{}, errors.New("unexpected error")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // incorrect Student PW
			ParentID: "jinhong07194",
			ParentPW: "incorrectPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetParentAuthWithID": {&model.ParentAuth{
					UUID:     "teacher-111111111111", // 중복 X !!
					ParentID: "jinhong07194",
					ParentPW: model.ParentPW(string(hashedByte)),
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeIncorrectParentPWForLogin,
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		var req = new(proto.LoginParentAuthRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		var resp = new(proto.LoginParentAuthResponse)
		_ = defaultHandler.LoginParentAuth(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedLoggedInParentUUID, resp.LoggedInParentUUID, "logged in uuid assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	newMock.AssertExpectations(t)
}
