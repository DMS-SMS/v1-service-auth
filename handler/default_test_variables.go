package handler

import (
	"auth/db"
	"auth/db/access"
	"fmt"
	"github.com/stretchr/testify/mock"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
)

var (
	defaultHandler *_default
	mockForDB *mock.Mock
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