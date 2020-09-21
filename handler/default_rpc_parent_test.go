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
	"time"
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
		}, { // GetParentAuthWithID unexpected error
			ParentID: "jinhong07193",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":             {},
				"GetParentAuthWithID": {&model.ParentAuth{}, errors.New("unexpected error")},
				"Rollback":            {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // incorrect Parent PW
			ParentID: "jinhong07194",
			ParentPW: "incorrectPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetParentAuthWithID": {&model.ParentAuth{
					UUID:     "parent-111111111111",
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

func Test_default_ChangeParentPW(t *testing.T) {
	newMock, defaultHandler := generateVarForTest()
	hashedTestPW, _ := bcrypt.GenerateFromPassword([]byte("testPW"), 1)

	tests := []test.ChangeParentPWCase{
		{ // success case
			UUID:       "parent-111111111111",
			ParentUUID: "parent-111111111111",
			CurrentPW:  "testPW",
			RevisionPW: "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetParentAuthWithUUID": {&model.ParentAuth{
					UUID:     "parent-111111111111",
					ParentPW: model.ParentPW(string(hashedTestPW)),
				}, nil},
				"ChangeParentPW": {nil},
				"Commit":         {&gorm.DB{}},
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
		}, { // forbidden (not parent)
			UUID:           "student-111111111112",
			ParentUUID:     "student-111111111112",
			CurrentPW:      "testPW",
			RevisionPW:     "NewPassword",
			ExpectedStatus: http.StatusForbidden,
		}, { // forbidden (not my auth)
			UUID:           "parent-111111111113",
			ParentUUID:     "parent-111111111114",
			CurrentPW:      "testPW",
			RevisionPW:     "NewPassword",
			ExpectedStatus: http.StatusForbidden,
		}, { // not exists parent
			UUID:       "parent-111111111115",
			ParentUUID: "parent-111111111115",
			CurrentPW:  "testPW",
			RevisionPW: "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
		}, { // 현재 Password 불일치
			UUID:       "parent-111111111116",
			ParentUUID: "parent-111111111116",
			CurrentPW:  "IncorrectPassword",
			RevisionPW: "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetParentAuthWithUUID": {&model.ParentAuth{
					UUID:     "parent-111111111116",
					ParentPW: model.ParentPW(string(hashedTestPW)),
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeIncorrectParentPWForChange,
		}, { // GetParentAuthWithUUID 에러 반환
			UUID:       "parent-111111111117",
			ParentUUID: "parent-111111111117",
			CurrentPW:  "testPW",
			RevisionPW: "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, errors.New("DB not connected")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // ChangeParentPW 에러 반환
			UUID:       "parent-111111111118",
			ParentUUID: "parent-111111111118",
			CurrentPW:  "testPW",
			RevisionPW: "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetParentAuthWithUUID": {&model.ParentAuth{
					UUID:     "parent-111111111118",
					ParentPW: model.ParentPW(string(hashedTestPW)),
				}, nil},
				"ChangeParentPW": {errors.New("DB not connected")},
				"Rollback":       {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetParentAuthWithUUID Short Hashed PW 반환
			UUID:       "parent-111111111119",
			ParentUUID: "parent-111111111119",
			CurrentPW:  "testPW",
			RevisionPW: "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetParentAuthWithUUID": {&model.ParentAuth{
					UUID:     "parent-111111111119",
					ParentPW: "TooShortHashedPasword",
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.ChangeParentPWRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.ChangeParentPWResponse)
		_ = defaultHandler.ChangeParentPW(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	newMock.AssertExpectations(t)
}

func Test_default_GetParentInformWithUUID(t *testing.T) {
	newMock, defaultHandler := generateVarForTest()
	now := time.Now()

	tests := []test.GetParentInformWithUUIDCase{
		{ // success case
			UUID:       "parent-111111111111",
			ParentUUID: "parent-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetParentInformWithUUID": {&model.ParentInform{
					Model:       gorm.Model{CreatedAt: now, UpdatedAt: now},
					ParentUUID:  "parent-111111111111",
					Name:        "박진홍",
					PhoneNumber: "01088378347",
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedInform: &model.ParentInform{
				Name:        "박진홍",
				PhoneNumber: "01088378347",
			},
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			XRequestID:      test.EmptyReplaceValueForString,
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
			ExpectedInform:  &model.ParentInform{},
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
			ExpectedInform:  &model.ParentInform{},
		}, { // no exist Span-Context -> Proxy Authorization Required
			SpanContextString: test.EmptyReplaceValueForString,
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
			ExpectedInform:    &model.ParentInform{},
		}, { // invalid Span-Context -> Proxy Authorization Required
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
			ExpectedInform:    &model.ParentInform{},
		}, { // forbidden (not parent)
			UUID:           "student-111111111112",
			ParentUUID:     "student-111111111112",
			ExpectedStatus: http.StatusForbidden,
			ExpectedInform: &model.ParentInform{},
		}, { // forbidden (not my auth)
			UUID:           "parent-111111111113",
			ParentUUID:     "parent-111111111114",
			ExpectedStatus: http.StatusForbidden,
			ExpectedInform: &model.ParentInform{},
		}, { // no exist parent uuid
			UUID:       "parent-111111111115",
			ParentUUID: "parent-111111111115",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                 {},
				"GetParentInformWithUUID": {&model.ParentInform{}, gorm.ErrRecordNotFound},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedInform: &model.ParentInform{},
		}, { // GetParentInformWithUUID error return
			UUID:       "parent-111111111116",
			ParentUUID: "parent-111111111116",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                 {},
				"GetParentInformWithUUID": {&model.ParentInform{}, errors.New("DB not connected")},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedInform: &model.ParentInform{},
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.GetParentInformWithUUIDRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.GetParentInformWithUUIDResponse)
		_ = defaultHandler.GetParentInformWithUUID(ctx, req, resp)

		resultInform := &model.ParentInform{
			Name:          model.Name(resp.Name),
			PhoneNumber:   model.PhoneNumber(resp.PhoneNumber),
		}

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedInform, resultInform, "result inform assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	newMock.AssertExpectations(t)
}

func Test_default_GetParentUUIDsWithInform(t *testing.T) {
	newMock, defaultHandler := generateVarForTest()

	tests := []test.GetParentUUIDsWithInformCase{
		{ // success case (for admin auth)
			UUID: "admin-111111111111",
			Name: "이성진",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetParentUUIDsWithInform": {[]string{"parent-123412341234", "parent-432143214321"}, nil},
				"Commit":                   {&gorm.DB{}},
			},
			ExpectedStatus:      http.StatusOK,
			ExpectedParentUUIDs: []string{"parent-123412341234", "parent-432143214321"},
		}, { // success case (for parent auth)
			UUID:        "parent-111111111111",
			PhoneNumber: "01088378347",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetParentUUIDsWithInform": {[]string{"parent-111111111111"}, nil},
				"Commit":                   {&gorm.DB{}},
			},
			ExpectedStatus:      http.StatusOK,
			ExpectedParentUUIDs: []string{"parent-111111111111"},
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			UUID:            "parent-111111111111",
			XRequestID:      test.EmptyReplaceValueForString,
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			UUID:            "parent-111111111111",
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // no exist Span-Context -> Proxy Authorization Required
			UUID:              "parent-111111111111",
			SpanContextString: test.EmptyReplaceValueForString,
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			UUID:              "parent-111111111111",
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // forbidden (not parent)
			UUID:           "student-111111111111",
			ExpectedStatus: http.StatusForbidden,
		}, { // no exist parent uuid with that inform
			UUID:        "parent-111111111111",
			PhoneNumber: "01011111111",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetParentUUIDsWithInform": {[]string{}, gorm.ErrRecordNotFound},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeParentWithThatInformNoExist,
		}, { // GetParentUUIDsWithInform error return
			UUID:        "parent-111111111111",
			PhoneNumber: "01012341234",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetParentUUIDsWithInform": {[]string{}, errors.New("I don't know about that error")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.GetParentUUIDsWithInformRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.GetParentUUIDsWithInformResponse)
		_ = defaultHandler.GetParentUUIDsWithInform(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedParentUUIDs, resp.ParentUUIDs, "result parentUUIDs assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	newMock.AssertExpectations(t)
}
