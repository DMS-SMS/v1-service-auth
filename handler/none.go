package handler

import (
	proto "auth/proto/golang/auth"
	"context"
)

type None struct {}

// About Admin RPC Service
func(n None) CreateStudentAuth(context.Context, *proto.CreateStudentAuthRequest, *proto.CreateStudentAuthResponse) (err error) { return }
func(n None) CreateTeacherAuth(context.Context, *proto.CreateTeacherAuthRequest, *proto.CreateTeacherAuthResponse) (err error) { return }
func(n None) CreateParentAuth(context.Context, *proto.CreateParentAuthRequest, *proto.CreateParentAuthResponse) (err error) { return }

// About Student RPC Service
func(n None) LoginStudentAuth(context.Context, *proto.LoginStudentAuthRequest, *proto.LoginStudentAuthResponse) (err error) { return }
func(n None) ChangeStudentAuthPw(context.Context, *proto.ChangeStudentAuthPwRequest, *proto.ChangeStudentAuthPwResponse) (err error) { return }
func(n None) GetStudentUserInform(context.Context, *proto.GetStudentUserInformRequest, *proto.GetStudentUserInformResponse) (err error) { return }

// About Teacher RPC Service
func(n None) LoginTeacherAuth(context.Context, *proto.LoginTeacherAuthRequest, *proto.LoginTeacherAuthResponse) (err error) { return }
func(n None) ChangeTeacherAuthPw(context.Context, *proto.ChangeTeacherAuthPwRequest, *proto.ChangeTeacherAuthPwResponse) (err error) { return }
func(n None) GetTeacherUserInform(context.Context, *proto.GetTeacherUserInformRequest, *proto.GetTeacherUserInformResponse) (err error) { return }

// About Parent RPC Service
func(n None) LoginParentAuth(context.Context, *proto.LoginParentAuthRequest, *proto.LoginParentAuthResponse) (err error) { return }
func(n None) ChangeParentAuthPw(context.Context, *proto.ChangeParentAuthPwRequest, *proto.ChangeParentAuthPwResponse) (err error) { return }
func(n None) GetParentUserInform(context.Context, *proto.GetParentUserInformRequest, *proto.GetParentUserInformResponse) (err error) { return }
