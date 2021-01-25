// Add package in v.1.1.6
// listener is function that return closure used in subscribe
// listener.go is file that declare closure listening message from aws sqs, rabbitMQ, etc ...

package subscriber

import (
	"github.com/aws/aws-sdk-go/service/sqs"
)

// function signature type for sqs message handler
type sqsMsgHandler func(*sqs.Message) error
