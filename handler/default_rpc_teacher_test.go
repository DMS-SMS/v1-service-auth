package handler

import (
	test "auth/handler/for_test"
	"auth/model"
	proto "auth/proto/golang/auth"
	code "auth/utils/code/golang"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"testing"
	"time"
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
			ExpectedCode:   code.TeacherIDNoExist,
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
			ExpectedCode:   code.IncorrectTeacherPWForLogin,
		},
	}

	for _, testCase := range tests {
		newMock, defaultHandler := generateVarForTest()

		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		var req = new(proto.LoginTeacherAuthRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		var resp = new(proto.LoginTeacherAuthResponse)
		_ = defaultHandler.LoginTeacherAuth(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedLoggedInTeacherUUID, resp.LoggedInTeacherUUID, "logged in uuid assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}

func Test_default_ChangeTeacherPW(t *testing.T) {
	hashedTestPW, _ := bcrypt.GenerateFromPassword([]byte("testPW"), 1)

	tests := []test.ChangeTeacherPWCase{
		{ // success case
			UUID:        "teacher-111111111111",
			TeacherUUID: "teacher-111111111111",
			CurrentPW:   "testPW",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetTeacherAuthWithUUID": {&model.TeacherAuth{
					UUID:     "teacher-111111111111",
					TeacherPW: model.TeacherPW(string(hashedTestPW)),
				}, nil},
				"ChangeTeacherPW": {nil},
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
		}, { // forbidden (not teacher)
			UUID:           "student-111111111112",
			TeacherUUID:    "student-111111111112",
			CurrentPW:      "testPW",
			RevisionPW:     "NewPassword",
			ExpectedStatus: http.StatusForbidden,
		}, { // forbidden (not my auth)
			UUID:           "teacher-111111111113",
			TeacherUUID:    "teacher-111111111114",
			CurrentPW:      "testPW",
			RevisionPW:     "NewPassword",
			ExpectedStatus: http.StatusForbidden,
		}, { // not exists student
			UUID:           "teacher-111111111115",
			TeacherUUID:    "teacher-111111111115",
			CurrentPW:      "testPW",
			RevisionPW:     "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetTeacherAuthWithUUID": {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
		}, { // 현재 Password 불일치
			UUID:        "teacher-111111111116",
			TeacherUUID: "teacher-111111111116",
			CurrentPW:   "IncorrectPassword",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetTeacherAuthWithUUID": {&model.TeacherAuth{
					UUID:      "teacher-111111111116",
					TeacherPW: model.TeacherPW(string(hashedTestPW)),
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.IncorrectTeacherPWForChange,
		}, { // GetStudentAuthWithUUID 에러 반환
			UUID:        "admin-111111111111",
			TeacherUUID: "teacher-111111111117",
			CurrentPW:   "testPW",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetTeacherAuthWithUUID": {&model.TeacherAuth{}, errors.New("DB not connected")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // ChangeStudentPW 에러 반환
			UUID:        "teacher-111111111118",
			TeacherUUID: "teacher-111111111118",
			CurrentPW:   "testPW",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetTeacherAuthWithUUID": {&model.TeacherAuth{
					UUID:      "teacher-111111111118",
					TeacherPW: model.TeacherPW(string(hashedTestPW)),
				}, nil},
				"ChangeTeacherPW": {errors.New("DB not connected")},
				"Rollback":        {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // GetStudentAuthWithUUID Short Hashed PW 반환
			UUID:        "admin-111111111112",
			TeacherUUID: "teacher-111111111119",
			CurrentPW:   "testPW",
			RevisionPW:  "NewPassword",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetTeacherAuthWithUUID": {&model.TeacherAuth{
					UUID:      "teacher-111111111119",
					TeacherPW: "TooShortHashedPasword",
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range tests {
		newMock, defaultHandler := generateVarForTest()

		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.ChangeTeacherPWRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.ChangeTeacherPWResponse)
		_ = defaultHandler.ChangeTeacherPW(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}

func Test_default_GetTeacherInformWithUUID(t *testing.T) {
	now := time.Now()

	tests := []test.GetTeacherInformWithUUIDCase{
		{ // success case
			UUID: "teacher-111111111111",
			TeacherUUID: "teacher-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetTeacherInformWithUUID": {&model.TeacherInform{
					Model:         gorm.Model{CreatedAt: now, UpdatedAt: now},
					TeacherUUID:   "teacher-111111111111",
					Grade:         2,
					Class:         2,
					Name:          "박진홍",
					PhoneNumber:   "01088378347",
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedInform: &model.TeacherInform{
				Grade:         2,
				Class:         2,
				Name:          "박진홍",
				PhoneNumber:   "01088378347",
			},
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			XRequestID:      test.EmptyReplaceValueForString,
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
			ExpectedInform:  &model.TeacherInform{},
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
			ExpectedInform:  &model.TeacherInform{},
		}, { // no exist Span-Context -> Proxy Authorization Required
			SpanContextString: test.EmptyReplaceValueForString,
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
			ExpectedInform:    &model.TeacherInform{},
		}, { // invalid Span-Context -> Proxy Authorization Required
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
			ExpectedInform:    &model.TeacherInform{},
		}, { // forbidden (not student)
			UUID:           "parent-111111111112",
			TeacherUUID:    "parent-111111111112",
			ExpectedStatus: http.StatusForbidden,
			ExpectedInform: &model.TeacherInform{},
		}, { // forbidden (not my auth)
			UUID:           "teacher-111111111113",
			TeacherUUID:    "teacher-111111111114",
			ExpectedStatus: http.StatusForbidden,
			ExpectedInform: &model.TeacherInform{},
		}, { // no exist student uuid
			UUID:        "admin-111111111111",
			TeacherUUID: "teacher-111111111115",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherInformWithUUID": {&model.TeacherInform{}, gorm.ErrRecordNotFound},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedInform: &model.TeacherInform{},
		}, { // GetStudentInformWithUUID error return
			UUID:        "admin-111111111112",
			TeacherUUID: "teacher-111111111116",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherInformWithUUID": {&model.TeacherInform{}, errors.New("DB not connected")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedInform: &model.TeacherInform{},
		},
	}

	for _, testCase := range tests {
		newMock, defaultHandler := generateVarForTest()

		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.GetTeacherInformWithUUIDRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.GetTeacherInformWithUUIDResponse)
		_ = defaultHandler.GetTeacherInformWithUUID(ctx, req, resp)

		resultInform := &model.TeacherInform{
			Grade:         model.Grade(int64(resp.Grade)),
			Class:         model.Class(int64(resp.Group)),
			Name:          model.Name(resp.Name),
			PhoneNumber:   model.PhoneNumber(resp.PhoneNumber),
		}

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedInform, resultInform, "result inform assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}

func Test_default_GetTeacherUUIDsWithInform(t *testing.T) {
	tests := []test.GetTeacherUUIDsWithInformCase{
		{ // success case (for admin auth)
			UUID: "admin-111111111111",
			Name: "이성진",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetTeacherUUIDsWithInform": {[]string{"teacher-123412341234", "teacher-123412341234"}, nil},
				"Commit":                    {&gorm.DB{}},
			},
			ExpectedStatus:       http.StatusOK,
			ExpectedTeacherUUIDs: []string{"teacher-123412341234", "teacher-123412341234"},
		}, { // success case (for teacher auth)
			UUID:  "teacher-111111111111",
			Class: 2,
			Grade: 2,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetTeacherUUIDsWithInform": {[]string{"teacher-111111111111"}, nil},
				"Commit":                    {&gorm.DB{}},
			},
			ExpectedStatus:       http.StatusOK,
			ExpectedTeacherUUIDs: []string{"teacher-111111111111"},
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			UUID:            "teacher-111111111111",
			XRequestID:      test.EmptyReplaceValueForString,
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			UUID:            "teacher-111111111111",
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // no exist Span-Context -> Proxy Authorization Required
			UUID:              "teacher-111111111111",
			SpanContextString: test.EmptyReplaceValueForString,
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			UUID:              "teacher-111111111111",
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // forbidden (not teacher)
			UUID:           "student-111111111111",
			ExpectedStatus: http.StatusForbidden,
		}, { // no exist parent uuid with that inform
			UUID:  "admin-111111111111",
			Class: 2,
			Grade: 12,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetTeacherUUIDsWithInform": {[]string{}, gorm.ErrRecordNotFound},
				"Rollback":                  {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.TeacherWithThatInformNoExist,
		}, { // GetTeacherUUIDsWithInform error return
			UUID: "admin-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetTeacherUUIDsWithInform": {[]string{}, errors.New("I don't know about that error")},
				"Rollback":                  {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range tests {
		newMock, defaultHandler := generateVarForTest()

		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.GetTeacherUUIDsWithInformRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.GetTeacherUUIDsWithInformResponse)
		_ = defaultHandler.GetTeacherUUIDsWithInform(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedTeacherUUIDs, resp.TeacherUUIDs, "result teacherUUIDs assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}

func Test_default_GetTeacherUUIDsWithInform(t *testing.T) {
	newMock, defaultHandler := generateVarForTest()

	tests := []test.GetTeacherUUIDsWithInformCase{
		{ // success case (for admin auth)
			UUID: "admin-111111111111",
			Name: "이성진",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetTeacherUUIDsWithInform": {[]string{"teacher-123412341234", "teacher-123412341234"}, nil},
				"Commit":                    {&gorm.DB{}},
			},
			ExpectedStatus:       http.StatusOK,
			ExpectedTeacherUUIDs: []string{"teacher-123412341234", "teacher-123412341234"},
		}, { // success case (for teacher auth)
			UUID:  "teacher-111111111111",
			Class: 2,
			Grade: 2,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetTeacherUUIDsWithInform": {[]string{"teacher-111111111111"}, nil},
				"Commit":                    {&gorm.DB{}},
			},
			ExpectedStatus:       http.StatusOK,
			ExpectedTeacherUUIDs: []string{"teacher-111111111111"},
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			UUID:            "teacher-111111111111",
			XRequestID:      test.EmptyReplaceValueForString,
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			UUID:            "teacher-111111111111",
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // no exist Span-Context -> Proxy Authorization Required
			UUID:              "teacher-111111111111",
			SpanContextString: test.EmptyReplaceValueForString,
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			UUID:              "teacher-111111111111",
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // forbidden (not teacher)
			UUID:           "student-111111111111",
			ExpectedStatus: http.StatusForbidden,
		}, { // no exist parent uuid with that inform
			UUID:  "teacher-111111111111",
			Class: 2,
			Grade: 12,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetTeacherUUIDsWithInform": {[]string{}, gorm.ErrRecordNotFound},
				"Rollback":                  {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.TeacherWithThatInformNoExist,
		}, { // GetTeacherUUIDsWithInform error return
			UUID: "teacher-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetTeacherUUIDsWithInform": {[]string{}, errors.New("I don't know about that error")},
				"Rollback":                  {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.GetTeacherUUIDsWithInformRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.GetTeacherUUIDsWithInformResponse)
		_ = defaultHandler.GetTeacherUUIDsWithInform(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedTeacherUUIDs, resp.TeacherUUIDs, "result teacherUUIDs assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	newMock.AssertExpectations(t)
}
