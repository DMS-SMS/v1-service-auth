// add package in v.1.0.5
// message package is used for sending sms or mms message from 'ALIGO 문자 서비스'
// send.go is file that sending message

package message

import (
	"log"
	"os"
)

var aligoAPIKey string
var aligoAccountID string
var aligoSender string

func init() {
	if aligoAPIKey = os.Getenv("ALIGO_API_KEY"); aligoAPIKey == "" {
		log.Fatal("please set ALIGO_API_KEY in environment variable")
	}

	if aligoAccountID = os.Getenv("ALIGO_ACCOUNT_ID"); aligoAccountID == "" {
		log.Fatal("please set ALIGO_ACCOUNT_ID in environment variable")
	}

	if aligoSender = os.Getenv("ALIGO_SENDER"); aligoSender == "" {
		log.Fatal("please set ALIGO_SENDER in environment variable")
	}
}
