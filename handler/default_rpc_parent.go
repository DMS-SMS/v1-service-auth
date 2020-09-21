package handler

import (
	"auth/model"
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

func (h _default) LoginParentAuth(ctx context.Context, req *proto.LoginParentAuthRequest, resp *proto.LoginParentAuthResponse) (err error) {
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

	spanForDB := opentracing.StartSpan("GetParentAuthWithID", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.GetParentAuthWithID(req.ParentID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", resultAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = CodeParentIDNoExist
			resp.Message = fmt.Sprintf(conflictErrorFormat, "parent id not exists")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(resultAuth.ParentPW), []byte(req.ParentPW))
	if err != nil {
		access.Rollback()
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = CodeIncorrectParentPWForLogin
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mismatched hash and password")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "hash compare error, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "succeed to login parent auth"
	resp.LoggedInParentUUID = string(resultAuth.UUID)

	return
}

func (h _default) ChangeParentPW(ctx context.Context, req *proto.ChangeParentPWRequest, resp *proto.ChangeParentPWResponse) (err error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	if !parentUUIDRegex.MatchString(req.ParentUUID) {
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "this API is for parents only")
		return
	}

	if req.UUID != req.ParentUUID {
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "not your auth, uuid: " + req.ParentUUID)
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

	spanForDB := opentracing.StartSpan("GetParentAuthWithUUID", opentracing.ChildOf(parentSpan))
	selectedAuth, err := access.GetParentAuthWithUUID(req.ParentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", selectedAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "not exist parent, uuid: " + req.ParentUUID)
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(selectedAuth.ParentPW), []byte(req.CurrentPW))
	if err != nil {
		access.Rollback()
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = CodeIncorrectParentPWForChange
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mismatched hash and password")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "hash compare error, err: " + err.Error())
		}
		return
	}

	spanForDB = opentracing.StartSpan("ChangeParentPW", opentracing.ChildOf(parentSpan))
	err = access.ChangeParentPW(string(selectedAuth.UUID), req.RevisionPW)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to update DB, err: " + err.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "parent pw change success"
	return
}

func (h _default) GetParentInformWithUUID(ctx context.Context, req *proto.GetParentInformWithUUIDRequest, resp *proto.GetParentInformWithUUIDResponse) (err error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	if !parentUUIDRegex.MatchString(req.UUID) {
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "this API is for parents only")
		return
	}

	if req.UUID != req.ParentUUID {
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "not your auth, uuid: " + req.ParentUUID)
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

	spanForDB := opentracing.StartSpan("GetParentInformWithUUID", opentracing.ChildOf(parentSpan))
	selectedAuth, err := access.GetParentInformWithUUID(req.ParentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", selectedAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "not exist parent, uuid: " + req.ParentUUID)
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.Name = string(selectedAuth.Name)
	resp.PhoneNumber = string(selectedAuth.PhoneNumber)

	resp.Status = http.StatusOK
	resp.Message = "get parent auth success"
	return
}

func (h _default) GetParentUUIDsWithInform(ctx context.Context, req *proto.GetParentUUIDsWithInformRequest, resp *proto.GetParentUUIDsWithInformResponse) (err error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case parentUUIDRegex.MatchString(req.UUID):
		break
	case adminUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "this API is for parents or admins only")
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

	spanForDB := opentracing.StartSpan("GetParentUUIDsWithInform", opentracing.ChildOf(parentSpan))
	selectedUUIDs, err := access.GetParentUUIDsWithInform(&model.ParentInform{
		Name:          model.Name(req.Name),
		PhoneNumber:   model.PhoneNumber(req.PhoneNumber),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedUUIDs", selectedUUIDs), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = CodeParentWithThatInformNoExist
			resp.Message = fmt.Sprintf(conflictErrorFormat, "no exist parent with that inform")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.ParentUUIDs = selectedUUIDs
	resp.Status = http.StatusOK
	resp.Message = "get parent uuids having that inform success"
	return
}
