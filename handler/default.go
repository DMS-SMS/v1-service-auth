package handler

import (
	"auth/db"
	"github.com/opentracing/opentracing-go"
)

type _default struct {
	None
	manager db.AccessorManage
	tracer  opentracing.Tracer
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

func Manager(m db.AccessorManage) FieldSetter {
	return func(h *_default) {
		h.manager = m
	}
}

func Tracer(t opentracing.Tracer) FieldSetter {
	return func(h *_default) {
		h.tracer = t
	}
}