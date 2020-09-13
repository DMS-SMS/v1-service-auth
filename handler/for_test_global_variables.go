package handler

import (
	"auth/db"
	"auth/db/access"
	"fmt"
	"github.com/stretchr/testify/mock"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
	//proto "auth/proto/golang/auth"
)

var (
	defaultHandler *_default
	mockForAccessor *mock.Mock
)

func init() {
	mockForAccessor = new(mock.Mock)

	exampleTracerForRPCService, closer, err := jaegercfg.Configuration{ServiceName: "DMS.SMS.v1.service.auth"}.NewTracer()
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new tracer for service, err: %v", err)) }
	defer func() { _ = closer.Close() }()

	mockAccessManager, err := db.NewAccessorManage(access.Mock(mockForAccessor))
	if err != nil { log.Fatal(fmt.Sprintf("error while creating new access manage with mock, err: %v", err)) }

	defaultHandler = &_default{
		manager: mockAccessManager,
		tracer:  exampleTracerForRPCService,
	}
}
