package handler

import (
	proto "auth/proto/golang/auth"
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func(h _default) LoginStudentAuth(ctx context.Context, req *proto.LoginStudentAuthRequest, resp *proto.LoginStudentAuthResponse) (err error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	reqID := ctx.Value("X-Request-Id").(string)
	parentSpan := ctx.Value("Span-Context").(jaeger.SpanContext)

	access, err := h.accessManage.BeginTx()
	if err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "tx begin fail, err: " + err.Error())
		return
	}

	spanForDB := opentracing.StartSpan("GetStudentAuthWithID", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.GetStudentAuthWithID(req.StudentID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", resultAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = CodeStudentIDNoExist
			resp.Message = fmt.Sprintf(conflictErrorFormat, "student id not exists")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(resultAuth.StudentPW), []byte(req.StudentPW))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = CodeIncorrectStudentPW
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mismatched hash and password")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "hash compare error, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "succeed to login student auth"
	resp.LoggedInStudentUUID = string(resultAuth.UUID)

	return
}
