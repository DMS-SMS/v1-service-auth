package handler

import (
	"github.com/stretchr/testify/mock"
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

func (test *createNewStudentTest) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *createNewStudentTest) onMethod(mock *mock.Mock, method method, returns returns) {
	switch method {

	}
}
