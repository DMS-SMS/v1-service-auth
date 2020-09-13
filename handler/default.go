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
