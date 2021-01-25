// Add package in v.1.1.6
// listener is function that return closure used in subscribe
// listener.go is file that declare closure listening message from aws sqs, rabbitMQ, etc ...

package subscriber

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/micro/go-micro/v2/util/log"
	systemlog "log"
)

// function signature type for sqs message handler
type sqsMsgHandler func(*sqs.Message) error

// function that returns closure listening aws sqs message & handling with function receive from parameter
func SqsMsgListener(queue string, handler sqsMsgHandler, rcvInput *sqs.ReceiveMessageInput) func() {
	sqsSrv := sqs.New(awsSession)
	urlResult, err := sqsSrv.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(queue),
	})
	if err != nil {
		systemlog.Fatalf("unable to get queue url from queue name, name: %s\n", queue)
	}

	if rcvInput == nil {
		rcvInput = &sqs.ReceiveMessageInput{}
	}
	rcvInput.QueueUrl = urlResult.QueueUrl

	return func() {
		var rcvOutput *sqs.ReceiveMessageOutput
		var err error
		var msg *sqs.Message

		for {
			rcvOutput, err = sqsSrv.ReceiveMessage(rcvInput)
			if err != nil {
				log.Errorf("some error occurs while pulling from aws sqs, queue: %s, err: %s\n", rcvInput.QueueUrl, err)
				return
			}

			for _, msg = range rcvOutput.Messages {
				go func(msg *sqs.Message) {
					if err := handler(msg); err != nil {
						log.Errorf("some error occurs while handling aws sqs message, queue: %s, msg id: %s err: %s\n", rcvInput.QueueUrl, msg.MessageId, err)
					}
				} (msg)
			}
		}
	}
}
