// Add package in v.1.1.6
// subscriber package is used for handling event message occurred by SNS, RabbitMQ, etc ...
// you can start subscribe by calling Start method with parameter, specific signature function

package subscriber

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/prometheus/common/log"
)

var awsSession *session.Session

func SetAwsSession(s *session.Session) {
	awsSession = s
}

type _default struct {
	awsSession *session.Session
	listeners  []func()
}

type FieldSetter func(*_default)

func Default(setters ...FieldSetter) *_default {
	return newDefault(setters...)
}

func newDefault(setters ...FieldSetter) (h *_default) {
	h = new(_default)
	h.awsSession = awsSession
	for _, setter := range setters {
		setter(h)
	}
	h.listeners = []func(){}
	return
}

func AwsSession(awsSession *session.Session) FieldSetter {
	return func(s *_default) {
		s.awsSession = awsSession
	}
}

func (s *_default) RegisterListeners(listeners ...func()) {
	s.listeners = append(s.listeners, listeners...)
}

func (s *_default) StartListening() (_ error) {
	log.Info("Default subscriber start listening!!")
	for _, listener := range s.listeners {
		go listener()
	}
	return
}
