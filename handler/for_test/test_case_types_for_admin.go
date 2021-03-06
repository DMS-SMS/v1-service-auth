package test

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/mock"
	"log"
)

type Method string
type Returns []interface{}

type CreateNewStudentCase struct {
	UUID                 string
	StudentID, StudentPW string
	ParentUUID           string
	Grade, Class         uint32
	StudentNumber        uint32
	Name, PhoneNumber    string
	Image                []byte
	StudentUUID          string
	XRequestID           string
	SpanContextString    string
	ExpectedMethods      map[Method]Returns
	ExpectedStatus       uint32
	ExpectedCode         int32
	ExpectedMessage      string
	ExpectedStudentUUID  string
}

func (test *CreateNewStudentCase) ChangeEmptyValueToValidValue() {
	//reflect.ValueOf(test).FieldByName().Type()
	if test.UUID == EmptyString              { test.UUID = validAdminUUID }
	if test.StudentID == EmptyString         { test.StudentID = validStudentID }
	if test.StudentPW == EmptyString         { test.StudentPW = validStudentPW }
	if test.ParentUUID == EmptyString        { test.ParentUUID = validParentUUID() }
	if test.Grade == EmptyUint32             { test.Grade = validGrade }
	if test.Class == EmptyUint32             { test.Class = validClass }
	if test.StudentNumber == EmptyUint32     { test.StudentNumber = validStudentNumber }
	if test.Name == EmptyString              { test.Name = validName }
	if test.PhoneNumber == EmptyString       { test.PhoneNumber = validPhoneNumber }
	if string(test.Image) == EmptyString     { test.Image = validImageByteArr }
	if test.StudentUUID == EmptyString       { test.StudentUUID = validStudentUUID() }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *CreateNewStudentCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString               { test.UUID = "" }
	if test.StudentID == EmptyReplaceValueForString          { test.StudentID = "" }
	if test.StudentPW == EmptyReplaceValueForString          { test.StudentPW = "" }
	if test.ParentUUID == EmptyReplaceValueForString         { test.ParentUUID = "" }
	if test.Grade == EmptyReplaceValueForUint32              { test.Grade = 0 }
	if test.Class == EmptyReplaceValueForUint32              { test.Class = 0 }
	if test.StudentNumber == EmptyReplaceValueForUint32      { test.StudentNumber = 0 }
	if test.Name == EmptyReplaceValueForString               { test.Name = "" }
	if test.PhoneNumber == EmptyReplaceValueForString        { test.PhoneNumber = "" }
	if string(test.Image) == EmptyReplaceValueForString	     { test.Image = []byte{} }
	if test.StudentUUID == EmptyReplaceValueForString        { test.StudentUUID = "" }
	if test.XRequestID == EmptyReplaceValueForString         { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString  { test.SpanContextString = "" }
}

func (test *CreateNewStudentCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *CreateNewStudentCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "CreateStudentAuth":
		const studentAuthIndex = 0
		const errorIndex = 1
		if _, ok := returns[studentAuthIndex].(*model.StudentAuth); ok && returns[errorIndex] == nil {
			modelToReturn := test.getStudentAuthModel()
			modelToReturn.Model = createGormModelOnCurrentTime()
			returns[studentAuthIndex] = modelToReturn
		}
		mock.On(string(method), test.getStudentAuthModel()).Return(returns...)

	case "CreateStudentInform":
		const studentInformIndex = 0
		const errorIndex = 1
		if _, ok := returns[studentInformIndex].(*model.StudentInform); ok && returns[errorIndex] == nil {
			modelToReturn := test.getStudentInformModel()
			modelToReturn.Model = createGormModelOnCurrentTime()
			returns[studentInformIndex] = modelToReturn
		}
		mock.On(string(method), test.getStudentInformModel()).Return(returns...)

	case "GetStudentAuthWithUUID":
		mock.On(string(method), test.StudentUUID).Return(returns...)

	case "BeginTx":
		mock.On(string(method)).Return(returns...)

	case "Commit":
		mock.On(string(method)).Return(returns...)

	case "Rollback":
		mock.On(string(method)).Return(returns...)

	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *CreateNewStudentCase) getStudentAuthModel() *model.StudentAuth {
	return &model.StudentAuth{
		UUID:       model.UUID(test.StudentUUID),
		StudentID:  model.StudentID(test.StudentID),
		StudentPW:  model.StudentPW(test.StudentPW),
		ParentUUID: model.ParentUUID(test.ParentUUID),
	}
}

func (test *CreateNewStudentCase) getStudentInformModel() *model.StudentInform {
	return &model.StudentInform{
		StudentUUID:   model.StudentUUID(test.StudentUUID),
		Grade:         model.Grade(int64(test.Grade)),
		Class:         model.Class(int64(test.Class)),
		StudentNumber: model.StudentNumber(int64(test.StudentNumber)),
		Name:          model.Name(test.Name),
		PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
		ProfileURI:    model.ProfileURI(fmt.Sprintf("profiles/%s", test.StudentUUID)),
	}
}

func (test *CreateNewStudentCase) SetRequestContextOf(req *proto.CreateNewStudentRequest) {
	req.UUID = test.UUID
	req.StudentID = test.StudentID
	req.StudentPW = test.StudentPW
	req.ParentUUID = test.ParentUUID
	req.Grade = test.Grade
	req.Group = test.Class
	req.StudentNumber = test.StudentNumber
	req.Name = test.Name
	req.PhoneNumber = test.PhoneNumber
	req.Image = test.Image
}

func (test *CreateNewStudentCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	ctx = metadata.Set(ctx, "StudentUUID", test.StudentUUID)
	return
}

type CreateNewTeacherCase struct {
	UUID                 string
	TeacherID, TeacherPW string
	Grade, Class         uint32
	Name, PhoneNumber    string
	TeacherUUID          string
	XRequestID           string
	SpanContextString    string
	ExpectedMethods      map[Method]Returns
	ExpectedStatus       uint32
	ExpectedCode         int32
	ExpectedMessage      string
	ExpectedStudentUUID  string
}

func (test *CreateNewTeacherCase) ChangeEmptyValueToValidValue() {
	//reflect.ValueOf(test).FieldByName().Type()
	if test.UUID == EmptyString              { test.UUID = validAdminUUID }
	if test.TeacherID == EmptyString         { test.TeacherID = validTeacherID }
	if test.TeacherPW == EmptyString         { test.TeacherPW = validTeacherPW }
	if test.Grade == EmptyUint32             { test.Grade = validGrade }
	if test.Class == EmptyUint32             { test.Class = validClass }
	if test.Name == EmptyString              { test.Name = validName }
	if test.PhoneNumber == EmptyString       { test.PhoneNumber = validPhoneNumber }
	if test.TeacherUUID == EmptyString       { test.TeacherUUID = validTeacherUUID() }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *CreateNewTeacherCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString               { test.UUID = "" }
	if test.TeacherID == EmptyReplaceValueForString          { test.TeacherID = "" }
	if test.TeacherPW == EmptyReplaceValueForString          { test.TeacherPW = "" }
	if test.Grade == EmptyReplaceValueForUint32              { test.Grade = 0 }
	if test.Class == EmptyReplaceValueForUint32              { test.Class = 0 }
	if test.Name == EmptyReplaceValueForString               { test.Name = "" }
	if test.PhoneNumber == EmptyReplaceValueForString        { test.PhoneNumber = "" }
	if test.TeacherUUID == EmptyReplaceValueForString        { test.TeacherUUID = "" }
	if test.XRequestID == EmptyReplaceValueForString         { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString  { test.SpanContextString = "" }
}

func (test *CreateNewTeacherCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *CreateNewTeacherCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "CreateTeacherAuth":
		const indexTeacherAuth = 0
		const indexError = 1
		if _, ok := returns[indexTeacherAuth].(*model.TeacherAuth); ok && returns[indexError] == nil {
			modelToReturn := test.getTeacherAuthModel()
			modelToReturn.Model = createGormModelOnCurrentTime()
			returns[indexTeacherAuth] = modelToReturn
		}
		mock.On(string(method), test.getTeacherAuthModel()).Return(returns...)

	case "CreateTeacherInform":
		const indexTeacherInform = 0
		const indexError = 1
		if _, ok := returns[indexTeacherInform].(*model.TeacherInform); ok && returns[indexError] == nil {
			modelToReturn := test.getTeacherInformModel()
			modelToReturn.Model = createGormModelOnCurrentTime()
			returns[indexTeacherInform] = modelToReturn
		}
		mock.On(string(method), test.getTeacherInformModel()).Return(returns...)

	case "GetTeacherAuthWithUUID":
		mock.On(string(method), test.TeacherUUID).Return(returns...)

	case "BeginTx":
		mock.On(string(method)).Return(returns...)

	case "Commit":
		mock.On(string(method)).Return(returns...)

	case "Rollback":
		mock.On(string(method)).Return(returns...)

	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *CreateNewTeacherCase) getTeacherAuthModel() *model.TeacherAuth {
	return &model.TeacherAuth{
		UUID:       model.UUID(test.TeacherUUID),
		TeacherID:  model.TeacherID(test.TeacherID),
		TeacherPW:  model.TeacherPW(test.TeacherPW),
	}
}

func (test *CreateNewTeacherCase) getTeacherInformModel() *model.TeacherInform {
	return &model.TeacherInform{
		TeacherUUID:   model.TeacherUUID(test.TeacherUUID),
		Grade:         model.Grade(int64(test.Grade)),
		Class:         model.Class(int64(test.Class)),
		Name:          model.Name(test.Name),
		PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
	}
}

func (test *CreateNewTeacherCase) SetRequestContextOf(req *proto.CreateNewTeacherRequest) {
	req.UUID = test.UUID
	req.TeacherID = test.TeacherID
	req.TeacherPW = test.TeacherPW
	req.Grade = test.Grade
	req.Group = test.Class
	req.Name = test.Name
	req.PhoneNumber = test.PhoneNumber
}

func (test *CreateNewTeacherCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	ctx = metadata.Set(ctx, "TeacherUUID", test.TeacherUUID)
	return
}

type CreateNewParentCase struct {
	UUID                string
	ParentID, ParentPW  string
	Name, PhoneNumber   string
	ParentUUID          string
	XRequestID          string
	SpanContextString   string
	ExpectedMethods     map[Method]Returns
	ExpectedStatus      uint32
	ExpectedCode        int32
	ExpectedMessage     string
	ExpectedStudentUUID string
}

func (test *CreateNewParentCase) ChangeEmptyValueToValidValue() {
	//reflect.ValueOf(test).FieldByName().Type()
	if test.UUID == EmptyString              { test.UUID = validAdminUUID }
	if test.ParentID == EmptyString          { test.ParentID = validParentID }
	if test.ParentPW == EmptyString          { test.ParentPW = validParentPW }
	if test.Name == EmptyString              { test.Name = validName }
	if test.PhoneNumber == EmptyString       { test.PhoneNumber = validPhoneNumber }
	if test.ParentUUID == EmptyString        { test.ParentUUID = validParentUUID() }
	if test.XRequestID == EmptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == EmptyString { test.SpanContextString = validSpanContextString }
}

func (test *CreateNewParentCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == EmptyReplaceValueForString               { test.UUID = "" }
	if test.ParentID == EmptyReplaceValueForString           { test.ParentID = "" }
	if test.ParentPW == EmptyReplaceValueForString           { test.ParentPW = "" }
	if test.Name == EmptyReplaceValueForString               { test.Name = "" }
	if test.PhoneNumber == EmptyReplaceValueForString        { test.PhoneNumber = "" }
	if test.ParentUUID == EmptyReplaceValueForString         { test.ParentUUID = "" }
	if test.XRequestID == EmptyReplaceValueForString         { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString  { test.SpanContextString = "" }
}

func (test *CreateNewParentCase) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *CreateNewParentCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "CreateParentAuth":
		const indexParentAuth = 0
		const indexError = 1
		if _, ok := returns[indexParentAuth].(*model.ParentAuth); ok && returns[indexError] == nil {
			modelToReturn := test.getParentAuthModel()
			modelToReturn.Model = createGormModelOnCurrentTime()
			returns[indexParentAuth] = modelToReturn
		}
		mock.On(string(method), test.getParentAuthModel()).Return(returns...)

	case "CreateParentInform":
		const indexParentInform = 0
		const indexError = 1
		if _, ok := returns[indexParentInform].(*model.ParentInform); ok && returns[indexError] == nil {
			modelToReturn := test.getParentInformModel()
			modelToReturn.Model = createGormModelOnCurrentTime()
			returns[indexParentInform] = modelToReturn
		}
		mock.On(string(method), test.getParentInformModel()).Return(returns...)

	case "GetParentAuthWithUUID":
		mock.On(string(method), test.ParentUUID).Return(returns...)

	case "BeginTx":
		mock.On(string(method)).Return(returns...)

	case "Commit":
		mock.On(string(method)).Return(returns...)

	case "Rollback":
		mock.On(string(method)).Return(returns...)

	default:
		log.Fatalf("this method cannot be registered, method name: %s", method)
	}
}

func (test *CreateNewParentCase) getParentAuthModel() *model.ParentAuth {
	return &model.ParentAuth{
		UUID:     model.UUID(test.ParentUUID),
		ParentID: model.ParentID(test.ParentID),
		ParentPW: model.ParentPW(test.ParentPW),
	}
}

func (test *CreateNewParentCase) getParentInformModel() *model.ParentInform {
	return &model.ParentInform{
		ParentUUID:  model.ParentUUID(test.ParentUUID),
		Name:        model.Name(test.Name),
		PhoneNumber: model.PhoneNumber(test.PhoneNumber),
	}
}

func (test *CreateNewParentCase) SetRequestContextOf(req *proto.CreateNewParentRequest) {
	req.UUID = test.UUID
	req.ParentID = test.ParentID
	req.ParentPW = test.ParentPW
	req.Name = test.Name
	req.PhoneNumber = test.PhoneNumber
}

func (test *CreateNewParentCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	ctx = metadata.Set(ctx, "ParentUUID", test.ParentUUID)
	return
}


type LoginAdminAuthCase struct {
	AdminID, AdminPW          string
	XRequestID                string
	SpanContextString         string
	ExpectedMethods           map[Method]Returns
	ExpectedStatus            uint32
	ExpectedCode              int32
	ExpectedMessage           string
	ExpectedAccessToken       string
	ExpectedLoggedInAdminUUID string
}

func (test *LoginAdminAuthCase) ChangeEmptyValueToValidValue() {
	if test.AdminID == ""           { test.AdminID = validAdminID }
	if test.AdminPW == ""           { test.AdminPW = validAdminPW }
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
}

func (test *LoginAdminAuthCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.AdminID == EmptyReplaceValueForString           { test.AdminID = "" }
	if test.AdminPW == EmptyReplaceValueForString           { test.AdminPW = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
}

func (test *LoginAdminAuthCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *LoginAdminAuthCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetAdminAuthWithID":
		mock.On(string(method), test.AdminID).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *LoginAdminAuthCase) SetRequestContextOf(req *proto.LoginAdminAuthRequest) {
	req.AdminID = test.AdminID
	req.AdminPW = test.AdminPW
}

func (test *LoginAdminAuthCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}