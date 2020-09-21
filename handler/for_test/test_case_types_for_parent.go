package test

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"context"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/mock"
)

type LoginParentAuthCase struct {
	ParentID, ParentPW         string
	XRequestID                 string
	SpanContextString          string
	ExpectedMethods            map[Method]Returns
	ExpectedStatus             uint32
	ExpectedCode               int32
	ExpectedMessage            string
	ExpectedAccessToken        string
	ExpectedLoggedInParentUUID string
}

func (test *LoginParentAuthCase) ChangeEmptyValueToValidValue() {
	if test.ParentID == ""          { test.ParentID = validParentID }
	if test.ParentPW == ""          { test.ParentPW = validParentPW }
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
}

func (test *LoginParentAuthCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
}

func (test *LoginParentAuthCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *LoginParentAuthCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetParentAuthWithID":
		mock.On(string(method), test.ParentID).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *LoginParentAuthCase) SetRequestContextOf(req *proto.LoginParentAuthRequest) {
	req.ParentID = test.ParentID
	req.ParentPW = test.ParentPW
}

func (test *LoginParentAuthCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}

type ChangeParentPWCase struct {
	UUID, ParentUUID      string
	CurrentPW, RevisionPW string
	XRequestID            string
	SpanContextString     string
	ExpectedMethods       map[Method]Returns
	ExpectedStatus        uint32
	ExpectedCode          int32
	ExpectedMessage       string
}

func (test *ChangeParentPWCase) ChangeEmptyValueToValidValue() {
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
}

func (test *ChangeParentPWCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
}

func (test *ChangeParentPWCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *ChangeParentPWCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetParentAuthWithUUID":
		mock.On(string(method), test.ParentUUID).Return(returns...)
	case "ChangeParentPW":
		mock.On(string(method), test.ParentUUID, "").Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *ChangeParentPWCase) SetRequestContextOf(req *proto.ChangeParentPWRequest) {
	req.UUID = test.UUID
	req.ParentUUID = test.ParentUUID
	req.CurrentPW = test.CurrentPW
	req.RevisionPW = test.RevisionPW
}

func (test *ChangeParentPWCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}

type GetParentInformWithUUIDCase struct {
	UUID, ParentUUID  string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedMessage   string
	ExpectedInform    *model.ParentInform
}

func (test *GetParentInformWithUUIDCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
}

func (test *GetParentInformWithUUIDCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetParentInformWithUUIDCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetParentInformWithUUIDCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetParentInformWithUUID":
		mock.On(string(method), test.ParentUUID).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *GetParentInformWithUUIDCase) SetRequestContextOf(req *proto.GetParentInformWithUUIDRequest) {
	req.UUID = test.UUID
	req.ParentUUID = test.ParentUUID
}

func (test *GetParentInformWithUUIDCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}
