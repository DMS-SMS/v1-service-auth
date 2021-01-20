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

type FieldSetter func(*_default)

func Default(setters ...FieldSetter) *_default {
	return newDefault(setters...)
}

func newDefault(setters ...FieldSetter) (h *_default) {
	h = new(_default)
	for _, setter := range setters {
		setter(h)
	}
	return
}

func AWSSession(awsSession *session.Session) FieldSetter {
	return func(s *_default) {
		s.awsSession = awsSession
	}
}
