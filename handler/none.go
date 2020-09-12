package handler

import (
	proto "auth/proto/golang/auth"
	"context"
)

type None struct {}

// About Admin RPC Service
func(n None) CreateNewAuth(context.Context, *proto.CreateNewAuthRequest, *proto.CreateNewAuthResponse) (err error) { return }
func(n None) CreateNewTeacher(context.Context, *proto.CreateNewTeacherRequest, *proto.CreateNewTeacherResponse) (err error) { return }
func(n None) CreateNewParent(context.Context, *proto.CreateNewParentRequest, *proto.CreateNewParentResponse) (err error) { return }

// About Student RPC Service
func(n None) LoginStudentAuth(context.Context, *proto.LoginStudentAuthRequest, *proto.LoginStudentAuthResponse) (err error) { return }
func(n None) ChangeStudentPW(context.Context, *proto.ChangeStudentPWRequest, *proto.ChangeStudentPWResponse) (err error) { return }
func(n None) GetStudentInformWithUUID(context.Context, *proto.GetStudentInformWithUUIDRequest, *proto.GetStudentInformWithUUIDResponse) (err error) { return }

// About Teacher RPC Service
func(n None) LoginTeacherAuth(context.Context, *proto.LoginTeacherAuthRequest, *proto.LoginTeacherAuthResponse) (err error) { return }
func(n None) ChangeTeacherPW(context.Context, *proto.ChangeTeacherPWRequest, *proto.ChangeTeacherPWResponse) (err error) { return }
func(n None) GetTeacherInformWithUUID(context.Context, *proto.GetTeacherInformWithUUIDRequest, *proto.GetTeacherInformWithUUIDResponse) (err error) { return }

// About Parent RPC Service
func(n None) LoginParentAuth(context.Context, *proto.LoginParentAuthRequest, *proto.LoginParentAuthResponse) (err error) { return }
func(n None) ChangeParentPW(context.Context, *proto.ChangeParentPWRequest, *proto.ChangeParentPWResponse) (err error) { return }
func(n None) GetParentInformWithUUID(context.Context, *proto.GetParentInformWithUUIDRequest, *proto.GetParentInformWithUUIDResponse) (err error) { return }
