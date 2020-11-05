package main

import (
	"auth/adapter"
	"auth/db"
	"auth/db/access"
	"auth/handler"
	proto "auth/proto/golang/auth"
	"auth/tool/closure"
	topic "auth/utils/topic/golang"
	"fmt"
	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/checkers"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"math/rand"
	"net"
	"os"
	"time"
)

func main() {
	// create consul connection
	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		log.Fatal("please set CONSUL_ADDRESS in environment variable")
	}
	consulCfg := api.DefaultConfig()
	consulCfg.Address = consulAddr
	consul, err := api.NewClient(consulCfg)
	if err != nil {
		log.Fatalf("consul connect fail, err: %v", err)
	}

	// create db access manager
	dbc, _, err := adapter.ConnectDBWithConsul(consul, "db/auth/local")
	if err != nil {
		log.Fatalf("db connect fail, err: %v", err)
	}
	db.Migrate(dbc)
	defaultAccessManage, err := db.NewAccessorManage(access.Default(dbc))
	if err != nil {
		log.Fatalf("db accessor create fail, err: %v", err)
	}

	// create jaeger connection
	jaegerAddr := os.Getenv("JAEGER_ADDRESS")
	if jaegerAddr == "" {
		log.Fatal("please set JAEGER_ADDRESS in environment variable")
	}
	authSrvTracer, closer, err := jaegercfg.Configuration{
		ServiceName: topic.AuthServiceName,
		Reporter:    &jaegercfg.ReporterConfig{LogSpans: true, LocalAgentHostPort: jaegerAddr},
		Sampler:     &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1},
	}.NewTracer()
	if err != nil {
		log.Fatalf("error while creating new tracer for service, err: %v", err)
	}
	defer func() {
		_ = closer.Close()
	}()

	// create gRPC handler
	rpcHandler := handler.Default(
		handler.AWSSession(nil),
		handler.Manager(defaultAccessManage),
		handler.Tracer(authSrvTracer),
	)

	// create service
	port := getRandomPortNotInUsedWithRange(10000, 10100)
	service := micro.NewService(
		micro.Name(topic.AuthServiceName),
		micro.Version("1.1.2"),
		micro.Transport(grpc.NewTransport()),
		micro.Address(fmt.Sprintf(":%d", port)),
	)

	// register initializer for service
	service.Init(
		micro.AfterStart(closure.ConsulServiceRegistrar(service.Server(), consul)),
		micro.BeforeStop(closure.ConsulServiceDeregistrar(service.Server(), consul)),
	)

	// register gRPC handler in service
	_ = proto.RegisterAuthAdminHandler(service.Server(), rpcHandler)
	_ = proto.RegisterAuthStudentHandler(service.Server(), rpcHandler)
	_ = proto.RegisterAuthTeacherHandler(service.Server(), rpcHandler)
	_ = proto.RegisterAuthParentHandler(service.Server(), rpcHandler)

	// run DB Health checker
	h := health.New()
	dbChecker, err := checkers.NewSQL(&checkers.SQLConfig{
		Pinger: dbc.DB(),
	})
	if err != nil {
		log.Fatalf("unable to create sql health checker, err: %v", err)
	}
	dbHealthCfg := &health.Config{
		Name:       "DB-Checker",
		Checker:    dbChecker,
		Interval:   time.Second * 5,
		OnComplete: closure.TTLCheckHandlerAboutDB(service.Server(), consul),
	}
	if err = h.AddChecks([]*health.Config{dbHealthCfg}); err != nil {
		log.Fatalf("unable to register health checks, err: %v", err)
	}
	if err = h.Start(); err != nil {
		log.Fatalf("unable to start health checks, err: %v", err)
	}

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func getRandomPortNotInUsedWithRange(min, max int) (port int) {
	for {
		port = rand.Intn(max - min) + min
		conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			continue
		}
		_ = conn.Close()
		break
	}
	return
}

//http.HandleFunc("/profiles", func(writer http.ResponseWriter, request *http.Request) {
//	file, fileHeader, err := request.FormFile("profile")
//	if err != nil {
//		http.Error(writer, fmt.Sprintf("%s, err: %v", http.StatusText(http.StatusBadRequest), err.Error()), http.StatusBadRequest)
//		return
//	}
//
//	buf := make([]byte, fileHeader.Size)
//	_, _ = file.Read(buf)
//
//	service := micro.NewService()
//	authService := proto.NewAuthAdminService("DMS.SMS.v1.service.auth", service.Client())
//	now := time.Now()
//	fmt.Println(authService.CreateNewStudent(context.Background(), &proto.CreateNewStudentRequest{
//		UUID:          "",
//		StudentID:     "",
//		StudentPW:     "",
//		ParentUUID:    "",
//		Grade:         0,
//		Class:         0,
//		StudentNumber: 0,
//		Name:          "",
//		PhoneNumber:   "",
//		Image:         buf,
//	}))
//	fmt.Println(time.Now().Sub(now).Seconds())
//	return
//})
//
//log.Fatal(http.ListenAndServe(":8080", nil))
