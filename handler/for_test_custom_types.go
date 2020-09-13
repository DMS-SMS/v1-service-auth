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