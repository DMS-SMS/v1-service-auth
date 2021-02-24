// add package in v.1.0.5
// message package is used for sending sms or mms message from 'ALIGO 문자 서비스'
// send.go is file that sending message

package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

type SendMassToReceiversResponse struct {
	ResultCode int    `json:"result_code"`
	Message    string `json:"message"`
	MsgID      string `json:"msg_id"`
	SuccessCnt int    `json:"success_cnt"`
	ErrorCnt   int    `json:"error_cnt"`
	MsgType    string `json:"msg_type"`
}

type SendToReceiversResponse struct {
	SendMassToReceiversResponse
}

func SendMassToReceivers(receivers, contents []string, _type, title string) (jsonResp SendMassToReceiversResponse, err error) {
	if (len(receivers) != len(contents)) || (len(receivers) < 1 || len(contents) < 1) {
		err = errors.New("receivers & contents must be same length bigger than 0")
		return
	}

	if _type != "" && _type != "SMS" && _type != "LMS" && _type != "MMS" {
		err = errors.New("type value must be blank or SMS or LMS or MMS")
		return
	}

	if (_type == "SMS" || _type == "") && title != "" {
		err = errors.New("cannot set title when type is black or SMS")
		return
	}

	req, err := http.NewRequest("POST", "https://apis.aligo.in/send_mass/", nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("some error occurs while creating request, err: %v", err))
		return
	}

	q := req.URL.Query()
	q.Add("key", aligoAPIKey)
	q.Add("user_id", aligoAccountID)
	q.Add("sender", aligoSender)
	if _type != "" {
		q.Add("msg_type", _type)
	}
	if title != "" {
		q.Add("title", title)
	}
	for i, receiver := range receivers {
		q.Add(fmt.Sprintf("rec_%d", i+1), receiver)
		q.Add(fmt.Sprintf("msg_%d", i+1), contents[i])
	}
	q.Add("cnt", strconv.Itoa(len(receivers)))
	req.URL.RawQuery = q.Encode()

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("some error occurs while sending request, err: %v", err))
		return
	}

	jsonResp = SendMassToReceiversResponse{}
	_ = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if resp.StatusCode != http.StatusOK || jsonResp.ResultCode != 0 {
		err = errors.New(fmt.Sprintf("failed to send mass message, status: %d, json response: %v", resp.StatusCode, jsonResp))
		return
	}

	return
}

func SendToReceivers(receivers []string, content , _type, title string) (jsonResp SendToReceiversResponse, err error) {
	if len(receivers) < 1 {
		err = errors.New("receivers must be longer than 0")
		return
	}

	if _type != "" && _type != "SMS" && _type != "LMS" && _type != "MMS" {
		err = errors.New("type value must be blank or SMS or LMS or MMS")
		return
	}

	if (_type == "SMS" || _type == "") && title != "" {
		err = errors.New("cannot set title when type is black or SMS")
		return
	}

	req, err := http.NewRequest("POST", "https://apis.aligo.in/send/", nil)
	if err != nil {
		err = errors.New(fmt.Sprintf("some error occurs while creating request, err: %v", err))
		return
	}

	q := req.URL.Query()
	q.Add("key", aligoAPIKey)
	q.Add("user_id", aligoAccountID)
	q.Add("sender", aligoSender)
	if _type != "" {
		q.Add("msg_type", _type)
	}
	if title != "" {
		q.Add("title", title)
	}
	q.Add("receiver", strings.Join(receivers, ","))
	q.Add("msg", content)
	req.URL.RawQuery = q.Encode()

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		err = errors.New(fmt.Sprintf("some error occurs while sending request, err: %v", err))
		return
	}

	jsonResp = SendToReceiversResponse{}
	_ = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if resp.StatusCode != http.StatusOK || jsonResp.ResultCode != 0 {
		err = errors.New(fmt.Sprintf("failed to send message, status: %d, json response: %v", resp.StatusCode, jsonResp))
		return
	}

	return
}
