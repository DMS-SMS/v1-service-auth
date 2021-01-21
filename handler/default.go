package handler

import (
	"auth/db"
	"auth/tool/consul"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/opentracing/opentracing-go"
)

type _default struct {
	accessManage db.AccessorManage
	tracer       opentracing.Tracer
	awsSession   *session.Session
	consulAgent  consul.Agent
}

// function signature used in subscriber (add in v.1.1.6)
type SQSMessageHandler func(*sqs.ReceiveMessageOutput)

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

func Manager(am db.AccessorManage) FieldSetter {
	return func(h *_default) {
		h.accessManage = am
	}
}

func Tracer(t opentracing.Tracer) FieldSetter {
	return func(h *_default) {
		h.tracer = t
	}
}

func AWSSession(s *session.Session) FieldSetter {
	return func(h *_default) {
		h.awsSession = s
	}
}

func ConsulAgent(a consul.Agent) FieldSetter {
	return func(h *_default) {
		h.consulAgent = a
	}
}
