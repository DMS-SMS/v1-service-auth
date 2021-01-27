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

func (h _default) LoginStudentAuth(ctx context.Context, req *proto.LoginStudentAuthRequest, resp *proto.LoginStudentAuthResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetStudentAuthWithID", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.GetStudentAuthWithID(req.StudentID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", resultAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = code.StudentIDNoExist
			resp.Message = fmt.Sprintf(conflictErrorFormat, "student id not exists")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		}
		return
	}

	spanForHash := h.tracer.StartSpan("CompareHashAndPassword", opentracing.ChildOf(parentSpan))
	err = bcrypt.CompareHashAndPassword([]byte(resultAuth.StudentPW), []byte(req.StudentPW))
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = code.IncorrectStudentPWForLogin
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

func (h _default) ChangeStudentPW(ctx context.Context, req *proto.ChangeStudentPWRequest, resp *proto.ChangeStudentPWResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case studentUUIDRegex.MatchString(req.StudentUUID) && req.UUID == req.StudentUUID:
		break
	case adminUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "not student or admin uuid OR not your student uuid")
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

	spanForDB := h.tracer.StartSpan("GetStudentAuthWithUUID", opentracing.ChildOf(parentSpan))
	selectedAuth, err := access.GetStudentAuthWithUUID(req.StudentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", selectedAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "not exist student, uuid: " + req.StudentUUID)
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(selectedAuth.StudentPW), []byte(req.CurrentPW))
	if err != nil {
		access.Rollback()
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = code.IncorrectStudentPWForChange
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mismatched hash and password")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "hash compare error, err: " + err.Error())
		}
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.RevisionPW), bcrypt.MinCost)
	spanForDB = h.tracer.StartSpan("ChangeStudentPW", opentracing.ChildOf(parentSpan))
	err = access.ChangeStudentPW(string(selectedAuth.UUID), string(hashedBytes))
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
	resp.Message = "student pw change success"
	return
}

func (h _default) GetStudentInformWithUUID(ctx context.Context, req *proto.GetStudentInformWithUUIDRequest, resp *proto.GetStudentInformWithUUIDResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetStudentInformWithUUID", opentracing.ChildOf(parentSpan))
	selectedAuth, err := access.GetStudentInformWithUUID(req.StudentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", selectedAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "not exist student, uuid: " + req.StudentUUID)
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.Grade = uint32(selectedAuth.Grade)
	resp.Group = uint32(selectedAuth.Class)
	resp.StudentNumber = uint32(selectedAuth.StudentNumber)
	resp.Name = string(selectedAuth.Name)
	resp.PhoneNumber = string(selectedAuth.PhoneNumber)
	resp.ImageURI = string(selectedAuth.ProfileURI)

	resp.Status = http.StatusOK
	resp.Message = "get student auth success"
	return
}

func (h _default) GetParentWithStudentUUID(ctx context.Context, req *proto.GetParentWithStudentUUIDRequest, resp *proto.GetParentWithStudentUUIDResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	switch true {
	case studentUUIDRegex.MatchString(req.StudentUUID) && req.UUID == req.StudentUUID:
		break
	case adminUUIDRegex.MatchString(req.UUID):
		break
	default:
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "not your student uuid or not admin uuid")
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

	spanForDB := h.tracer.StartSpan("GetStudentAuthWithUUID", opentracing.ChildOf(parentSpan))
	selectedAuth, err := access.GetStudentAuthWithUUID(req.StudentUUID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", selectedAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "not exist student, uuid: " + req.StudentUUID)
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	spanForDB = h.tracer.StartSpan("GetParentInformWithUUID", opentracing.ChildOf(parentSpan))
	selectedParent, err := access.GetParentInformWithUUID(string(selectedAuth.ParentUUID))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("selectedParent", selectedParent), log.Error(err))
	spanForDB.Finish()

	access.Commit()
	resp.ParentUUID = string(selectedParent.ParentUUID)
	resp.Name = string(selectedParent.Name)
	resp.PhoneNumber = string(selectedParent.PhoneNumber)

	resp.Status = http.StatusOK
	resp.Message = "succeed to get parent with student uuid"
	return
}

func (h _default) GetStudentInformsWithUUIDs(ctx context.Context, req *proto.GetStudentInformsWithUUIDsRequest, resp *proto.GetStudentInformsWithUUIDsResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetStudentInformsWithUUIDs", opentracing.ChildOf(parentSpan))
	selectedInforms, err := access.GetStudentInformsWithUUIDs(req.StudentUUIDs)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedInforms", selectedInforms), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = code.StudentUUIDsContainNoExistUUID
			resp.Message = fmt.Sprintf(conflictErrorFormat, "student uuid array contain no exist uuid")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "some error occurs while quering DB, err: " + err.Error())
		}
		return
	}

	access.Commit()
	for _, inform := range selectedInforms {
		resp.StudentInforms = append(resp.StudentInforms, &proto.StudentInform{
			StudentUUID:   string(inform.StudentUUID),
			Grade:         uint32(inform.Grade),
			Group:         uint32(inform.Class),
			StudentNumber: uint32(inform.StudentNumber),
			Name:          string(inform.Name),
			PhoneNumber:   string(inform.PhoneNumber),
			ImageURI:      string(inform.ProfileURI),
		})
	}

	resp.Status = http.StatusOK
	resp.Message = "get student informs success"
	return
}

func (h _default) GetStudentUUIDsWithInform(ctx context.Context, req *proto.GetStudentUUIDsWithInformRequest, resp *proto.GetStudentUUIDsWithInformResponse) (_ error) {
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

	informToSelect := &model.StudentInform{
		Grade:         model.Grade(int64(req.Grade)),
		Class:         model.Class(int64(req.Group)),
		StudentNumber: model.StudentNumber(int64(req.StudentNumber)),
		Name:          model.Name(req.Name),
		PhoneNumber:   model.PhoneNumber(req.PhoneNumber),
		ProfileURI:    model.ProfileURI(req.ImageURI),
	}
	if reflect.DeepEqual(informToSelect, &model.StudentInform{}) {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "bad reqeust")
		return
	}

	spanForDB := h.tracer.StartSpan("GetStudentUUIDsWithInform", opentracing.ChildOf(parentSpan))
	selectedUUIDs, err := access.GetStudentUUIDsWithInform(informToSelect)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedUUIDs", selectedUUIDs), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = code.StudentWithThatInformNoExist
			resp.Message = fmt.Sprintf(conflictErrorFormat, "no exist student with that inform")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.StudentUUIDs = selectedUUIDs
	resp.Status = http.StatusOK
	resp.Message = "get student uuids having that inform success"
	return
}
