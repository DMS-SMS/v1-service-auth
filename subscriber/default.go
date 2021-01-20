// Add package in v.1.1.6
// subscriber package is used for handling event message occurred by SNS, RabbitMQ, etc ...
// you can start subscribe by calling Start method with parameter, specific signature function

package subscriber

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type _default struct {
	awsSession *session.Session
}
