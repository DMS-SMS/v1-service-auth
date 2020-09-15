package handler

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"auth/tool/mysqlerr"
	"auth/tool/random"
	"bytes"
	"context"
	"fmt"
	mysqlcode "github.com/VividCortex/mysqlerr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/uber/jaeger-client-go"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
)

func(h _default) CreateNewStudent(ctx context.Context, req *proto.CreateNewStudentRequest, resp *proto.CreateNewStudentResponse) (_ error) {
	const (
		forbiddenMessageFormat = "forbidden (reason: %s)"
		proxyAuthRequiredMessageFormat = "proxy auth required (reason: %s)"
		internalServerErrorFormat = "internal server error (reason: %s)"
		conflictErrorFormat = "conflict (reason: %s)"
	)

	adminUUIDRegex := regexp.MustCompile("^admin-\\d{12}")
	if !adminUUIDRegex.MatchString(req.UUID) {
		resp.Status = http.StatusForbidden
		resp.Message = fmt.Sprintf(forbiddenMessageFormat, "you are not admin")
		return
	}

	md, ok := metadata.FromContext(ctx)
	if !ok {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "metadata not exists")
		return
	}

	reqID, ok := md.Get("X-Request-Id")
	if !ok {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "X-Request-Id not exists")
		return
	}

	_, err := uuid.Parse(reqID)
	if err != nil {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "X-Request-ID invalid, err: " + err.Error())
		return
	}

	spanCtx, ok := md.Get("Span-Context")
	if !ok {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "Span-Context not exists")
		return
	}

	parentSpan, err := jaeger.ContextFromString(spanCtx)
	if err != nil {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "Span-Context invalid, err: " + err.Error())
		return
	}

	access, err := h.manager.BeginTx()
	if err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "tx begin fail, err: " + err.Error())
		return
	}

	sUUID, ok := md.Get("StudentUUID")
	if !ok {
		sUUID = fmt.Sprintf("student-%s", random.StringConsistOfIntWithLength(12))
	}

	for {
		spanForDB := h.tracer.StartSpan("CheckIfStudentAuthExists", opentracing.ChildOf(parentSpan))
		exist, err := access.CheckIfStudentAuthExists(sUUID)
		spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Bool("exist", exist), log.Error(err))
		spanForDB.Finish()
		if err != nil {
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " + err.Error())
			return
		}
		if !exist {
			break
		}
		sUUID = fmt.Sprintf("student-%s", random.StringConsistOfIntWithLength(12))
		continue
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.StudentPW), 1)
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
		switch assertedError.Number{
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
				resp.Code = CodeStudentIDDuplicate
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
				resp.Code = CodeParentUUIDNoExist
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
	}

	if h.awsSession != nil {
		_, err = s3.New(h.awsSession).PutObject(&s3.PutObjectInput{
			Bucket: aws.String("dms-sms"),
			Key:    aws.String(fmt.Sprintf("profiles/%s", string(resultAuth.UUID))),
			Body:   bytes.NewReader(req.Image),
		})
		if err != nil {
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to upload profile to s3, err: " + err.Error())
			return
		}
	}

	spanForDB = h.tracer.StartSpan("CreateStudentInform", opentracing.ChildOf(parentSpan))
	resultInform, err := access.CreateStudentInform(&model.StudentInform{
		StudentUUID:   model.StudentUUID(string(resultAuth.UUID)),
		Grade:         model.Grade(int64(req.Grade)),
		Class:         model.Class(int64(req.Class)),
		StudentNumber: model.StudentNumber(int64(req.StudentNumber)),
		Name:          model.Name(req.Name),
		PhoneNumber:   model.PhoneNumber(req.PhoneNumber),
		ProfileURI:    model.ProfileURI(fmt.Sprintf("profiles/%s", string(resultAuth.UUID))),
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
				resp.Code = CodeStudentNumberDuplicate
				resp.Message = fmt.Sprintf(conflictErrorFormat, "student number duplicate, entry: " + entry)
			case model.StudentInformInstance.PhoneNumber.KeyName():
				resp.Status = http.StatusConflict
				resp.Code = CodePhoneNumberDuplicate
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
	}

	access.Commit()
	resp.Status = http.StatusCreated
	resp.Message = "new student create success"
	resp.CreatedStudentUUID = sUUID

	return
}
