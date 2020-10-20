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

func Test_default_LoginStudentAuth(t *testing.T) {
	newMock, defaultHandler := generateVarForTest()
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
			ExpectedCode:   code.StudentIDNoExist,
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
			ExpectedCode:   code.IncorrectStudentPWForLogin,
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		var req = new(proto.LoginStudentAuthRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		var resp = new(proto.LoginStudentAuthResponse)
		_ = defaultHandler.LoginStudentAuth(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedLoggedInStudentUUID, resp.LoggedInStudentUUID, "student uuid assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	newMock.AssertExpectations(t)
}

func Test_default_ChangeStudentPW(t *testing.T) {
	newMock, defaultHandler := generateVarForTest()

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
			UUID:        "admin-111111111111",
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
			ExpectedCode:   code.IncorrectStudentPWForChange,
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
			UUID:        "admin-111111111112",
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
		testCase.OnExpectMethods(newMock)

		req := new(proto.ChangeStudentPWRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.ChangeStudentPWResponse)
		_ = defaultHandler.ChangeStudentPW(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	newMock.AssertExpectations(t)
}

func Test_default_GetStudentInformWithUUID(t *testing.T) {
	newMock, defaultHandler := generateVarForTest()
	now := time.Now()

	tests := []test.GetStudentInformWithUUIDCase{
		{ // success case
			UUID: "student-111111111111",
			StudentUUID: "student-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentInformWithUUID": {&model.StudentInform{
					Model:         gorm.Model{CreatedAt: now, UpdatedAt: now},
					StudentUUID:   "student-111111111111",
					Grade:         2,
					Class:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01088378347",
					ProfileURI:    "/profiles/student-111111111111",
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedInform: &model.StudentInform{
				Grade:         2,
				Class:         2,
				StudentNumber: 7,
				Name:          "박진홍",
				PhoneNumber:   "01088378347",
				ProfileURI:    "/profiles/student-111111111111",
			},
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			XRequestID:      test.EmptyReplaceValueForString,
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
			ExpectedInform:  &model.StudentInform{},
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
			ExpectedInform:  &model.StudentInform{},
		}, { // no exist Span-Context -> Proxy Authorization Required
			SpanContextString: test.EmptyReplaceValueForString,
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
			ExpectedInform:    &model.StudentInform{},
		}, { // invalid Span-Context -> Proxy Authorization Required
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
			ExpectedInform:    &model.StudentInform{},
		}, { // forbidden (not student)
			UUID:           "parent-111111111112",
			StudentUUID:    "student-111111111112",
			ExpectedStatus: http.StatusForbidden,
			ExpectedInform: &model.StudentInform{},
		}, { // forbidden (not my auth)
			UUID:           "student-111111111113",
			StudentUUID:    "student-111111111114",
			ExpectedStatus: http.StatusForbidden,
			ExpectedInform: &model.StudentInform{},
		}, { // no exist student uuid
			UUID:        "admin-111111111111",
			StudentUUID: "student-111111111115",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetStudentInformWithUUID": {&model.StudentInform{}, gorm.ErrRecordNotFound},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedInform: &model.StudentInform{},
		}, { // GetStudentInformWithUUID error return
			UUID:        "admin-111111111112",
			StudentUUID: "student-111111111116",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetStudentInformWithUUID": {&model.StudentInform{}, errors.New("DB not connected")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedInform: &model.StudentInform{},
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.GetStudentInformWithUUIDRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.GetStudentInformWithUUIDResponse)
		_ = defaultHandler.GetStudentInformWithUUID(ctx, req, resp)

		resultInform := &model.StudentInform{
			Grade:         model.Grade(int64(resp.Grade)),
			Class:         model.Class(int64(resp.Group)),
			StudentNumber: model.StudentNumber(int64(resp.StudentNumber)),
			Name:          model.Name(resp.Name),
			PhoneNumber:   model.PhoneNumber(resp.PhoneNumber),
			ProfileURI:    model.ProfileURI(resp.ImageURI),
		}

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedInform, resultInform, "result inform assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	newMock.AssertExpectations(t)
}

func Test_default_GetStudentUUIDsWithInform(t *testing.T) {
	newMock, defaultHandler := generateVarForTest()

	tests := []test.GetStudentUUIDsWithInformCase{
		{ // success case (for admin auth)
			UUID: "admin-111111111111",
			Name: "이성진",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetStudentUUIDsWithInform": {[]string{"student-123412341234", "student-123412341234"}, nil},
				"Commit":                    {&gorm.DB{}},
			},
			ExpectedStatus:       http.StatusOK,
			ExpectedStudentUUIDs: []string{"student-123412341234", "student-123412341234"},
		}, { // success case (for student auth)
			UUID:          "student-111111111111",
			Class:         2,
			Grade:         2,
			StudentNumber: 7,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetStudentUUIDsWithInform": {[]string{"student-111111111111"}, nil},
				"Commit":                    {&gorm.DB{}},
			},
			ExpectedStatus:       http.StatusOK,
			ExpectedStudentUUIDs: []string{"student-111111111111"},
		}, { // no exist X-Request-ID -> Proxy Authorization Required
			UUID:            "student-111111111111",
			XRequestID:      test.EmptyReplaceValueForString,
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // invalid X-Request-ID -> Proxy Authorization Required
			UUID:            "student-111111111111",
			XRequestID:      "InvalidXRequestID",
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusProxyAuthRequired,
		}, { // no exist Span-Context -> Proxy Authorization Required
			UUID:              "student-111111111111",
			SpanContextString: test.EmptyReplaceValueForString,
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			UUID:              "student-111111111111",
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[test.Method]test.Returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // forbidden (not student)
			UUID:           "parent-111111111111",
			ExpectedStatus: http.StatusForbidden,
		}, { // no exist student uuid with that inform
			UUID:          "student-111111111111",
			Class:         2,
			Grade:         2,
			StudentNumber: 21,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetStudentUUIDsWithInform": {[]string{}, gorm.ErrRecordNotFound},
				"Rollback":                  {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.StudentWithThatInformNoExist,
		}, { // GetStudentInformWithUUID error return
			UUID:          "student-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                   {},
				"GetStudentUUIDsWithInform": {[]string{}, errors.New("I don't know about that error")},
				"Rollback":                  {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range tests {
		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.GetStudentUUIDsWithInformRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.GetStudentUUIDsWithInformResponse)
		_ = defaultHandler.GetStudentUUIDsWithInform(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedStudentUUIDs, resp.StudentUUIDs, "result studentUUIDs assertion error (test case: %v, message: %s)", testCase, resp.Message)
	}

	newMock.AssertExpectations(t)
}

func Test_default_GetStudentInformsWithUUIDs(t *testing.T) {
	now := time.Now()

	tests := []test.GetStudentInformsWithUUIDsCase{
		{ // success case
			UUID:         "student-111111111111",
			StudentUUIDs: []string{"student-111111111111", "student-222222222222"},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentInformsWithUUIDs": {[]*model.StudentInform{
					{
						Model:         gorm.Model{CreatedAt: now, UpdatedAt: now},
						StudentUUID:   "student-111111111111",
						Grade:         2,
						Class:         2,
						StudentNumber: 7,
						Name:          "박진홍",
						PhoneNumber:   "01011111111",
						ProfileURI:    "/profiles/student-111111111111",
					}, {
						Model:         gorm.Model{CreatedAt: now, UpdatedAt: now},
						StudentUUID:   "student-222222222222",
						Grade:         1,
						Class:         2,
						StudentNumber: 8,
						Name:          "박진홍",
						PhoneNumber:   "01022222222",
						ProfileURI:    "/profiles/student-222222222222",
					},
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusOK,
			ExpectedInforms: []*proto.StudentInform{
				{
					StudentUUID:   "student-111111111111",
					Grade:         2,
					Group:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01011111111",
					ImageURI:      "/profiles/student-111111111111",
				}, {
					StudentUUID:   "student-222222222222",
					Grade:         1,
					Group:         2,
					StudentNumber: 8,
					Name:          "박진홍",
					PhoneNumber:   "01022222222",
					ImageURI:      "/profiles/student-222222222222",
				},
			},
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
			StudentUUIDs:   []string{"student-111111111112"},
			ExpectedStatus: http.StatusForbidden,
		}, { // no exist student uuid
			UUID:         "admin-111111111111",
			StudentUUIDs: []string{"student-111111111115", "student-111111111111"},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetStudentInformsWithUUIDs": {[]*model.StudentInform{{}, {
					Model:         gorm.Model{CreatedAt: now, UpdatedAt: now},
					StudentUUID:   "student-111111111111",
					Grade:         2,
					Class:         2,
					StudentNumber: 7,
					Name:          "박진홍",
					PhoneNumber:   "01011111111",
					ProfileURI:    "/profiles/student-111111111111",
				}}, gorm.ErrRecordNotFound},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.StudentUUIDsContainNoExistUUID,
		}, { // GetStudentInformWithUUID error return
			UUID:         "admin-111111111112",
			StudentUUIDs: []string{},
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                    {},
				"GetStudentInformsWithUUIDs": {[]*model.StudentInform{}, errors.New("DB not connected")},
				"Rollback":                   {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range tests {
		newMock, defaultHandler := generateVarForTest()

		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		req := new(proto.GetStudentInformsWithUUIDsRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		resp := new(proto.GetStudentInformsWithUUIDsResponse)
		_ = defaultHandler.GetStudentInformsWithUUIDs(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedInforms, resp.StudentInforms, "result informs assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}
