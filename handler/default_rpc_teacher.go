package handler

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"auth/tool/mysqlerr"
	"auth/tool/random"
	code "auth/utils/code/golang"
	"context"
	"fmt"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"reflect"
)

func (h _default) CreateNewTeacher(ctx context.Context, req *proto.CreateNewTeacherRequest, resp *proto.CreateNewTeacherResponse) (_ error) {
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

	tUUID, ok := ctx.Value("TeacherUUID").(string)
	if !ok || tUUID == "" {
		tUUID = fmt.Sprintf("teacher-%s", random.StringConsistOfIntWithLength(12))
	}

	for {
		spanForDB := h.tracer.StartSpan("GetTeacherAuthWithUUID", opentracing.ChildOf(parentSpan))
		selectedAuth, err := access.GetTeacherAuthWithUUID(tUUID)
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
		tUUID = fmt.Sprintf("teacher-%s", random.StringConsistOfIntWithLength(12))
		continue
	}

	spanForHash := h.tracer.StartSpan("GenerateFromPassword", opentracing.ChildOf(parentSpan))
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.TeacherPW), bcrypt.MinCost)
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to hash pw, err: " + err.Error())
		return
	}

	spanForDB := h.tracer.StartSpan("CreateTeacherAuth", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.CreateTeacherAuth(&model.TeacherAuth{
		UUID:      model.UUID(tUUID),
		TeacherID: model.TeacherID(req.TeacherID),
		TeacherPW: model.TeacherPW(string(hashedBytes)),
		Certified: false,
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedAuth", resultAuth), log.Error(err))
	spanForDB.Finish()

	switch assertedError := err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for teacher auth model, err: " + err.Error())
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
			case model.TeacherAuthInstance.TeacherID.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.TeacherIDDuplicate
				resp.Message = fmt.Sprintf(conflictErrorFormat, "teacher id duplicate, entry: " + entry)
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected duplicate error, key: " + key)
			}
			return
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected CreateTeacberAuth error, err: " + assertedError.Error())
			return
		}
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "CreateTeacberAuth returns unexpected type of error, err: " + assertedError.Error())
		return
	}

	spanForDB = h.tracer.StartSpan("CreateTeacherInform", opentracing.ChildOf(parentSpan))
	resultInform, err := access.CreateTeacherInform(&model.TeacherInform{
		TeacherUUID:   model.TeacherUUID(string(resultAuth.UUID)),
		Grade:         model.Grade(int64(req.Grade)),
		Class:         model.Class(int64(req.Group)),
		Name:          model.Name(req.Name),
		PhoneNumber:   model.PhoneNumber(req.PhoneNumber),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedInform", resultInform), log.Error(err))
	spanForDB.Finish()

	switch assertedError := err.(type) {
	case nil:
		break
	case validator.ValidationErrors:
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for teacher inform, err: " + err.Error())
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
			case model.TeacherInformInstance.PhoneNumber.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.TeacherPhoneNumberDuplicate
				resp.Message = fmt.Sprintf(conflictErrorFormat, "phone number duplicate, entry: " + entry)
			case model.TeacherInformInstance.Class.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.TeacherClassDuplicate
				resp.Message = fmt.Sprintf(conflictErrorFormat, "class duplicate, entry: " + entry)
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected duplicate error, key: " + key)
			}
			return
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected CreateTeacherInform error, err: " + assertedError.Error())
			return
		}
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "CreateTeacherInform returns unexpected type of error, err: " + assertedError.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusCreated
	resp.Message = "new teacher create success"
	resp.CreatedTeacherUUID = tUUID

	return
}

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
