package handler

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"auth/tool/hash"
	"auth/tool/message"
	"auth/tool/mysqlerr"
	"auth/tool/random"
	code "auth/utils/code/golang"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
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
	err = hash.CompareHashAndPassword(string(resultAuth.StudentPW), req.StudentPW)
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case hash.ErrMismatchedHashAndPassword:
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

	spanForHash := h.tracer.StartSpan("CompareHashAndPassword", opentracing.ChildOf(parentSpan))
	err = hash.CompareHashAndPassword(string(selectedAuth.StudentPW), req.CurrentPW)
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case hash.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = code.IncorrectStudentPWForChange
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mismatched hash and password")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "hash compare error, err: " + err.Error())
		}
		return
	}

	spanForHash = h.tracer.StartSpan("GenerateFromPassword", opentracing.ChildOf(parentSpan))
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.RevisionPW), bcrypt.MinCost)
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to hash pw, err: " + err.Error())
		return
	}

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

	var parentStatus string
	if conn, notify := selectedAuth.ParentStatus.GetBool(); !notify {
		if conn {
			parentStatus = "CONNECTED"
		} else {
			parentStatus = "UN_CONNECTED"
		}
		revisionInform := &model.StudentInform{}
		revisionInform.ParentStatus.SetWithBool(conn, true)

		spanForDB := h.tracer.StartSpan("ModifyStudentInform", opentracing.ChildOf(parentSpan))
		err := access.ModifyStudentInform(string(selectedAuth.StudentUUID), revisionInform)
		spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
		spanForDB.Finish()
		if err != nil {
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "some error occurs in ModifyStudentInform, err: " + err.Error())
			return
		}
	}

	access.Commit()
	resp.Grade = uint32(selectedAuth.Grade)
	resp.Group = uint32(selectedAuth.Class)
	resp.StudentNumber = uint32(selectedAuth.StudentNumber)
	resp.Name = string(selectedAuth.Name)
	resp.PhoneNumber = string(selectedAuth.PhoneNumber)
	resp.ImageURI = string(selectedAuth.ProfileURI)
	resp.ParentStatus = parentStatus

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

	if selectedAuth.ParentUUID == "" {
		access.Commit()
		resp.Status = http.StatusConflict
		resp.Message = fmt.Sprintf(conflictErrorFormat, "not exist teacher uuid in student, uuid: " + req.StudentUUID)
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

func (h _default) GetUnsignedStudentWithAuthCode(ctx context.Context, req *proto.GetUnsignedStudentWithAuthCodeRequest, resp *proto.GetUnsignedStudentWithAuthCodeResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetUnsignedStudentWithAuthCode", opentracing.ChildOf(parentSpan))
	selectedStudent, err := access.GetUnsignedStudentWithAuthCode(int64(req.AuthCode))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedStudent", selectedStudent), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "unsigned student with that auth code is not exist")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		}
		return
	}

	access.Commit()
	resp.AuthCode = uint32(selectedStudent.AuthCode)
	resp.Name = string(selectedStudent.Name)
	resp.Grade = uint32(selectedStudent.Grade)
	resp.Group = uint32(selectedStudent.Class)
	resp.StudentNumber = uint32(selectedStudent.StudentNumber)
	resp.PhoneNumber = string(selectedStudent.PhoneNumber)
	resp.Group = uint32(selectedStudent.Class)

	resp.Status = http.StatusOK
	resp.Message = "succeed to get unsigned student with auth code"
	return
}

func (h _default) CreateNewStudentWithAuthCode(ctx context.Context, req *proto.CreateNewStudentWithAuthCodeRequest, resp *proto.CreateNewStudentWithAuthCodeResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetUnsignedStudentWithAuthCode", opentracing.ChildOf(parentSpan))
	student, err := access.GetUnsignedStudentWithAuthCode(int64(req.AuthCode))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedStudent", student), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "unsigned student with that auth code is not exist")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		}
		return
	}

	var sUUID string
	for {
		sUUID = fmt.Sprintf("student-%s", random.StringConsistOfIntWithLength(12))
		spanForDB := h.tracer.StartSpan("GetStudentAuthWithUUID", opentracing.ChildOf(parentSpan))
		selectedAuth, err := access.GetStudentAuthWithUUID(sUUID)
		spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("selectedAuth", selectedAuth), log.Error(err))
		spanForDB.Finish()
		if err == gorm.ErrRecordNotFound {
			break
		}
		if err != nil {
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
			return
		}
		continue
	}

	spanForDB = h.tracer.StartSpan("GetParentChildWithInform", opentracing.ChildOf(parentSpan))
	child, err := access.GetParentChildWithInform(int64(student.Grade), int64(student.Class), int64(student.StudentNumber), string(student.Name))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedChild", child), log.Error(err))
	spanForDB.Finish()

	var parentUUID string
	var parentConn bool
	switch err {
	case nil:
		parentUUID = string(child.ParentUUID)
		parentConn = true
	case gorm.ErrRecordNotFound:
		parentUUID = ""
		parentConn = false
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		return
	}

	var hashedPW string
	if regexp.MustCompile("^pbkdf2:sha\\d+(:\\d+)?\\$.*\\$.*$").MatchString(req.StudentPW) {
		hashedPW = req.StudentPW
	} else {
		spanForHash := h.tracer.StartSpan("GenerateFromPassword", opentracing.ChildOf(parentSpan))
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.StudentPW), bcrypt.MinCost)
		hashedPW = string(hashedBytes)
		spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
		spanForHash.Finish()

		if err != nil {
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to hash pw, err: " + err.Error())
			return
		}
	}

	spanForDB = h.tracer.StartSpan("CreateStudentAuth", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.CreateStudentAuth(&model.StudentAuth{
		UUID:       model.UUID(sUUID),
		StudentID:  model.StudentID(req.StudentID),
		StudentPW:  model.StudentPW(hashedPW),
		ParentUUID: model.ParentUUID(parentUUID),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreateStudentAuth", resultAuth), log.Error(err))
	spanForDB.Finish()

	switch assertedError := err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for student auth model, err: " + err.Error())
		return
	case *mysql.MySQLError:
		access.Rollback()
		switch assertedError.Number {
		case mysqlcode.ER_DUP_ENTRY:
			key, entry, err := mysqlerr.ParseDuplicateEntryErrorFrom(assertedError)
			if err != nil {
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to parse duplicate error, err: " + err.Error())
				return
			}
			switch key {
			case model.StudentAuthInstance.StudentID.KeyName():
				resp.Status = http.StatusConflict
				resp.Message = fmt.Sprintf(conflictErrorFormat, "student id duplicate, entry: " + entry)
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected duplicate error, key: " + key)
			}
			return
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected CreateStudentAuth error, err: " + assertedError.Error())
			return
		}
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "CreateStudentAuth returns unexpected type of error, err: " + assertedError.Error())
		return
	}

	profileURI := fmt.Sprintf("profiles/uuids/%s", string(resultAuth.UUID))
	spanForDB = h.tracer.StartSpan("CreateStudentInform", opentracing.ChildOf(parentSpan))
	studentInform := &model.StudentInform{
		StudentUUID:   model.StudentUUID(string(resultAuth.UUID)),
		Grade:         student.Grade,
		Class:         student.Class,
		StudentNumber: student.StudentNumber,
		Name:          student.Name,
		PhoneNumber:   student.PhoneNumber,
		ProfileURI:    model.ProfileURI(profileURI),
	}
	studentInform.ParentStatus.SetWithBool(parentConn, false)
	resultInform, err := access.CreateStudentInform(studentInform)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedInform", resultInform), log.Error(err))
	spanForDB.Finish()

	switch assertedError := err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for student inform, err: " + err.Error())
		return
	case *mysql.MySQLError:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected CreateStudentInform error, err: " + assertedError.Error())
		return
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "CreateStudentInform returns unexpected type of error, err: " + assertedError.Error())
		return
	}

	if parentUUID != "" {
		spanForDB = h.tracer.StartSpan("ModifyParentChildren", opentracing.ChildOf(parentSpan))
		revision := &model.ParentChildren{
			StudentUUID: model.StudentUUID(string(resultAuth.UUID)),
		}
		err := access.ModifyParentChildren(child, revision)
		spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
		spanForDB.Finish()

		if err != nil {
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "some error occurs in ModifyParentChildren, err: " + err.Error())
			return
		}
	}

	spanForDB = h.tracer.StartSpan("DeleteUnsignedStudent", opentracing.ChildOf(parentSpan))
	studentInform.ParentStatus.SetWithBool(parentConn, false)
	err = access.DeleteUnsignedStudent(int64(student.AuthCode))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "some error occurs in DeleteUnsignedStudent, err: " + err.Error())
		return
	}

	spanForS3 := h.tracer.StartSpan("CopyObject", opentracing.ChildOf(parentSpan))
	preProfileUri := fmt.Sprintf("profiles/years/2021/grades/%d/groups/%d/numbers/%d", student.Grade, student.Class, student.StudentNumber)
	source := s3Bucket + "/" + preProfileUri
	_, err = s3.New(h.awsSession).CopyObject(&s3.CopyObjectInput{
		Bucket:     aws.String(s3Bucket),
		CopySource: aws.String(url.PathEscape(source)),
		Key:        aws.String(profileURI),
		ACL:        aws.String("public-read"),
	})
	spanForS3.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForS3.Finish()
	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, fmt.Sprintf( "unable to copy s3, err: %v, copy: %s, key: %s", err, preProfileUri, profileURI))
		return
	}

	access.Commit()
	resp.Status = http.StatusCreated
	resp.Message = "succeed to create new student with auth code"
	resp.StudentUUID = string(resultAuth.UUID)

	if studentInform.Grade != 1 {
		return
	}

	smsContent := `
[대덕소프트웨어마이스터고등학교]

신입생 대상 기숙사 지원 시스템(DMS) 안내 문자입니다.

앞서, 저희 학교 지원 시스템(SMS)에 가입해주셔서 감사합니다.

저희는 이 외에도 'DMS'라는 기숙사 지원 시스템을 제공하여 현재 모든 재학생분들이 사용중입니다.

여러분들의 빠른 회원가입을 위하여 SMS에서 입력하신 계정 정보로 DMS 계정을 발급하였습니다.

Play 스토어 또는 App Store에서 'DMS - 기숙사 지원 시스템' 앱을 다운 받아 사용해보세요!

PC 전용 웹 사이트 또한 제공중이니 많이 방문해주세요.
https://www.dsm-dms.com

* 해당 문자는 전공동아리 DMS에서 발신되었습니다.
`

	spanForMsg := h.tracer.StartSpan("SendToReceivers", opentracing.ChildOf(parentSpan))
	jsonResp, err := message.SendToReceivers([]string{string(student.PhoneNumber)}, smsContent, "LMS", "DSM 신입생 대상 기숙사 지원 시스템(DMS) 안내 문자")
	spanForMsg.SetTag("X-Request-Id", reqID).LogFields(log.Object("JsonResponse", jsonResp), log.Error(err))
	spanForMsg.Finish()
	if err != nil {
		resp.Message += fmt.Sprintf("SendToReceivers error: %v", err)
	}

	go func() {
		number, _ := strconv.Atoi(fmt.Sprintf("%d%d%02d", studentInform.Grade, studentInform.Class, studentInform.StudentNumber))
		dmsReq := map[string]interface{}{
			"key":      dmsAPIKey,
			"id":       req.StudentID,
			"password": req.StudentPW,
			"number":   number,
			"name":     string(studentInform.Name),
		}
		dmsReqJson, _ := json.Marshal(dmsReq)
		spanForDMS := h.tracer.StartSpan("PostToDMS", opentracing.ChildOf(parentSpan))
		dmsResp, err := http.Post("https://api.dsm-dms.com/account/signup","application/json", bytes.NewBuffer(dmsReqJson))
		spanForDMS.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
		if err == nil {
			spanForDMS.LogFields(log.Int("status", dmsResp.StatusCode), log.Object("DMSResp", dmsResp))
		}
		spanForDMS.Finish()
	}()

	return
}
