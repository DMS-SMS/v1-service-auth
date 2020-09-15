package handler

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"auth/tool/mysqlerr"
	"errors"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

//func init() {
//	s, err := session.NewSession(&aws.Config{
//		Credentials: credentials.NewStaticCredentials("eks-user-1-id", "eks-user-1-key", ""),
//		Region:      aws.String("ap-northeast-2"),
//	})
//	if err != nil { panic(err) }
//
//	_, err = s3.New(s).PutObject(&s3.PutObjectInput{
//		Bucket:      aws.String("dms-sms"),
//		Key:         aws.String(fmt.Sprintf("profiles/%s", "student-111111111111")),
//		Body:        bytes.NewReader(validImageByteArr),
//	})
//	if err != nil { panic(err) }
//}

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
			SpanContextString: emptyReplaceValueForString,
			ExpectedMethods:   map[method]returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[method]returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
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
			ExpectedCode:   CodeStudentPhoneNumberDuplicate,
		}, { // CheckIfStudentAuthExists error occur
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, errors.New("unexpected error from DB Connection")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return invalid duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return unexpected key duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, mysqlerr.DuplicateEntry("UnexpectedKey", "error")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return invalid Fk Constraint Fail error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_NO_REFERENCED_ROW_2, Message: "InvalidMessage"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return unexpected constraint name error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
					ConstraintName: "unexpected constraint name",
					AttrName:       "unexpected attr",
				}, mysqlerr.RefInform{})},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return unexpected constraint name error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
					ConstraintName: "unexpected constraint name",
					AttrName:       "unexpected attr",
				}, mysqlerr.RefInform{})},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return unexpected error code
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentInform return invalid duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentInform return unexpected duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, mysqlerr.DuplicateEntry("UnexpectedKey", "duplicated")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentInform return unexpected error code
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfStudentAuthExists": {false, nil},
				"CreateStudentAuth":        {&model.StudentAuth{}, nil},
				"CreateStudentInform":      {&model.StudentInform{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		test.ChangeEmptyValueToValidValue()
		test.ChangeEmptyReplaceValueToEmptyValue()
		test.OnExpectMethodsTo(mockForDB)

		req := new(proto.CreateNewStudentRequest)
		test.SetRequestContextOf(req)
		ctx := test.GetMetadataContext()

		resp := new(proto.CreateNewStudentResponse)
		_ = defaultHandler.CreateNewStudent(ctx, req, resp)

		test.Image = nil
		assert.Equalf(t, int(test.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", test, resp.Message)
		assert.Equalf(t, test.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", test, resp.Message)
		assert.Regexpf(t, test.ExpectedStudentUUID, resp.CreatedStudentUUID, "student uuid assertion error (test case: %v, message: %s)", test, resp.Message)
	}

	mockForDB.AssertExpectations(t)
}

func Test_default_CreateNewTeacher(t *testing.T) {
	const teacherUUIDRegexString = "^teacher-\\d{12}"

	tests := []createNewTeacherTest{
		{ // success case
			Grade: emptyReplaceValueForUint32,
			Class: emptyReplaceValueForUint32,
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, nil},
				"Commit":                   {&gorm.DB{}},
			},
			ExpectedStatus:      http.StatusCreated,
			ExpectedStudentUUID: teacherUUIDRegexString,
		}, { // not admin uuid -> forbidden
			UUID:            "NotAdminAuthUUID", // (admin-숫자 12개의 형식이여야 함)
			ExpectedMethods: map[method]returns{},
			ExpectedStatus:  http.StatusForbidden,
		}, { // invalid request value -> Proxy Authorization Required
			TeacherID: "유효하지 않은 아이디", // ASCII, 4~16 사이 문자열이여야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Grade: 100, // 1~3 사이의 숫자여야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Name: "Invalid Name", // 2~4 글자의 한글이어야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, (validator.ValidationErrors)(nil)},
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
			SpanContextString: emptyReplaceValueForString,
			ExpectedMethods:   map[method]returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[method]returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // student id duplicate -> Conflict -201
			TeacherID: "duplicateID",
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, mysqlerr.DuplicateEntry(model.TeacherAuthInstance.TeacherID.KeyName(), "duplicateID")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeTeacherIDDuplicate,
		}, { // phone number duplicate -> Conflict -202
			PhoneNumber: "01088378347",
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, mysqlerr.DuplicateEntry(model.TeacherInformInstance.PhoneNumber.KeyName(), "01088378347")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeTeacherPhoneNumberDuplicate,
		}, { // CheckIfTeacherAuthExists error occur
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, errors.New("unexpected error from DB Connection")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return invalid duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return unexpected key duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, mysqlerr.DuplicateEntry("UnexpectedKey", "error")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return unexpected error code
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return invalid duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return unexpected duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, mysqlerr.DuplicateEntry("UnexpectedKey", "duplicated")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return unexpected error code
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfTeacherAuthExists": {false, nil},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		test.ChangeEmptyValueToValidValue()
		test.ChangeEmptyReplaceValueToEmptyValue()
		test.OnExpectMethodsTo(mockForDB)

		req := new(proto.CreateNewTeacherRequest)
		test.SetRequestContextOf(req)
		ctx := test.GetMetadataContext()

		resp := new(proto.CreateNewTeacherResponse)
		_ = defaultHandler.CreateNewTeacher(ctx, req, resp)

		assert.Equalf(t, int(test.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", test, resp.Message)
		assert.Equalf(t, test.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", test, resp.Message)
		assert.Regexpf(t, test.ExpectedStudentUUID, resp.CreatedTeacherUUID, "teacher uuid assertion error (test case: %v, message: %s)", test, resp.Message)
	}

	mockForDB.AssertExpectations(t)
}

func Test_default_CreateNewParent(t *testing.T) {
	const parentUUIDRegexString = "^parent-\\d{12}"

	tests := []createNewParentTest{
		{ // success case
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, nil},
				"CreateParentInform":      {&model.ParentInform{}, nil},
				"Commit":                  {&gorm.DB{}},
			},
			ExpectedStatus:      http.StatusCreated,
			ExpectedStudentUUID: parentUUIDRegexString,
		}, { // not admin uuid -> forbidden
			UUID:            "NotAdminAuthUUID", // (admin-숫자 12개의 형식이여야 함)
			ExpectedMethods: map[method]returns{},
			ExpectedStatus:  http.StatusForbidden,
		}, { // invalid request value -> Proxy Authorization Required
			ParentID: "유효하지 않은 아이디", // ASCII, 4~16 사이 문자열이여야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, (validator.ValidationErrors)(nil)},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Name: "Invalid Name", // 2~4 글자의 한글이어야 함
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, nil},
				"CreateParentInform":      {&model.ParentInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":                {&gorm.DB{}},
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
			SpanContextString: emptyReplaceValueForString,
			ExpectedMethods:   map[method]returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // invalid Span-Context -> Proxy Authorization Required
			SpanContextString: "InvalidSpanContext",
			ExpectedMethods:   map[method]returns{},
			ExpectedStatus:    http.StatusProxyAuthRequired,
		}, { // student id duplicate -> Conflict -201
			ParentID: "duplicateID",
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, mysqlerr.DuplicateEntry(model.ParentAuthInstance.ParentID.KeyName(), "duplicateID")},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeParentIDDuplicate,
		}, { // phone number duplicate -> Conflict -202
			PhoneNumber: "01088378347",
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, nil},
				"CreateParentInform":      {&model.ParentInform{}, mysqlerr.DuplicateEntry(model.ParentInformInstance.PhoneNumber.KeyName(), "01088378347")},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   CodeParentPhoneNumberDuplicate,
		}, { // CheckIfTeacherAuthExists error occur
			ExpectedMethods: map[method]returns{
				"BeginTx":                  {},
				"CheckIfParentAuthExists": {false, errors.New("unexpected error from DB Connection")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return invalid duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return unexpected key duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, mysqlerr.DuplicateEntry("UnexpectedKey", "error")},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return unexpected error code
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return invalid duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, nil},
				"CreateParentInform":      {&model.ParentInform{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return unexpected duplicate error
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, nil},
				"CreateParentInform":      {&model.ParentInform{}, mysqlerr.DuplicateEntry("UnexpectedKey", "duplicated")},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return unexpected error code
			ExpectedMethods: map[method]returns{
				"BeginTx":                 {},
				"CheckIfParentAuthExists": {false, nil},
				"CreateParentAuth":        {&model.ParentAuth{}, nil},
				"CreateParentInform":      {&model.ParentInform{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":                {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		test.ChangeEmptyValueToValidValue()
		test.ChangeEmptyReplaceValueToEmptyValue()
		test.OnExpectMethodsTo(mockForDB)

		req := new(proto.CreateNewParentRequest)
		test.SetRequestContextOf(req)
		ctx := test.GetMetadataContext()

		resp := new(proto.CreateNewParentResponse)
		_ = defaultHandler.CreateNewParent(ctx, req, resp)

		assert.Equalf(t, int(test.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", test, resp.Message)
		assert.Equalf(t, test.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", test, resp.Message)
		assert.Regexpf(t, test.ExpectedStudentUUID, resp.CreatedParentUUID, "parent uuid assertion error (test case: %v, message: %s)", test, resp.Message)
	}

	mockForDB.AssertExpectations(t)
}