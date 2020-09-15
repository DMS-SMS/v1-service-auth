package test

type LoginStudentAuthCase struct {
	StudentID, StudentPW        string
	XRequestID                  string
	SpanContextString           string
	ExpectedMethods             map[Method]Returns
	ExpectedStatus              uint32
	ExpectedCode                int32
	ExpectedMessage             string
	ExpectedAccessToken			string
	ExpectedLoggedInStudentUUID string
}

func (test *LoginStudentAuthCase) ChangeEmptyValueToValidValue() {
	if test.StudentID == ""         { test.StudentID = validStudentID }
	if test.StudentPW == ""         { test.StudentPW = validStudentPW }
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
}

func (test *LoginStudentAuthCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.StudentID == EmptyReplaceValueForString         { test.StudentID = "" }
	if test.StudentPW == EmptyReplaceValueForString         { test.StudentPW = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
}
