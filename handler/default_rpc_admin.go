package handler

import (
	proto "auth/proto/golang/auth"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/uber/jaeger-client-go"
	"net/http"
)

func(h _default) CreateNewStudent(ctx context.Context, req *proto.CreateNewStudentRequest, resp *proto.CreateNewStudentResponse) (_ error) {
	const (
		proxyAuthRequiredMessageFormat = "proxy auth required (reason: %s)"
	)

	md, ok := metadata.FromContext(ctx)
	if !ok {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "metadata not exists")
		return
	}

	reqID, ok := md.Get("X-Request-Id")
	if !ok {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "X-Request-Id not exists")
		return
	}

	_, err := uuid.Parse(reqID)
	if err != nil {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "X-Request-ID invalid, err: " + err.Error())
	}

	spanCtx, ok := md.Get("Span-Context")
	if !ok {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "Span-Context not exists")
		return
	}

	parentSpan, err := jaeger.ContextFromString(spanCtx)
	if err != nil {
		resp.Status = http.StatusProxyAuthRequired
		resp.Message = fmt.Sprintf(proxyAuthRequiredMessageFormat, "Span-Context invalid, err: " + err.Error())
		return
	}

	return
}
