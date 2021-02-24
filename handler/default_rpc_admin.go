package handler

import (
	"auth/model"
	proto "auth/proto/golang/auth"
	"auth/tool/message"
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
	"math/rand"
	"net/http"
	"time"
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

	spanForHash := h.tracer.StartSpan("GenerateFromPassword", opentracing.ChildOf(parentSpan))
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.StudentPW), bcrypt.MinCost)
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

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

	profileURI := fmt.Sprintf("profiles/uuids/%s", string(resultAuth.UUID))
	spanForDB = h.tracer.StartSpan("CreateStudentInform", opentracing.ChildOf(parentSpan))
	studentInform := &model.StudentInform{
		StudentUUID:   model.StudentUUID(string(resultAuth.UUID)),
		Grade:         model.Grade(int64(req.Grade)),
		Class:         model.Class(int64(req.Group)),
		StudentNumber: model.StudentNumber(int64(req.StudentNumber)),
		Name:          model.Name(req.Name),
		PhoneNumber:   model.PhoneNumber(req.PhoneNumber),
		ProfileURI:    model.ProfileURI(profileURI),
	}
	if req.ParentUUID != "" {
		studentInform.ParentStatus.SetWithBool(true, false)
	} else {
		studentInform.ParentStatus.SetWithBool(false, false)
	}
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

	var pUUID string
	for {
		pUUID = fmt.Sprintf("parent-%s", random.StringConsistOfIntWithLength(12))
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
		continue
	}

	spanForHash := h.tracer.StartSpan("GenerateFromPassword", opentracing.ChildOf(parentSpan))
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.ParentPW), bcrypt.MinCost)
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

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

	for _, child := range req.ChildrenInform {
		spanForDB := h.tracer.StartSpan("GetStudentUUIDsWithInform", opentracing.ChildOf(parentSpan))
		uuidArr, err := access.GetStudentUUIDsWithInform(&model.StudentInform{
			Grade:         model.Grade(int64(child.Grade)),
			Class:         model.Class(int64(child.Group)),
			StudentNumber: model.StudentNumber(int64(child.StudentNumber)),
			Name:          model.Name(child.Name),
		})
		spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("UUIDs", uuidArr), log.Error(err))
		spanForDB.Finish()
		
		var childUUID string
		if len(uuidArr) >= 1 {
			childUUID = uuidArr[0]
		}

		spanForDB = h.tracer.StartSpan("CreateParentChildren", opentracing.ChildOf(parentSpan))
		resultChild, err := access.CreateParentChildren(&model.ParentChildren{
			ParentUUID:    model.ParentUUID(string(resultAuth.UUID)),
			Grade:         model.Grade(int64(child.Grade)),
			Class:         model.Class(int64(child.Group)),
			StudentNumber: model.StudentNumber(int64(child.StudentNumber)),
			Name:          model.Name(child.Name),
			StudentUUID:   model.StudentUUID(childUUID),
		})
		spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("ResultChild", resultChild), log.Error(err))
		spanForDB.Finish()

		switch assertedError := err.(type) {
		case nil:
			break
		case validator.ValidationErrors:
			access.Rollback()
			resp.Status = http.StatusProxyAuthRequired
			resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for parent children, err: " + err.Error())
			return
		case *mysql.MySQLError:
			access.Rollback()
			resp.Status = http.StatusConflict
			resp.Message = fmt.Sprintf(conflictErrorFormat, "mysql error occurs in CreateParentChildren, err: " + err.Error())
			return
		default:
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "CreateParentInform returns unexpected type of error, err: " + assertedError.Error())
			return
		}

		if childUUID != "" {
			spanForDB = h.tracer.StartSpan("ModifyStudentInform", opentracing.ChildOf(parentSpan))
			revisionStudent := &model.StudentInform{}
			revisionStudent.ParentStatus.SetWithBool(true, false)
			err = access.ModifyStudentInform(uuidArr[0], revisionStudent)
			spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
			spanForDB.Finish()

			if err != nil {
				access.Rollback()
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "ModifyStudentInform returns error, err: " + err.Error())
				return
			}

			spanForDB = h.tracer.StartSpan("ChangeParentUUID", opentracing.ChildOf(parentSpan))
			err = access.ChangeParentUUID(uuidArr[0], string(resultAuth.UUID))
			spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
			spanForDB.Finish()

			if err != nil {
				access.Rollback()
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "ChangeParentUUID returns error, err: " + err.Error())
				return
			}
		}
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

	spanForHash := h.tracer.StartSpan("CompareHashAndPassword", opentracing.ChildOf(parentSpan))
	err = bcrypt.CompareHashAndPassword([]byte(resultAuth.AdminPW), []byte(req.AdminPW))
	spanForHash.SetTag("X-Request-Id", reqID).LogFields(log.Error(err))
	spanForHash.Finish()

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

func (h _default) AddUnsignedStudents(ctx context.Context, req *proto.AddUnsignedStudentsRequest, resp *proto.AddUnsignedStudentsResponse) (_ error) {
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

	access, err := h.accessManage.BeginTx()
	if err != nil {
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "tx begin fail, err: " + err.Error())
		return
	}

	svc := s3.New(h.awsSession)
	var addCount uint32 = 0
	var noAddCount uint32 = 0
	var duplicateLog string

	for _, student := range req.Students {
		preProfileUri := fmt.Sprintf("profiles/years/2021/grades/%d/groups/%d/numbers/%d", student.Grade, student.Group, student.StudentNumber)
		_, err := svc.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(preProfileUri),
		})
		if err != nil {
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf("pre profile not exist in s3, uri: %s, name: %s. (not save anything)", preProfileUri, student.Name)
			access.Rollback()
			return
		}

		rand.Seed(time.Now().UnixNano())
		min := 100000
		max := 999999
		authCode := rand.Intn(max - min + 1) + min

		_, err = access.AddUnsignedStudent(&model.UnsignedStudent{
			AuthCode:      model.AuthCode(int64(authCode)),
			Grade:         model.Grade(int64(student.Grade)),
			Class:         model.Class(int64(student.Group)),
			StudentNumber: model.StudentNumber(int64(student.StudentNumber)),
			Name:          model.Name(student.Name),
			PhoneNumber:   model.PhoneNumber(student.PhoneNumber),
			PreProfileURI: model.PreProfileURI(preProfileUri),
		})

		switch assertedError := err.(type) {
		case nil:
			addCount++
			continue
		case validator.ValidationErrors:
			access.Rollback()
			resp.Status = http.StatusProxyAuthRequired
			resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "invalid data for unsigned students, err: " + err.Error())
			return
		case *mysql.MySQLError:
			switch assertedError.Number {
			case mysqlcode.ER_DUP_ENTRY:
				noAddCount++
				duplicateLog += fmt.Sprintf("\nuri: %s, name: %s, duplicate err: %v", preProfileUri, student.Name, assertedError)
				continue
			default:
				access.Rollback()
				resp.Status = http.StatusInternalServerError
				resp.Message = fmt.Sprintf(internalServerErrorFormat, "unexpected AddUnsignedStudent error, err: " + assertedError.Error())
				return
			}
		default:
			access.Rollback()
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "AddUnsignedStudent returns unexpected type of error, err: " + assertedError.Error())
			return
		}
	}
	access.Commit()

	resp.Status = http.StatusCreated
	resp.Message = fmt.Sprintf("succeed to add unsigned students. duplicate log: %s", duplicateLog)
	resp.AddCount = addCount
	resp.NoAddCount = noAddCount
	return
}

func (h _default) SendJoinSMSToUnsignedStudents(ctx context.Context, req *proto.SendJoinSMSToUnsignedStudentsRequest, resp *proto.SendJoinSMSToUnsignedStudentsResponse) (_ error) {
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

	spanForDB := h.tracer.StartSpan("GetUnsignedStudents", opentracing.ChildOf(parentSpan))
	selectedStudents, err := access.GetUnsignedStudents(int64(req.TargetGrade), int64(req.TargetGroup))
	spanForDB.SetTag("X-Request-Id", reqID).LogFields(log.Object("SelectedStudents", selectedStudents), log.Error(err))
	spanForDB.Finish()

	if err != nil {
		access.Rollback()
		switch err {
		case gorm.ErrRecordNotFound:
			resp.Status = http.StatusNotFound
			resp.Message = fmt.Sprintf(notFoundMessageFormat, "unsigned student not exists with that grade & group")
		default:
			resp.Status = http.StatusInternalServerError
			resp.Message = fmt.Sprintf(internalServerErrorFormat, "unable to query DB, err: " +err.Error())
		}
		return
	}

	smsFormat := `
[대덕소프트웨어마이스터고등학교]

학교 지원 시스템(SMS) 회원가입 안내 문자입니다.

[가입 대상: %d%d%02d %s]
[인증 번호: %d]

Play 스토어 또는 앱스토어에서 'SMS - 학교 지원 시스템' 앱 다운로드 후 진행해주세요.

모든 재학생분들(신입생 포함)은 3/5(금)까지 회원가입을 완료해주세요.

* 해당 문자는 전공동아리 DMS에서 발신되었습니다.
* 페이스북 'DSM 기숙사 지원 시스템' 페이지 팔로우!`

	receivers := make([]string, len(selectedStudents))
	contents := make([]string, len(selectedStudents))
	for i, student := range selectedStudents {
		receivers[i] = string(student.PhoneNumber)
		contents[i] = fmt.Sprintf(smsFormat, student.Grade, student.Class, student.StudentNumber, student.Name, student.AuthCode)
	}

	spanForMsg := h.tracer.StartSpan("SendMassToReceivers", opentracing.ChildOf(parentSpan))
	jsonResp, err := message.SendMassToReceivers(receivers, contents, "LMS", "DSM 학교 지원 시스템(SMS) 회원가입 안내")
	spanForMsg.SetTag("X-Request-Id", reqID).LogFields(log.Object("JsonResponse", jsonResp), log.Error(err))
	spanForMsg.Finish()

	if err != nil {
		access.Rollback()
		resp.Status = http.StatusInternalServerError
		resp.Message = fmt.Sprintf(internalServerErrorFormat, "fail to send mass message, err: " + err.Error())
		return
	}

	resp.Status = http.StatusOK
	resp.Message = "succeed to send mass message"
	resp.SendCount = uint32(jsonResp.SuccessCnt)
	resp.NoSendCount = uint32(jsonResp.ErrorCnt)
	access.Commit()
	return
}
