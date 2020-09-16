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

func Test_default_LoginStudentAuth(t *testing.T) {
	hashedByte, _ := bcrypt.GenerateFromPassword([]byte("testPW"), 1)

	tests := []test.LoginStudentAuthCase{
		{ // success case
			StudentID: "jinhong07191",
			StudentPW: "testPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentAuthWithID": {&model.StudentAuth{
					UUID:       "student-111111111111", // 중복 X !!
					StudentID:  "jinhong0719",
					StudentPW:  model.StudentPW(string(hashedByte)),
					ParentUUID: "parent-111111111111",
				}, nil},
				"Commit": {&gorm.DB{}},
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
			StudentID: "jinhong07192",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":              {},
				"GetStudentAuthWithID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"Rollback":             {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeStudentIDNoExist,
		}, { // GetStudentAuthWithID unexpected error
			StudentID: "jinhong07193",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":              {},
				"GetStudentAuthWithID": {&model.StudentAuth{}, errors.New("unexpected error")},
				"Rollback":             {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // incorrect Student PW
			StudentID: "jinhong07194",
			StudentPW: "incorrectPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentAuthWithID": {&model.StudentAuth{
					UUID:       "student-111111111111", // 중복 X !!
					StudentID:  "jinhong07194",
					StudentPW:  model.StudentPW(string(hashedByte)),
					ParentUUID: "parent-111111111111",
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeIncorrectStudentPWForLogin,
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

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedLoggedInStudentUUID, resp.LoggedInStudentUUID, "student uuid assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	mockForDB.AssertExpectations(t)
}

func Test_default_ChangeStudentPW(t *testing.T) {
	hashedTestPW1, _ := bcrypt.GenerateFromPassword([]byte("testPW1"), 1)
	hashedTestPW2, _ := bcrypt.GenerateFromPassword([]byte("testPW2"), 1)


	tests := []test.ChangeStudentPWCase{
		{ // success case
			UUID:        "student-111111111111",
			StudentUUID: "student-111111111111",
			CurrentPW:   "testPW1",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{
					UUID:      "student-111111111111",
					StudentPW: model.StudentPW(string(hashedTestPW1)),
				}, nil},
				"ChangeStudentPW": {nil},
				"Commit":          {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
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
		}, { // forbidden (not student)
			UUID:           "parent-111111111112",
			StudentUUID:    "parent-111111111112",
			CurrentPW:      "testPW1",
			RevisionPW:     "NewPassword",
			ExpectedStatus: http.StatusForbidden,
		}, { // forbidden (not my auth)
			UUID:           "student-111111111113",
			StudentUUID:    "student-111111111114",
			CurrentPW:      "testPW1",
			RevisionPW:     "NewPassword",
			ExpectedStatus: http.StatusForbidden,
		}, { // not exists student
			UUID:           "student-111111111115",
			StudentUUID:    "student-111111111115",
			CurrentPW:      "testPW1",
			RevisionPW:     "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
		}, { // 현재 Password 불일치
			UUID:        "student-111111111116",
			StudentUUID: "student-111111111116",
			CurrentPW:   "testPW1",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{
					UUID:      "student-111111111116",
					StudentPW: model.StudentPW(string(hashedTestPW2)),
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeIncorrectStudentPWForChange,
		}, { // GetStudentAuthWithUUID 에러 반환
			UUID:        "student-111111111117",
			StudentUUID: "student-111111111117",
			CurrentPW:   "testPW1",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, errors.New("DB not connected")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // ChangeStudentPW 에러 반환
			UUID:        "student-111111111118",
			StudentUUID: "student-111111111118",
			CurrentPW:   "testPW1",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{
					UUID:      "student-111111111118",
					StudentPW: model.StudentPW(string(hashedTestPW1)),
				}, nil},
				"ChangeStudentPW": {errors.New("DB not connected")},
				"Rollback":        {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetStudentAuthWithUUID Short Hashed PW 반환
			UUID:        "student-111111111119",
			StudentUUID: "student-111111111119",
			CurrentPW:   "testPW1",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{
					UUID:      "student-111111111119",
					StudentPW: "TooShortHashedPasword",
				}, nil},
				"Rollback":        {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(mockForDB)

		req := new(proto.ChangeStudentPWRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.ChangeStudentPWResponse)
		_ = defaultHandler.ChangeStudentPW(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}
}