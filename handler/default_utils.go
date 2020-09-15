package handler

import (
	"context"
	"github.com/google/uuid"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/uber/jaeger-client-go"
)

func (_ _default) getContextFromMetadata(ctx context.Context) (parsedCtx context.Context, proxyAuthenticated bool, reason string) {
	md, ok := metadata.FromContext(ctx)
	if !ok {
		proxyAuthenticated = false
		reason = "metadata not exists"
		return
	}

	reqID, ok := md.Get("X-Request-Id")
	if !ok {
		proxyAuthenticated = false
		reason = "X-Request-Id not exists"
		return
	}

	_, err := uuid.Parse(reqID)
	if err != nil {
		proxyAuthenticated = false
		reason = "X-Request-ID invalid, err: " + err.Error()
		return
	}

	spanCtx, ok := md.Get("Span-Context")
	if !ok {
		proxyAuthenticated = false
		reason = "Span-Context not exists"
		return
	}

	parentSpan, err := jaeger.ContextFromString(spanCtx)
	if err != nil {
		proxyAuthenticated = false
		reason = "Span-Context invalid, err: " + err.Error()
		return
	}

	proxyAuthenticated = true
	reason = ""

	parsedCtx = context.Background()
	parsedCtx = context.WithValue(parsedCtx, "X-Request-Id", reqID)
	parsedCtx = context.WithValue(parsedCtx, "Span-Context", parentSpan)

	if sUUID, ok := md.Get("StudentUUID"); ok { parsedCtx = context.WithValue(parsedCtx, "StudentUUID", sUUID) }
	if tUUID, ok := md.Get("TeacherUUID"); ok { parsedCtx = context.WithValue(parsedCtx, "TeacherUUID", tUUID) }
	if pUUID, ok := md.Get("ParentUUID"); ok  { parsedCtx = context.WithValue(parsedCtx, "ParentUUID", pUUID) }

	return
}