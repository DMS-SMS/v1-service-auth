package handler

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	code "auth/utils/code/golang"
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"reflect"
)

func (h _default) LoginParentAuth(ctx context.Context, req *proto.LoginParentAuthRequest, resp *proto.LoginParentAuthResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetParentAuthWithID", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.GetParentAuthWithID(req.ParentID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", resultAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = code.ParentIDNoExist
			resp.Message = fmt.Sprintf(conflictErrorFormat, "parent id not exists")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		}
		return
	}

	spanForHash := h.tracer.StartSpan("CompareHashAndPassword", opentracing.ChildOf(parentSpan))
	err = bcrypt.CompareHashAndPassword([]byte(resultAuth.ParentPW), []byte(req.ParentPW))
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = code.IncorrectParentPWForLogin
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

func (h _default) ChangeParentPW(ctx context.Context, req *proto.ChangeParentPWRequest, resp *proto.ChangeParentPWResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case parentUUIDRegex.MatchString(req.ParentUUID) && req.UUID == req.ParentUUID:
		break
	case adminUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "not parent or admin uuid OR not your student uuid")
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

	spanForDB := h.tracer.StartSpan("GetParentAuthWithUUID", opentracing.ChildOf(parentSpan))
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
			resp.Code = code.IncorrectParentPWForChange
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mismatched hash and password")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "hash compare error, err: " + err.Error())
		}
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.RevisionPW), 3)
	spanForDB = h.tracer.StartSpan("ChangeParentPW", opentracing.ChildOf(parentSpan))
	err = access.ChangeParentPW(string(selectedAuth.UUID), string(hashedBytes))
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

func (h _default) GetParentInformWithUUID(ctx context.Context, req *proto.GetParentInformWithUUIDRequest, resp *proto.GetParentInformWithUUIDResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case studentUUIDRegex.MatchString(req.UUID):
		break
	case adminUUIDRegex.MatchString(req.UUID):
		break
	case teacherUUIDRegex.MatchString(req.UUID):
		break
	case parentUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "not student or admin or teacher or parent uuid")
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

	spanForDB := h.tracer.StartSpan("GetParentInformWithUUID", opentracing.ChildOf(parentSpan))
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

func (h _default) GetParentUUIDsWithInform(ctx context.Context, req *proto.GetParentUUIDsWithInformRequest, resp *proto.GetParentUUIDsWithInformResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case studentUUIDRegex.MatchString(req.UUID):
		break
	case adminUUIDRegex.MatchString(req.UUID):
		break
	case teacherUUIDRegex.MatchString(req.UUID):
		break
	case parentUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "not student or admin or teacher or parent uuid")
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

	informToSelect := &model.ParentInform{
		Name:          model.Name(req.Name),
		PhoneNumber:   model.PhoneNumber(req.PhoneNumber),
	}
	if reflect.DeepEqual(informToSelect, &model.ParentInform{}) {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "bad reqeust")
		return
	}

	spanForDB := h.tracer.StartSpan("GetParentUUIDsWithInform", opentracing.ChildOf(parentSpan))
	selectedUUIDs, err := access.GetParentUUIDsWithInform(informToSelect)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedUUIDs", selectedUUIDs), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = code.ParentWithThatInformNoExist
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

func (h _default) GetChildrenInformsWithUUID(ctx context.Context, req *proto.GetChildrenInformsWithUUIDRequest, resp *proto.GetChildrenInformsWithUUIDResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case adminUUIDRegex.MatchString(req.UUID):
		break
	case parentUUIDRegex.MatchString(req.UUID) && req.UUID == req.ParentUUID:
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "not admin or your parent uuid")
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

	spanForDB := h.tracer.StartSpan("GetStudentInformsWithParentUUID", opentracing.ChildOf(parentSpan))
	selectedInforms, err := access.GetStudentInformsWithParentUUID(req.ParentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedInforms", selectedInforms), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "children not exist parent, uuid: " + req.ParentUUID)
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.ChildrenInform = make([]*proto.StudentInform, len(selectedInforms))
	for index, inform := range selectedInforms {
		resp.ChildrenInform[index] = &proto.StudentInform{
			StudentUUID:   string(inform.StudentUUID),
			Grade:         uint32(inform.Grade),
			Group:         uint32(inform.Class),
			StudentNumber: uint32(inform.StudentNumber),
			Name:          string(inform.Name),
			PhoneNumber:   string(inform.PhoneNumber),
			ImageURI:      string(inform.ProfileURI),
		}
	}

	resp.Status = http.StatusOK
	resp.Message = "get children informs success"
	return
}
