package handler

import (
	test "auth/handler/for_test"
	"auth/model"
	proto "auth/proto/golang/auth"
	"auth/tool/mysqlerr"
	code "auth/utils/code/golang"
	"errors"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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

	tests := []test.CreateNewStudentCase{
		{ // success case
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"CreateStudentInform":    {&model.StudentInform{}, nil},
				"Commit":                 {&gorm.DB{}},
			},
			ExpectedStatus:      http.StatusCreated,
			ExpectedStudentUUID: studentUUIDRegexString,
		}, { // not admin uuid -> forbidden
			UUID:            "NotAdminAuthUUID", // (admin-숫자 12개의 형식이여야 함)
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusForbidden,
		}, { // invalid request value -> Proxy Authorization Required
			StudentID: "유효하지 않은 아이디", // ASCII, 4~16 사이 문자열이여야 함
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, (validator.ValidationErrors)(nil)},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Grade: 100, // 1~3 사이의 숫자여야 함
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"CreateStudentInform":    {&model.StudentInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Name: "Invalid Name", // 2~4 글자의 한글이어야 함
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"CreateStudentInform":    {&model.StudentInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
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
		}, { // student id duplicate -> Conflict -101
			StudentID: "jinhong0719",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, mysqlerr.DuplicateEntry(model.StudentAuthInstance.StudentID.KeyName(), "jinhong0719")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.StudentIDDuplicate,
		}, { // parent uuid fk constraint fail -> Conflict -102
			ParentUUID: "parent-111111111111",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, test.StudentAuthParentUUIDFKConstraintFailError},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.ParentUUIDNoExist,
		}, { // image empty byte array
			Image: []byte(test.EmptyReplaceValueForString),
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // student number duplicate -> Conflict -103
			Grade:         2,
			Class:         2,
			StudentNumber: 7,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"CreateStudentInform":    {&model.StudentInform{}, mysqlerr.DuplicateEntry(model.StudentInformInstance.StudentNumber.KeyName(), "2207")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.StudentNumberDuplicate,
		}, { // phone number duplicate -> Conflict -104
			PhoneNumber: "01088378347",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"CreateStudentInform":    {&model.StudentInform{}, mysqlerr.DuplicateEntry(model.StudentInformInstance.PhoneNumber.KeyName(), "01088378347")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.StudentPhoneNumberDuplicate,
		}, { // CheckIfStudentAuthExists error occur
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, errors.New("unexpected error from DB Connection")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return invalid duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return unexpected key duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, mysqlerr.DuplicateEntry("UnexpectedKey", "error")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return invalid Fk Constraint Fail error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_NO_REFERENCED_ROW_2, Message: "InvalidMessage"}},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return unexpected constraint name error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth": {&model.StudentAuth{}, mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
					ConstraintName: "unexpected constraint name",
					AttrName:       "unexpected attr",
				}, mysqlerr.RefInform{})},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return unexpected constraint name error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth": {&model.StudentAuth{}, mysqlerr.FKConstraintFailWithoutReferenceInform(mysqlerr.FKInform{
					ConstraintName: "unexpected constraint name",
					AttrName:       "unexpected attr",
				}, mysqlerr.RefInform{})},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return unexpected error code
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentAuth return unexpected type of error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, errors.New("unexpected type of error")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentInform return invalid duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"CreateStudentInform":    {&model.StudentInform{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentInform return unexpected duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"CreateStudentInform":    {&model.StudentInform{}, mysqlerr.DuplicateEntry("UnexpectedKey", "duplicated")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentInform return unexpected error code
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"CreateStudentInform":    {&model.StudentInform{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateStudentInform return unexpected type of error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                {},
				"GetStudentAuthWithUUID": {&model.StudentAuth{}, gorm.ErrRecordNotFound},
				"CreateStudentAuth":      {&model.StudentAuth{}, nil},
				"CreateStudentInform":    {&model.StudentInform{}, errors.New("unexpected error")},
				"Rollback":               {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, createNewStudentTest := range tests {
		newMock, defaultHandler := generateVarForTest()

		createNewStudentTest.ChangeEmptyValueToValidValue()
		createNewStudentTest.ChangeEmptyReplaceValueToEmptyValue()
		createNewStudentTest.OnExpectMethodsTo(newMock)

		req := new(proto.CreateNewStudentRequest)
		createNewStudentTest.SetRequestContextOf(req)
		ctx := createNewStudentTest.GetMetadataContext()

		resp := new(proto.CreateNewStudentResponse)
		_ = defaultHandler.CreateNewStudent(ctx, req, resp)

		createNewStudentTest.Image = nil
		assert.Equalf(t, int(createNewStudentTest.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", createNewStudentTest, resp.Message)
		assert.Equalf(t, createNewStudentTest.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", createNewStudentTest, resp.Message)
		assert.Regexpf(t, createNewStudentTest.ExpectedStudentUUID, resp.CreatedStudentUUID, "student uuid assertion error (test case: %v, message: %s)", createNewStudentTest, resp.Message)

		newMock.AssertExpectations(t)
	}
}

func Test_default_CreateNewTeacher(t *testing.T) {
	const teacherUUIDRegexString = "^teacher-\\d{12}"

	tests := []test.CreateNewTeacherCase{
		{ // success case
			Grade: test.EmptyReplaceValueForUint32,
			Class: test.EmptyReplaceValueForUint32,
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, nil},
				"Commit":                   {&gorm.DB{}},
			},
			ExpectedStatus:      http.StatusCreated,
			ExpectedStudentUUID: teacherUUIDRegexString,
		}, { // not admin uuid -> forbidden
			UUID:            "NotAdminAuthUUID", // (admin-숫자 12개의 형식이여야 함)
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusForbidden,
		}, { // invalid request value -> Proxy Authorization Required
			TeacherID: "유효하지 않은 아이디", // ASCII, 4~16 사이 문자열이여야 함
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Grade: 100, // 1~3 사이의 숫자여야 함
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Name: "Invalid Name", // 2~4 글자의 한글이어야 함
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
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
		}, { // student id duplicate -> Conflict -201
			TeacherID: "duplicateID",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, mysqlerr.DuplicateEntry(model.TeacherAuthInstance.TeacherID.KeyName(), "duplicateID")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.TeacherIDDuplicate,
		}, { // phone number duplicate -> Conflict -202
			PhoneNumber: "01088378347",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, mysqlerr.DuplicateEntry(model.TeacherInformInstance.PhoneNumber.KeyName(), "01088378347")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.TeacherPhoneNumberDuplicate,
		}, { // CheckIfTeacherAuth1Exists error occur
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, errors.New("unexpected error from DB Connection")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return invalid duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return unexpected key duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, mysqlerr.DuplicateEntry("UnexpectedKey", "error")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return unexpected error code
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherAuth return unexpected type of error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, errors.New("unexpected type of error")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return invalid duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return unexpected duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, mysqlerr.DuplicateEntry("UnexpectedKey", "duplicated")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return unexpected error code
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateTeacherInform return unexpected type of error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":                  {},
				"GetTeacherAuthWithUUID":   {&model.TeacherAuth{}, gorm.ErrRecordNotFound},
				"CreateTeacherAuth":        {&model.TeacherAuth{}, nil},
				"CreateTeacherInform":      {&model.TeacherInform{}, errors.New("unexpected type of error")},
				"Rollback":                 {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, createNewTeacherTest := range tests {
		newMock, defaultHandler := generateVarForTest()

		createNewTeacherTest.ChangeEmptyValueToValidValue()
		createNewTeacherTest.ChangeEmptyReplaceValueToEmptyValue()
		createNewTeacherTest.OnExpectMethodsTo(newMock)

		req := new(proto.CreateNewTeacherRequest)
		createNewTeacherTest.SetRequestContextOf(req)
		ctx := createNewTeacherTest.GetMetadataContext()

		resp := new(proto.CreateNewTeacherResponse)
		_ = defaultHandler.CreateNewTeacher(ctx, req, resp)

		assert.Equalf(t, int(createNewTeacherTest.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", createNewTeacherTest, resp.Message)
		assert.Equalf(t, createNewTeacherTest.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", createNewTeacherTest, resp.Message)
		assert.Regexpf(t, createNewTeacherTest.ExpectedStudentUUID, resp.CreatedTeacherUUID, "teacher uuid assertion error (test case: %v, message: %s)", createNewTeacherTest, resp.Message)

		newMock.AssertExpectations(t)
	}
}

func Test_default_CreateNewParent(t *testing.T) {
	const parentUUIDRegexString = "^parent-\\d{12}"

	tests := []test.CreateNewParentCase{
		{ // success case
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, nil},
				"CreateParentInform":    {&model.ParentInform{}, nil},
				"Commit":                {&gorm.DB{}},
			},
			ExpectedStatus:      http.StatusCreated,
			ExpectedStudentUUID: parentUUIDRegexString,
		}, { // not admin uuid -> forbidden
			UUID:            "NotAdminAuthUUID", // (admin-숫자 12개의 형식이여야 함)
			ExpectedMethods: map[test.Method]test.Returns{},
			ExpectedStatus:  http.StatusForbidden,
		}, { // invalid request value -> Proxy Authorization Required
			ParentID: "유효하지 않은 아이디", // ASCII, 4~16 사이 문자열이여야 함
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, (validator.ValidationErrors)(nil)},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
		}, { // invalid request value -> Proxy Authorization Required
			Name: "Invalid Name", // 2~4 글자의 한글이어야 함
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, nil},
				"CreateParentInform":    {&model.ParentInform{}, (validator.ValidationErrors)(nil)},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusProxyAuthRequired,
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
		}, { // student id duplicate -> Conflict -201
			ParentID: "duplicateID",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, mysqlerr.DuplicateEntry(model.ParentAuthInstance.ParentID.KeyName(), "duplicateID")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.ParentIDDuplicate,
		}, { // phone number duplicate -> Conflict -202
			PhoneNumber: "01088378347",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, nil},
				"CreateParentInform":    {&model.ParentInform{}, mysqlerr.DuplicateEntry(model.ParentInformInstance.PhoneNumber.KeyName(), "01088378347")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.ParentPhoneNumberDuplicate,
		}, { // GetParentAuth1WithUUID error occur
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, errors.New("unexpected error from DB Connection")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateParentAuth return invalid duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateParentAuth return unexpected key duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, mysqlerr.DuplicateEntry("UnexpectedKey", "error")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateParentAuth return unexpected error code
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateParentAuth return unexpected type of error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, errors.New("unexpected type of error")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateParentInform return invalid duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, nil},
				"CreateParentInform":    {&model.ParentInform{}, &mysql.MySQLError{Number: mysqlcode.ER_DUP_ENTRY, Message: "InvalidMessage"}},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateParentInform return unexpected duplicate error
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, nil},
				"CreateParentInform":    {&model.ParentInform{}, mysqlerr.DuplicateEntry("UnexpectedKey", "duplicated")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateParentInform return unexpected error code
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, nil},
				"CreateParentInform":    {&model.ParentInform{}, &mysql.MySQLError{Number: mysqlcode.ER_BAD_NULL_ERROR, Message: "unexpected code"}},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // CreateParentInform return unexpected error code
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":               {},
				"GetParentAuthWithUUID": {&model.ParentAuth{}, gorm.ErrRecordNotFound},
				"CreateParentAuth":      {&model.ParentAuth{}, nil},
				"CreateParentInform":    {&model.ParentInform{}, errors.New("unexpected type of error")},
				"Rollback":              {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		},
	}

	for _, createNewParentTest := range tests {
		newMock, defaultHandler := generateVarForTest()

		createNewParentTest.ChangeEmptyValueToValidValue()
		createNewParentTest.ChangeEmptyReplaceValueToEmptyValue()
		createNewParentTest.OnExpectMethodsTo(newMock)

		req := new(proto.CreateNewParentRequest)
		createNewParentTest.SetRequestContextOf(req)
		ctx := createNewParentTest.GetMetadataContext()

		resp := new(proto.CreateNewParentResponse)
		_ = defaultHandler.CreateNewParent(ctx, req, resp)

		assert.Equalf(t, int(createNewParentTest.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", createNewParentTest, resp.Message)
		assert.Equalf(t, createNewParentTest.ExpectedCode, resp.Code, "code assertion error (test case: %v, message: %s)", createNewParentTest, resp.Message)
		assert.Regexpf(t, createNewParentTest.ExpectedStudentUUID, resp.CreatedParentUUID, "parent uuid assertion error (test case: %v, message: %s)", createNewParentTest, resp.Message)

		newMock.AssertExpectations(t)
	}
}

func Test_default_LoginAdminAuth(t *testing.T) {
	hashedByte, _ := bcrypt.GenerateFromPassword([]byte("testPW"), 1)

	tests := []test.LoginAdminAuthCase{
		{ // success case
			AdminID: "jinhong07191",
			AdminPW: "testPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetAdminAuthWithID": {&model.AdminAuth{
					UUID:    "admin-111111111111",
					AdminID: "jinhong07191",
					AdminPW: model.AdminPW(string(hashedByte)),
				}, nil},
				"Commit": {&gorm.DB{}},
			},
			ExpectedStatus:            http.StatusOK,
			ExpectedLoggedInAdminUUID: "admin-111111111111",
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
			AdminID: "jinhong07192",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":            {},
				"GetAdminAuthWithID": {&model.AdminAuth{}, gorm.ErrRecordNotFound},
				"Rollback":           {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.AdminIDNoExist,
		}, { // GetParentAuthWithID unexpected error
			AdminID: "jinhong07193",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx":            {},
				"GetAdminAuthWithID": {&model.AdminAuth{}, errors.New("unexpected error")},
				"Rollback":           {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusInternalServerError,
		}, { // incorrect Parent PW
			AdminID: "jinhong07194",
			AdminPW: "incorrectPW",
			ExpectedMethods: map[test.Method]test.Returns{
				"BeginTx": {},
				"GetAdminAuthWithID": {&model.AdminAuth{
					UUID:    "admin-111111111111",
					AdminID: "jinhong07194",
					AdminPW: model.AdminPW(string(hashedByte)),
				}, nil},
				"Rollback": {&gorm.DB{}},
			},
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   code.IncorrectAdminPWForLogin,
		},
	}

	for _, testCase := range tests {
		newMock, defaultHandler := generateVarForTest()

		testCase.ChangeEmptyValueToValidValue()
		testCase.ChangeEmptyReplaceValueToEmptyValue()
		testCase.OnExpectMethods(newMock)

		var req = new(proto.LoginAdminAuthRequest)
		testCase.SetRequestContextOf(req)
		ctx := testCase.GetMetadataContext()

		var resp = new(proto.LoginAdminAuthResponse)
		_ = defaultHandler.LoginAdminAuth(ctx, req, resp)

		assert.Equalf(t, int(testCase.ExpectedStatus), int(resp.Status), "status assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, int(testCase.ExpectedCode), int(resp.Code), "code assertion error (test case: %v, message: %s)", testCase, resp.Message)
		assert.Equalf(t, testCase.ExpectedLoggedInAdminUUID, resp.LoggedInAdminUUID, "logged in uuid assertion error (test case: %v, message: %s)", testCase, resp.Message)

		newMock.AssertExpectations(t)
	}
}
