package handler

import (
	"auth/db"
	"auth/db/access"
	"auth/tool/random"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
	"time"
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

func createGormModelOnCurrentTime() gorm.Model {
	currentTime := time.Now()
	return gorm.Model{
		ID:        uint(random.Int64WithLength(3)),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		DeletedAt: nil,
	}
}