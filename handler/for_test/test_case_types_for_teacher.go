package test

type LoginTeacherAuthCase struct {
	TeacherID, TeacherPW        string
	XRequestID                  string
	SpanContextString           string
	ExpectedMethods             map[Method]Returns
	ExpectedStatus              uint32
	ExpectedCode                int32
	ExpectedMessage             string
	ExpectedAccessToken			string
	ExpectedLoggedInTeacherUUID string
}
