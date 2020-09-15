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

type method string
type returns []interface{}

type CreateNewStudentTest struct {
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
	ExpectedMethods      map[method]returns
	ExpectedStatus       uint32
	ExpectedCode         int32
	ExpectedMessage      string
	ExpectedStudentUUID  string
}

func (test *CreateNewStudentTest) ChangeEmptyValueToValidValue() {
	//reflect.ValueOf(test).FieldByName().Type()
	if test.UUID == emptyString              { test.UUID = validAdminUUID }
	if test.StudentID == emptyString         { test.StudentID = validStudentID }
	if test.StudentPW == emptyString         { test.StudentPW = validStudentPW }
	if test.ParentUUID == emptyString        { test.ParentUUID = validParentUUID() }
	if test.Grade == emptyUint32             { test.Grade = validGrade }
	if test.Class == emptyUint32             { test.Class = validClass }
	if test.StudentNumber == emptyUint32     { test.StudentNumber = validStudentNumber }
	if test.Name == emptyString              { test.Name = validName }
	if test.PhoneNumber == emptyString       { test.PhoneNumber = validPhoneNumber }
	if string(test.Image) == emptyString     { test.Image = validImageByteArr }
	if test.StudentUUID == emptyString       { test.StudentUUID = validStudentUUID() }
	if test.XRequestID == emptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == emptyString { test.SpanContextString = validSpanContextString }
}

func (test *CreateNewStudentTest) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == emptyReplaceValueForString               { test.UUID = "" }
	if test.StudentID == emptyReplaceValueForString          { test.StudentID = "" }
	if test.StudentPW == emptyReplaceValueForString          { test.StudentPW = "" }
	if test.ParentUUID == emptyReplaceValueForString         { test.ParentUUID = "" }
	if test.Grade == emptyReplaceValueForUint32              { test.Grade = 0 }
	if test.Class == emptyReplaceValueForUint32              { test.Class = 0 }
	if test.StudentNumber == emptyReplaceValueForUint32      { test.StudentNumber = 0 }
	if test.Name == emptyReplaceValueForString               { test.Name = "" }
	if test.PhoneNumber == emptyReplaceValueForString        { test.PhoneNumber = "" }
	if string(test.Image) == emptyReplaceValueForString	     { test.Image = []byte{} }
	if test.StudentUUID == emptyReplaceValueForString        { test.StudentUUID = "" }
	if test.XRequestID == emptyReplaceValueForString         { test.XRequestID = "" }
	if test.SpanContextString == emptyReplaceValueForString  { test.SpanContextString = "" }
}

func (test *CreateNewStudentTest) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *CreateNewStudentTest) onMethod(mock *mock.Mock, method method, returns returns) {
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

	case "CheckIfStudentAuthExists":
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

func (test *CreateNewStudentTest) getStudentAuthModel() *model.StudentAuth {
	return &model.StudentAuth{
		UUID:       model.UUID(test.StudentUUID),
		StudentID:  model.StudentID(test.StudentID),
		StudentPW:  model.StudentPW(test.StudentPW),
		ParentUUID: model.ParentUUID(test.ParentUUID),
	}
}

func (test *CreateNewStudentTest) getStudentInformModel() *model.StudentInform {
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

func (test *CreateNewStudentTest) SetRequestContextOf(req *proto.CreateNewStudentRequest) {
	req.UUID = test.UUID
	req.StudentID = test.StudentID
	req.StudentPW = test.StudentPW
	req.ParentUUID = test.ParentUUID
	req.Grade = test.Grade
	req.Class = test.Class
	req.StudentNumber = test.StudentNumber
	req.Name = test.Name
	req.PhoneNumber = test.PhoneNumber
	req.Image = test.Image
}

func (test *CreateNewStudentTest) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	ctx = metadata.Set(ctx, "StudentUUID", test.StudentUUID)
	return
}

type CreateNewTeacherTest struct {
	UUID                 string
	TeacherID, TeacherPW string
	Grade, Class         uint32
	Name, PhoneNumber    string
	TeacherUUID          string
	XRequestID           string
	SpanContextString    string
	ExpectedMethods      map[method]returns
	ExpectedStatus       uint32
	ExpectedCode         int32
	ExpectedMessage      string
	ExpectedStudentUUID  string
}

func (test *CreateNewTeacherTest) ChangeEmptyValueToValidValue() {
	//reflect.ValueOf(test).FieldByName().Type()
	if test.UUID == emptyString              { test.UUID = validAdminUUID }
	if test.TeacherID == emptyString         { test.TeacherID = validTeacherID }
	if test.TeacherPW == emptyString         { test.TeacherPW = validTeacherPW }
	if test.Grade == emptyUint32             { test.Grade = validGrade }
	if test.Class == emptyUint32             { test.Class = validClass }
	if test.Name == emptyString              { test.Name = validName }
	if test.PhoneNumber == emptyString       { test.PhoneNumber = validPhoneNumber }
	if test.TeacherUUID == emptyString       { test.TeacherUUID = validTeacherUUID() }
	if test.XRequestID == emptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == emptyString { test.SpanContextString = validSpanContextString }
}

func (test *CreateNewTeacherTest) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == emptyReplaceValueForString               { test.UUID = "" }
	if test.TeacherID == emptyReplaceValueForString          { test.TeacherID = "" }
	if test.TeacherPW == emptyReplaceValueForString          { test.TeacherPW = "" }
	if test.Grade == emptyReplaceValueForUint32              { test.Grade = 0 }
	if test.Class == emptyReplaceValueForUint32              { test.Class = 0 }
	if test.Name == emptyReplaceValueForString               { test.Name = "" }
	if test.PhoneNumber == emptyReplaceValueForString        { test.PhoneNumber = "" }
	if test.TeacherUUID == emptyReplaceValueForString        { test.TeacherUUID = "" }
	if test.XRequestID == emptyReplaceValueForString         { test.XRequestID = "" }
	if test.SpanContextString == emptyReplaceValueForString  { test.SpanContextString = "" }
}

func (test *CreateNewTeacherTest) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *CreateNewTeacherTest) onMethod(mock *mock.Mock, method method, returns returns) {
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

	case "CheckIfTeacherAuthExists":
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

func (test *CreateNewTeacherTest) getTeacherAuthModel() *model.TeacherAuth {
	return &model.TeacherAuth{
		UUID:       model.UUID(test.TeacherUUID),
		TeacherID:  model.TeacherID(test.TeacherID),
		TeacherPW:  model.TeacherPW(test.TeacherPW),
	}
}

func (test *CreateNewTeacherTest) getTeacherInformModel() *model.TeacherInform {
	return &model.TeacherInform{
		TeacherUUID:   model.TeacherUUID(test.TeacherUUID),
		Grade:         model.Grade(int64(test.Grade)),
		Class:         model.Class(int64(test.Class)),
		Name:          model.Name(test.Name),
		PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
	}
}

func (test *CreateNewTeacherTest) SetRequestContextOf(req *proto.CreateNewTeacherRequest) {
	req.UUID = test.UUID
	req.TeacherID = test.TeacherID
	req.TeacherPW = test.TeacherPW
	req.Grade = test.Grade
	req.Class = test.Class
	req.Name = test.Name
	req.PhoneNumber = test.PhoneNumber
}

func (test *CreateNewTeacherTest) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	ctx = metadata.Set(ctx, "TeacherUUID", test.TeacherUUID)
	return
}

type CreateNewParentTest struct {
	UUID                string
	ParentID, ParentPW  string
	Name, PhoneNumber   string
	ParentUUID          string
	XRequestID          string
	SpanContextString   string
	ExpectedMethods     map[method]returns
	ExpectedStatus      uint32
	ExpectedCode        int32
	ExpectedMessage     string
	ExpectedStudentUUID string
}

func (test *CreateNewParentTest) ChangeEmptyValueToValidValue() {
	//reflect.ValueOf(test).FieldByName().Type()
	if test.UUID == emptyString              { test.UUID = validAdminUUID }
	if test.ParentID == emptyString          { test.ParentID = validParentID }
	if test.ParentPW == emptyString          { test.ParentPW = validParentPW }
	if test.Name == emptyString              { test.Name = validName }
	if test.PhoneNumber == emptyString       { test.PhoneNumber = validPhoneNumber }
	if test.ParentUUID == emptyString        { test.ParentUUID = validParentUUID() }
	if test.XRequestID == emptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == emptyString { test.SpanContextString = validSpanContextString }
}

func (test *CreateNewParentTest) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == emptyReplaceValueForString               { test.UUID = "" }
	if test.ParentID == emptyReplaceValueForString           { test.ParentID = "" }
	if test.ParentPW == emptyReplaceValueForString           { test.ParentPW = "" }
	if test.Name == emptyReplaceValueForString               { test.Name = "" }
	if test.PhoneNumber == emptyReplaceValueForString        { test.PhoneNumber = "" }
	if test.ParentUUID == emptyReplaceValueForString         { test.ParentUUID = "" }
	if test.XRequestID == emptyReplaceValueForString         { test.XRequestID = "" }
	if test.SpanContextString == emptyReplaceValueForString  { test.SpanContextString = "" }
}

func (test *CreateNewParentTest) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *CreateNewParentTest) onMethod(mock *mock.Mock, method method, returns returns) {
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

	case "CheckIfParentAuthExists":
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

func (test *CreateNewParentTest) getParentAuthModel() *model.ParentAuth {
	return &model.ParentAuth{
		UUID:     model.UUID(test.ParentUUID),
		ParentID: model.ParentID(test.ParentID),
		ParentPW: model.ParentPW(test.ParentPW),
	}
}

func (test *CreateNewParentTest) getParentInformModel() *model.ParentInform {
	return &model.ParentInform{
		ParentUUID:  model.ParentUUID(test.ParentUUID),
		Name:        model.Name(test.Name),
		PhoneNumber: model.PhoneNumber(test.PhoneNumber),
	}
}

func (test *CreateNewParentTest) SetRequestContextOf(req *proto.CreateNewParentRequest) {
	req.UUID = test.UUID
	req.ParentID = test.ParentID
	req.ParentPW = test.ParentPW
	req.Name = test.Name
	req.PhoneNumber = test.PhoneNumber
}

func (test *CreateNewParentTest) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()
	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)
	ctx = metadata.Set(ctx, "ParentUUID", test.ParentUUID)
	return
}
