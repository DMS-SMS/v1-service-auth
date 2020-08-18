package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"auth/handler"
	"auth/subscriber"

	auth "auth/proto/auth"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("DMS.SMS.v1.service.auth"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	auth.RegisterAuthHandler(service.Server(), new(handler.Auth))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("DMS.SMS.v1.service.auth", service.Server(), new(subscriber.Auth))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
