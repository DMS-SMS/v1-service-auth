package handler

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
	ExpectMethod         map[method]returns
	ExpectedStatus       uint32
	ExpectedCode         int32
	ExpectedMessage      string
}

func (test *createNewStudentTest) ChangeEmptyValueToValidValue() {
	if test.UUID == emptyString          { test.UUID = validAdminUUID }
	if test.StudentID == emptyString     { test.StudentID = validStudentID }
	if test.StudentPW == emptyString     { test.StudentPW = validStudentPW }
	if test.ParentUUID == emptyString    { test.ParentUUID = validParentUUID }
	if test.Grade == emptyUint32         { test.Grade = validGrade }
	if test.Class == emptyUint32         { test.Class = validClass }
	if test.StudentNumber == emptyUint32 { test.StudentNumber = validStudentNumber }
	if test.Name == emptyString          { test.Name = validName }
	if test.PhoneNumber == emptyString   { test.PhoneNumber = validPhoneNumber }
	if string(test.Image) == emptyString { test.Image = validImageByteArr }
}
