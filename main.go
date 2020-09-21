package main

import (
	"auth/adapter"
	"auth/db"
	"auth/db/access"
	"auth/handler"
	proto "auth/proto/golang/auth"
	"auth/tool/closure"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/transport/grpc"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"os"
)

func main() {
	consulCfg := api.DefaultConfig()
	if addr := os.Getenv("CONSUL_ADDRESS"); addr != "" {
		consulCfg.Address = addr
	}

	consul, err := api.NewClient(consulCfg)
	if err != nil {
		log.Fatalf("consul connect fail, err: %v", err)
	}

	dbc, _, err := adapter.ConnectDBWithConsul(consul, "db/auth/local")
	if err != nil {
		log.Fatalf("db connect fail, err: %v", err)
	}
	db.Migrate(dbc)

	defaultAccessManage, err := db.NewAccessorManage(access.Default(dbc))
	if err != nil {
		log.Fatalf("db accessor create fail, err: %v", err)
	}

	// 이 부분은 나중에 Consul 조회로 변경 예
	agentHost := "localhost:6831"
	if addr := os.Getenv("JAEGER_ADDRESS"); addr != "" {
		agentHost = addr
	}

	authSrvTracer, closer, err := jaegercfg.Configuration{
		ServiceName: "DMS.SMS.v1.service.auth",
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: agentHost,
		},
	}.NewTracer()
	if err != nil {
		log.Fatalf("error while creating new tracer for service, err: %v", err)
	}
	defer func() {
		_ = closer.Close()
	}()

	rpcHandler := handler.Default(
		handler.AWSSession(nil),
		handler.Manager(defaultAccessManage),
		handler.Tracer(authSrvTracer),
	)

	service := micro.NewService(
		micro.Name("DMS.SMS.v1.service.auth"),
		micro.Version("1.0.0"),
		micro.Transport(grpc.NewTransport()),
	)

	service.Init(
		micro.AfterStart(closure.ConsulServiceRegistrar(service.Server(), consul)),
		micro.BeforeStop(closure.ConsulServiceDeregistrar(service.Server(), consul)),
	)

	_ = proto.RegisterAuthAdminHandler(service.Server(), rpcHandler)
	_ = proto.RegisterAuthStudentHandler(service.Server(), rpcHandler)
	_ = proto.RegisterAuthTeacherHandler(service.Server(), rpcHandler)
	_ = proto.RegisterAuthParentHandler(service.Server(), rpcHandler)

	// health checker 실행 추가

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
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
