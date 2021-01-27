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

func (h _default) LoginTeacherAuth(ctx context.Context, req *proto.LoginTeacherAuthRequest, resp *proto.LoginTeacherAuthResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetTeacherAuthWithID", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.GetTeacherAuthWithID(req.TeacherID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", resultAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = code.TeacherIDNoExist
			resp.Message = fmt.Sprintf(conflictErrorFormat, "teacher id not exists")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		}
		return
	}

	spanForHash := h.tracer.StartSpan("CompareHashAndPassword", opentracing.ChildOf(parentSpan))
	err = bcrypt.CompareHashAndPassword([]byte(resultAuth.TeacherPW), []byte(req.TeacherPW))
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = code.IncorrectTeacherPWForLogin
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mismatched hash and password")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "hash compare error, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "succeed to login teacher auth"
	resp.LoggedInTeacherUUID = string(resultAuth.UUID)

	return
}

func (h _default) ChangeTeacherPW(ctx context.Context, req *proto.ChangeTeacherPWRequest, resp *proto.ChangeTeacherPWResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case teacherUUIDRegex.MatchString(req.TeacherUUID) && req.UUID == req.TeacherUUID:
		break
	case adminUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "not teacher or admin uuid OR not your student uuid")
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

	spanForDB := h.tracer.StartSpan("GetTeacherAuthWithUUID", opentracing.ChildOf(parentSpan))
	selectedAuth, err := access.GetTeacherAuthWithUUID(req.TeacherUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", selectedAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "not exist teacher, uuid: " + req.TeacherUUID)
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	spanForHash := h.tracer.StartSpan("CompareHashAndPassword", opentracing.ChildOf(parentSpan))
	err = bcrypt.CompareHashAndPassword([]byte(selectedAuth.TeacherPW), []byte(req.CurrentPW))
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = code.IncorrectTeacherPWForChange
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mismatched hash and password")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "hash compare error, err: " + err.Error())
		}
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.RevisionPW), bcrypt.MinCost)
	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to hash pw, err: " + err.Error())
		return
	}

	spanForDB = h.tracer.StartSpan("ChangeTeacherPW", opentracing.ChildOf(parentSpan))
	err = access.ChangeTeacherPW(string(selectedAuth.UUID), string(hashedBytes))
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
	resp.Message = "teacher pw change success"
	return
}

func (h _default) GetTeacherInformWithUUID(ctx context.Context, req *proto.GetTeacherInformWithUUIDRequest, resp *proto.GetTeacherInformWithUUIDResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetTeacherInformWithUUID", opentracing.ChildOf(parentSpan))
	selectedAuth, err := access.GetTeacherInformWithUUID(req.TeacherUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", selectedAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "not exist teacher, uuid: " + req.TeacherUUID)
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.Grade = uint32(selectedAuth.Grade)
	resp.Group = uint32(selectedAuth.Class)
	resp.Name = string(selectedAuth.Name)
	resp.PhoneNumber = string(selectedAuth.PhoneNumber)

	resp.Status = http.StatusOK
	resp.Message = "get teacher auth success"
	return
}

func (h _default) GetTeacherUUIDsWithInform(ctx context.Context, req *proto.GetTeacherUUIDsWithInformRequest, resp *proto.GetTeacherUUIDsWithInformResponse) (_ error) {
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

	informToSelect := &model.TeacherInform{
		Grade:         model.Grade(int64(req.Grade)),
		Class:         model.Class(int64(req.Group)),
		Name:          model.Name(req.Name),
		PhoneNumber:   model.PhoneNumber(req.PhoneNumber),
	}
	if reflect.DeepEqual(informToSelect, &model.TeacherInform{}) {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "bad reqeust")
		return
	}

	spanForDB := h.tracer.StartSpan("GetTeacherUUIDsWithInform", opentracing.ChildOf(parentSpan))
	selectedUUIDs, err := access.GetTeacherUUIDsWithInform(informToSelect)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedUUIDs", selectedUUIDs), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = code.TeacherWithThatInformNoExist
			resp.Message = fmt.Sprintf(conflictErrorFormat, "no exist teacher with that inform")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.TeacherUUIDs = selectedUUIDs
	resp.Status = http.StatusOK
	resp.Message = "get teacher uuids having that inform success"
	return
}
