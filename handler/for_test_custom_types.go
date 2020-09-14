package handler

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"fmt"
	"github.com/stretchr/testify/mock"
	"log"
)

type method string
type returns []interface{}

type createNewStudentTest struct {
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

func (test *createNewStudentTest) ChangeEmptyValueToValidValue() {
	//reflect.ValueOf(test).FieldByName().Type()
	if test.UUID == emptyString              { test.UUID = validAdminUUID }
	if test.StudentID == emptyString         { test.StudentID = validStudentID }
	if test.StudentPW == emptyString         { test.StudentPW = validStudentPW }
	if test.ParentUUID == emptyString        { test.ParentUUID = validParentUUID }
	if test.Grade == emptyUint32             { test.Grade = validGrade }
	if test.Class == emptyUint32             { test.Class = validClass }
	if test.StudentNumber == emptyUint32     { test.StudentNumber = validStudentNumber }
	if test.Name == emptyString              { test.Name = validName }
	if test.PhoneNumber == emptyString       { test.PhoneNumber = validPhoneNumber }
	if string(test.Image) == emptyString     { test.Image = validImageByteArr }
	if test.StudentUUID == emptyString       { test.StudentUUID = validStudentUUID }
	if test.XRequestID == emptyString        { test.XRequestID = validXRequestID }
	if test.SpanContextString == emptyString { test.SpanContextString = validSpanContextString }
}

func (test *createNewStudentTest) ChangeEmptyReplaceValueToEmptyValue() {
	if test.UUID == emptyReplaceValueForString               { test.UUID = "" }
	if test.StudentID == emptyReplaceValueForString          { test.StudentID = "" }
	if test.StudentPW == emptyReplaceValueForString          { test.StudentPW = "" }
	if test.ParentUUID == emptyReplaceValueForString         { test.ParentUUID = "" }
	if test.Grade == uint32(emptyReplaceValueForInt)         { test.Grade = 0 }
	if test.Class == uint32(emptyReplaceValueForInt)         { test.Class = 0 }
	if test.StudentNumber == uint32(emptyReplaceValueForInt) { test.StudentNumber = 0 }
	if test.Name == emptyReplaceValueForString               { test.Name = "" }
	if test.PhoneNumber == emptyReplaceValueForString        { test.PhoneNumber = "" }
	if string(test.Image) == emptyReplaceValueForString	     { test.Image = []byte{} }
	if test.StudentUUID == emptyReplaceValueForString        { test.StudentUUID = "" }
	if test.XRequestID == emptyReplaceValueForString         { test.XRequestID = "" }
	if test.SpanContextString == emptyReplaceValueForString  { test.SpanContextString = "" }
}

func (test *createNewStudentTest) OnExpectMethodsTo(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *createNewStudentTest) onMethod(mock *mock.Mock, method method, returns returns) {
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

func (test *createNewStudentTest) getStudentAuthModel() *model.StudentAuth {
	return &model.StudentAuth{
		UUID:       model.UUID(test.StudentUUID),
		StudentID:  model.StudentID(test.StudentID),
		StudentPW:  model.StudentPW(test.StudentPW),
		ParentUUID: model.ParentUUID(test.ParentUUID),
	}
}

func (test *createNewStudentTest) getStudentInformModel() *model.StudentInform {
	return &model.StudentInform{
		StudentUUID:   model.StudentUUID(test.StudentUUID),
		Grade:         model.Grade(int64(test.Grade)),
		Class:         model.Class(int64(test.Class)),
		StudentNumber: model.StudentNumber(int64(test.StudentNumber)),
		Name:          model.Name(test.Name),
		PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
		ProfileURI:    model.ProfileURI(fmt.Sprintf("example.com/profiles/%s", test.StudentUUID)),
	}
}

func (test *createNewStudentTest) SetRequestContextOf(req *proto.CreateNewStudentRequest) {
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