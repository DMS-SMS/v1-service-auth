package handler

import (
	"auth/db"
	"auth/db/access"
	"fmt"
	"github.com/stretchr/testify/mock"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
	"regexp"
)

var (
	defaultHandler *_default
	mockForDB *mock.Mock

	adminUUIDRegex = regexp.MustCompile("^admin-\\d{12}")
	studentUUIDRegex = regexp.MustCompile("^student-\\d{12}")
	teacherUUIDRegex = regexp.MustCompile("^teacher-\\d{12}")
	parentUUIDRegex = regexp.MustCompile("^parent-\\d{12}")
)

const (
	forbiddenMessageFormat = "forbidden (reason: %s)"
	notFoundMessageFormat = "not found (reason: %s)"
	proxyAuthRequiredMessageFormat = "proxy auth required (reason: %s)"
	conflictErrorFormat = "conflict (reason: %s)"
	internalServerErrorFormat = "internal server error (reason: %s)"

)

func init() {
	mockForDB = new(mock.Mock)

	exampleTracerForRPCService, closer, err := jaegercfg.Configuration{ServiceName: "DMS.SMS.v1.service.auth"}.NewTracer()
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new tracer for service, err: %v", err)) }
	defer func() { _ = closer.Close() }()

	mockAccessManage, err := db.NewAccessorManage(access.Mock(mockForDB))
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new access manage with mock, err: %v", err)) }

	defaultHandler = &_default{
		accessManage: mockAccessManage,
		tracer:       exampleTracerForRPCService,
	}
}
