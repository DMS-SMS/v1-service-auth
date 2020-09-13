package handler

import (
	"bufio"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
	"os"
	"path/filepath"
)

const (
	validAdminUUID = "admin-111111111111"
	validParentUUID = "parent-111111111111"
	validStudentID = "jinhong0719"
	validStudentPW = "testPW"
	validGrade = 2
	validClass = 2
	validStudentNumber = 7
	validName = "박진홍"
	validPhoneNumber = "01088378347"
)

var (
	validImageByteArr []byte
	validSpanContextString string
)

func init() {
	exampleTracerForAPIGateway, closer, err := jaegercfg.Configuration{ServiceName: "DMS.SMS.v1.api.gateway"}.NewTracer()
	if err != nil { log.Fatal(err) }
	defer func() { _ = closer.Close() }()
	exampleTracerForRPCService, closer, err := jaegercfg.Configuration{ServiceName: "DMS.SMS.v1.service.auth"}.NewTracer()
	if err != nil { log.Fatal(err) }
	defer func() { _ = closer.Close() }()

	exampleSpanForAPIGateway := exampleTracerForAPIGateway.StartSpan("v1/students")
	exampleSpanForRPCService := exampleTracerForRPCService.StartSpan("CreateNewStudent", opentracing.ChildOf(exampleSpanForAPIGateway.Context()))
	validSpanContextString = exampleSpanForRPCService.Context().(jaeger.SpanContext).String()

	absPath, err := filepath.Abs("./images_for_test/doraemon.png")
	if err != nil { log.Fatal(err) }
	file, err := os.Open(absPath)
	if err != nil { log.Fatal(err) }
	fileInfo, err := file.Stat()
	if err != nil { log.Fatal(err) }
	validImageByteArr = make([]byte, fileInfo.Size())
	fileReader := bufio.NewReader(file)
	_, err = fileReader.Read(validImageByteArr)
	if err != nil { log.Fatal(err) }
}