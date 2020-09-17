package test

import (
	proto "auth/proto/golang/auth"
	"context"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/mock"
)

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

func (test *LoginTeacherAuthCase) ChangeEmptyValueToValidValue() {
	if test.TeacherID == ""         { test.TeacherID = validTeacherID }
	if test.TeacherPW == ""         { test.TeacherPW = validTeacherPW }
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
}

func (test *LoginTeacherAuthCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
}

func (test *LoginTeacherAuthCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *LoginTeacherAuthCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetTeacherAuthWithID":
		mock.On(string(method), test.TeacherID).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *LoginTeacherAuthCase) SetRequestContextOf(req *proto.LoginTeacherAuthRequest) {
	req.TeacherID = test.TeacherID
	req.TeacherPW = test.TeacherPW
}

func (test *LoginTeacherAuthCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}