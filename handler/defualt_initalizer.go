package handler

import (
	"log"
	"os"
)

var s3Bucket string
var dmsAPIKey string

func init() {
	if s3Bucket = os.Getenv("SMS_AWS_BUCKET"); s3Bucket == "" {
		log.Fatal("please set SMS_AWS_BUCKET in environment variable")
	}
	if dmsAPIKey = os.Getenv("DMS_API_KEY"); s3Bucket == "" {
		log.Fatal("please set DMS_API_KEY in environment variable")
	}
}
