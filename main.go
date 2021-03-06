package main

import (
	"auth/consul"
	consulagent "auth/consul/agent"
	"auth/db"
	"auth/db/access"
	"auth/handler"
	proto "auth/proto/golang/auth"
	"auth/subscriber"
	"auth/tool/closure"
	"auth/tool/network"
	topic "auth/utils/topic/golang"
	"fmt"
	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/checkers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client/selector"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/transport/grpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"os"
	"time"
)

func main() {
	// create service
	port := network.GetRandomPortNotInUsedWithRange(10000, 10100) // change from function to method (in v.1.1.6)
	service := micro.NewService(
		micro.Name(topic.AuthServiceName),
		micro.Version("1.1.6"),
		micro.Transport(grpc.NewTransport()),
		micro.Address(fmt.Sprintf(":%d", port)),
	)
	srvID := fmt.Sprintf("%s-%s", service.Server().Options().Name, service.Server().Options().Id)

	// create consul connection
	consulAddr := os.Getenv("CONSUL_ADDRESS")
	if consulAddr == "" {
		log.Fatal("please set CONSUL_ADDRESS in environment variable")
	}
	consulCfg := api.DefaultConfig()
	consulCfg.Address = consulAddr
	consulCli, err := api.NewClient(consulCfg)
	if err != nil {
		log.Fatalf("consul connect fail, err: %v", err)
	}
	consulAgent := consulagent.Default( // add in v.1.1.6
		consulagent.Strategy(selector.RoundRobin),
		consulagent.Client(consulCli),
		consulagent.Services([]consul.ServiceName{topic.AuthServiceName, topic.ClubServiceName,
			topic.OutingServiceName, topic.ScheduleServiceName, topic.AnnouncementServiceName}),
	)

	// create db access manager
	dbc, _, err := db.ConnectWithConsul(consulCli, "db/auth/local")
	if err != nil {
		log.Fatalf("db connect fail, err: %v", err)
	}
	db.Migrate(dbc)
	accessManage, err := db.NewAccessorManage(access.Default(dbc))
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
		Tags:        []opentracing.Tag{{"sid", srvID}},
		Reporter:    &jaegercfg.ReporterConfig{LogSpans: true, LocalAgentHostPort: jaegerAddr},
		Sampler:     &jaegercfg.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1},
	}.NewTracer()
	if err != nil {
		log.Fatalf("error while creating new tracer for service, err: %v", err)
	}
	defer func() {
		_ = closer.Close()
	}()

	// create AWS session
	awsId := os.Getenv("SMS_AWS_ID")
	if awsId == "" {
		log.Fatal("please set SMS_AWS_ID in environment variable")
	}
	awsKey := os.Getenv("SMS_AWS_KEY")
	if awsKey == "" {
		log.Fatal("please set SMS_AWS_KEY in environment variable")
	}
	s3Region := os.Getenv("SMS_AWS_REGION")
	if s3Region == "" {
		log.Fatal("please set SMS_AWS_REGION in environment variable")
	}
	awsSession, err := session.NewSession(&aws.Config{
		Region:      aws.String(s3Region),
		Credentials: credentials.NewStaticCredentials(awsId, awsKey, ""),
	})
	if err != nil {
		log.Fatalf("error while creating new aws session, err: %v", err)
	}

	// create gRPC handler
	defaultHandler := handler.Default(
		handler.Manager(accessManage),
		handler.Tracer(authSrvTracer),
		handler.AWSSession(awsSession),
		handler.ConsulAgent(consulAgent),
	)

	// create subscriber & register listener (add in v.1.1.6)
	consulChangeQueue := os.Getenv("CHANGE_CONSUL_SQS_AUTH")
	if consulChangeQueue == "" {
		log.Fatal("please set CHANGE_CONSUL_SQS_AUTH in environment variable")
	}
	subscriber.SetAwsSession(awsSession)
	defaultSubscriber := subscriber.Default()
	//defaultSubscriber.RegisterBeforeStart(
	//	subscriber.SqsQueuePurger(consulChangeQueue),
	//)
	//defaultSubscriber.RegisterListeners(
	//	subscriber.SqsMsgListener(consulChangeQueue, defaultHandler.ChangeConsulNodes, &sqs.ReceiveMessageInput{
	//		MaxNumberOfMessages: aws.Int64(10),
	//		WaitTimeSeconds:     aws.Int64(2),
	//	}),
	//)

	// register initializer for service
	service.Init(
		micro.BeforeStart(consulAgent.ChangeAllServiceNodes),
		micro.AfterStart(consulAgent.ChangeAllServiceNodes),
		micro.AfterStart(defaultSubscriber.StartListening),
		micro.AfterStart(consulAgent.ServiceNodeRegistry(service.Server())),
		micro.BeforeStop(consulAgent.ServiceNodeDeregistry(service.Server())),
	)

	// register gRPC handler in service
	_ = proto.RegisterAuthAdminHandler(service.Server(), defaultHandler)
	_ = proto.RegisterAuthStudentHandler(service.Server(), defaultHandler)
	_ = proto.RegisterAuthTeacherHandler(service.Server(), defaultHandler)
	_ = proto.RegisterAuthParentHandler(service.Server(), defaultHandler)
	_ = proto.RegisterAuthEventHandler(service.Server(), defaultHandler)

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
		OnComplete: closure.TTLCheckHandlerAboutDB(service.Server(), consulCli),
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
