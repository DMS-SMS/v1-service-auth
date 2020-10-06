package test

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"context"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/stretchr/testify/mock"
)

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

func (test *LoginStudentAuthCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *LoginStudentAuthCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetStudentAuthWithID":
		mock.On(string(method), test.StudentID).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *LoginStudentAuthCase) SetRequestContextOf(req *proto.LoginStudentAuthRequest) {
	req.StudentID = test.StudentID
	req.StudentPW = test.StudentPW
}

func (test *LoginStudentAuthCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}

type ChangeStudentPWCase struct {
	UUID, StudentUUID     string
	CurrentPW, RevisionPW string
	XRequestID            string
	SpanContextString     string
	ExpectedMethods       map[Method]Returns
	ExpectedStatus        uint32
	ExpectedCode          int32
	ExpectedMessage       string
}

func (test *ChangeStudentPWCase) ChangeEmptyValueToValidValue() {
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
}

func (test *ChangeStudentPWCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
}

func (test *ChangeStudentPWCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *ChangeStudentPWCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetStudentAuthWithUUID": // 추가 구현 필요
		mock.On(string(method), test.StudentUUID).Return(returns...)
	case "ChangeStudentPW":
		mock.On(string(method), test.StudentUUID, "").Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *ChangeStudentPWCase) SetRequestContextOf(req *proto.ChangeStudentPWRequest) {
	req.UUID = test.UUID
	req.StudentUUID = test.StudentUUID
	req.CurrentPW = test.CurrentPW
	req.RevisionPW = test.RevisionPW
}

func (test *ChangeStudentPWCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}

type GetStudentInformWithUUIDCase struct {
	UUID, StudentUUID string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedMessage   string
	ExpectedInform    *model.StudentInform
}

func (test *GetStudentInformWithUUIDCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
}

func (test *GetStudentInformWithUUIDCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetStudentInformWithUUIDCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetStudentInformWithUUIDCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetStudentInformWithUUID":
		mock.On(string(method), test.StudentUUID).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *GetStudentInformWithUUIDCase) SetRequestContextOf(req *proto.GetStudentInformWithUUIDRequest) {
	req.UUID = test.UUID
	req.StudentUUID = test.StudentUUID
}

func (test *GetStudentInformWithUUIDCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}

type GetStudentUUIDsWithInformCase struct {
	UUID                 string
	Grade, Class         int64
	StudentNumber        int64
	Name, PhoneNumber    string
	ImageURI             string
	XRequestID           string
	SpanContextString    string
	ExpectedMethods      map[Method]Returns
	ExpectedStatus       uint32
	ExpectedCode         int32
	ExpectedMessage      string
	ExpectedStudentUUIDs []string
}

func (test *GetStudentUUIDsWithInformCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
}

func (test *GetStudentUUIDsWithInformCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetStudentUUIDsWithInformCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetStudentUUIDsWithInformCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetStudentUUIDsWithInform":
		mock.On(string(method), &model.StudentInform{
			Grade:         model.Grade(test.Grade),
			Class:         model.Class(test.Class),
			StudentNumber: model.StudentNumber(test.StudentNumber),
			Name:          model.Name(test.Name),
			PhoneNumber:   model.PhoneNumber(test.PhoneNumber),
			ProfileURI:    model.ProfileURI(test.ImageURI),
		}).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *GetStudentUUIDsWithInformCase) SetRequestContextOf(req *proto.GetStudentUUIDsWithInformRequest) {
	req.UUID = test.UUID
	req.Grade = uint32(test.Grade)
	req.Group = uint32(test.Class)
	req.StudentNumber = uint32(test.StudentNumber)
	req.Name = test.Name
	req.PhoneNumber = test.PhoneNumber
	req.ImageURI = test.ImageURI
}

func (test *GetStudentUUIDsWithInformCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}

type GetStudentInformsWithUUIDsCase struct {
	UUID              string
	StudentUUIDs      []string
	XRequestID        string
	SpanContextString string
	ExpectedMethods   map[Method]Returns
	ExpectedStatus    uint32
	ExpectedCode      int32
	ExpectedMessage   string
	ExpectedInforms   []*proto.StudentInform
}

func (test *GetStudentInformsWithUUIDsCase) ChangeEmptyValueToValidValue() {
	if test.XRequestID == ""        { test.XRequestID = validXRequestID }
	if test.SpanContextString == "" { test.SpanContextString = validSpanContextString }
}

func (test *GetStudentInformsWithUUIDsCase) ChangeEmptyReplaceValueToEmptyValue() {
	if test.XRequestID == EmptyReplaceValueForString        { test.XRequestID = "" }
	if test.SpanContextString == EmptyReplaceValueForString { test.SpanContextString = "" }
}

func (test *GetStudentInformsWithUUIDsCase) OnExpectMethods(mock *mock.Mock) {
	for method, returns := range test.ExpectedMethods {
		test.onMethod(mock, method, returns)
	}
}

func (test *GetStudentInformsWithUUIDsCase) onMethod(mock *mock.Mock, method Method, returns Returns) {
	switch method {
	case "BeginTx":
		mock.On(string(method)).Return(returns...)
	case "GetStudentInformsWithUUIDs":
		mock.On(string(method), test.StudentUUIDs).Return(returns...)
	case "Commit":
		mock.On(string(method)).Return(returns...)
	case "Rollback":
		mock.On(string(method)).Return(returns...)
	}
}

func (test *GetStudentInformsWithUUIDsCase) SetRequestContextOf(req *proto.GetStudentInformsWithUUIDsRequest) {
	req.UUID = test.UUID
	req.StudentUUIDs = test.StudentUUIDs
}

func (test *GetStudentInformsWithUUIDsCase) GetMetadataContext() (ctx context.Context) {
	ctx = context.Background()

	ctx = metadata.Set(ctx, "X-Request-Id", test.XRequestID)
	ctx = metadata.Set(ctx, "Span-Context", test.SpanContextString)

	return
}
