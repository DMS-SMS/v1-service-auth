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
