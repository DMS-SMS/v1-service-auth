package handler

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"auth/tool/mysqlerr"
	"auth/tool/random"
	code "auth/utils/code/golang"
	"bytes"
	"context"
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
)

func (h _default) CreateNewStudent(ctx context.Context, req *proto.CreateNewStudentRequest, resp *proto.CreateNewStudentResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	if !adminUUIDRegex.MatchString(req.UUID) {
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not admin")
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

	sUUID, ok := ctx.Value("StudentUUID").(string)
	if !ok || sUUID == "" {
		sUUID = fmt.Sprintf("student-%s", random.StringConsistOfIntWithLength(12))
	}

	for {
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
		sUUID = fmt.Sprintf("student-%s", random.StringConsistOfIntWithLength(12))
		continue
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.StudentPW), 3)
	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to hash pw, err: " + err.Error())
		return
	}

	spanForDB := h.tracer.StartSpan("CreateStudentAuth", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.CreateStudentAuth(&model.StudentAuth{
		UUID:       model.UUID(sUUID),
		StudentID:  model.StudentID(req.StudentID),
		StudentPW:  model.StudentPW(string(hashedBytes)),
		ParentUUID: model.ParentUUID(req.ParentUUID),
	})
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("CreatedAuth", resultAuth), log.Error(err))
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
				resp.Code = code.StudentIDDuplicate
				resp.Message = fmt.Sprintf(conflictErrorFormat, "student id duplicate, entry: " + entry)
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected duplicate error, key: " + key)
			}
			return
		case mysqlcode.ER_NO_REFERENCED_ROW_2:
			fkInform, _, err := mysqlerr.ParseFKConstraintFailErrorFrom(assertedError)
			if err != nil {
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to parse fk contraint error, err: " + err.Error())
				return
			}
			switch fkInform.ConstraintName {
			case model.StudentAuthInstance.ParentUUIDConstraintName():
				resp.Status = http.StatusConflict
				resp.Code = code.ParentUUIDNoExist
				resp.Message = fmt.Sprintf(conflictErrorFormat, "FK constraint fail, FK name: " + fkInform.AttrName)
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected FK constraint fail, FK name: " + fkInform.AttrName)
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

	if string(req.Image) == "" {
		access.Rollback()
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "image is empty byte array")
		return
	}

	profileURI := fmt.Sprintf("profiles/%s", string(resultAuth.UUID))
	spanForDB = h.tracer.StartSpan("CreateStudentInform", opentracing.ChildOf(parentSpan))
	resultInform, err := access.CreateStudentInform(&model.StudentInform{
		StudentUUID:   model.StudentUUID(string(resultAuth.UUID)),
		Grade:         model.Grade(int64(req.Grade)),
		Class:         model.Class(int64(req.Group)),
		StudentNumber: model.StudentNumber(int64(req.StudentNumber)),
		Name:          model.Name(req.Name),
		PhoneNumber:   model.PhoneNumber(req.PhoneNumber),
		ProfileURI:    model.ProfileURI(profileURI),
	})
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
		switch assertedError.Number {
		case mysqlcode.ER_DUP_ENTRY:
			key, entry, err := mysqlerr.ParseDuplicateEntryErrorFrom(assertedError)
			if err != nil {
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to parse duplicate error, err: " + err.Error())
				return
			}
			switch key {
			case model.StudentInformInstance.StudentNumber.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.StudentNumberDuplicate
				resp.Message = fmt.Sprintf(conflictErrorFormat, "student number duplicate, entry: " + entry)
			case model.StudentInformInstance.PhoneNumber.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.StudentPhoneNumberDuplicate
				resp.Message = fmt.Sprintf(conflictErrorFormat, "phone number duplicate entry: " + entry)
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected duplicate error, key: " + key)
			}
			return
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected CreateStudentInform error, err: " + assertedError.Error())
			return
		}
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "CreateStudentInform returns unexpected type of error, err: " + assertedError.Error())
		return
	}

	if h.awsSession != nil {
		spanForS3 := h.tracer.StartSpan("PutObject", opentracing.ChildOf(parentSpan))
		_, err = s3.New(h.awsSession).PutObject(&s3.PutObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(profileURI),
			Body:   bytes.NewReader(req.Image),
			ACL:    aws.String("public-read"),
		})
		spanForS3.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
		spanForS3.Finish()
		if err != nil {
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to upload profile to s3, err: " + err.Error())
			return
		}
	}

	access.Commit()
	resp.Status = http.StatusCreated
	resp.Message = "new student create success"
	resp.CreatedStudentUUID = sUUID

	return
}

func (h _default) CreateNewTeacher(ctx context.Context, req *proto.CreateNewTeacherRequest, resp *proto.CreateNewTeacherResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	if !adminUUIDRegex.MatchString(req.UUID) {
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not admin")
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
	
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.TeacherPW), 3)
	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to hash pw, err: " + err.Error())
		return
	}

	spanForDB := h.tracer.StartSpan("CreateTeacherAuth", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.CreateTeacherAuth(&model.TeacherAuth{
		UUID:       model.UUID(tUUID),
		TeacherID:  model.TeacherID(req.TeacherID),
		TeacherPW:  model.TeacherPW(string(hashedBytes)),
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

func (h _default) CreateNewParent(ctx context.Context, req *proto.CreateNewParentRequest, resp *proto.CreateNewParentResponse) (_ error) {
	ctx, proxyAuthenticated, reason := h.getContextFromMetadata(ctx)
	if !proxyAuthenticated {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, reason)
		return
	}

	if !adminUUIDRegex.MatchString(req.UUID) {
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not admin")
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

	pUUID, ok := ctx.Value("ParentUUID").(string)
	if !ok || pUUID == "" {
		pUUID = fmt.Sprintf("parent-%s", random.StringConsistOfIntWithLength(12))
	}

	for {
		spanForDB := h.tracer.StartSpan("GetParentAuthWithUUID", opentracing.ChildOf(parentSpan))
		selectedAuth, err := access.GetParentAuthWithUUID(pUUID)
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
		pUUID = fmt.Sprintf("parent-%s", random.StringConsistOfIntWithLength(12))
		continue
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.ParentPW), 3)
	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to hash pw, err: " + err.Error())
		return
	}

	spanForDB := h.tracer.StartSpan("CreateParentAuth", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.CreateParentAuth(&model.ParentAuth{
		UUID:     model.UUID(pUUID),
		ParentID: model.ParentID(req.ParentID),
		ParentPW: model.ParentPW(string(hashedBytes)),
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
			case model.ParentAuthInstance.ParentID.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.ParentIDDuplicate
				resp.Message = fmt.Sprintf(conflictErrorFormat, "parent id duplicate, entry: " + entry)
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
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "CreateParentAuth returns unexpected type of error, err: " + assertedError.Error())
		return
	}

	spanForDB = h.tracer.StartSpan("CreateParentInform", opentracing.ChildOf(parentSpan))
	resultInform, err := access.CreateParentInform(&model.ParentInform{
		ParentUUID:  model.ParentUUID(string(resultAuth.UUID)),
		Name:        model.Name(req.Name),
		PhoneNumber: model.PhoneNumber(req.PhoneNumber),
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
			case model.ParentInformInstance.PhoneNumber.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = code.ParentPhoneNumberDuplicate
				resp.Message = fmt.Sprintf(conflictErrorFormat, "phone number duplicate, entry: " + entry)
			default:
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected duplicate error, key: " + key)
			}
			return
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected CreateParentInform error, err: " + assertedError.Error())
			return
		}
	default:
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "CreateParentInform returns unexpected type of error, err: " + assertedError.Error())
		return
	}

	access.Commit()
	resp.Status = http.StatusCreated
	resp.Message = "new parent create success"
	resp.CreatedParentUUID = pUUID

	return
}

func (h _default) LoginAdminAuth(ctx context.Context, req *proto.LoginAdminAuthRequest, resp *proto.LoginAdminAuthResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetAdminAuthWithID", opentracing.ChildOf(parentSpan))
	resultAuth, err := access.GetAdminAuthWithID(req.AdminID)
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedAuth", resultAuth), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusConflict
			resp.Code = code.AdminIDNoExist
			resp.Message = fmt.Sprintf(conflictErrorFormat, "admin id not exists")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(resultAuth.AdminPW), []byte(req.AdminPW))
	if err != nil {
		access.Rollback()
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			resp.Status = http.StatusConflict
			resp.Code = code.IncorrectAdminPWForLogin
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mismatched hash and password")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "hash compare error, err: " + err.Error())
		}
		return
	}

	access.Commit()
	resp.Status = http.StatusOK
	resp.Message = "succeed to login admin auth"
	resp.LoggedInAdminUUID = string(resultAuth.UUID)

	return
}
