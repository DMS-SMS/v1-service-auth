package handler

import (
	proto "auth/proto/golang/auth"
	"context"
	log "github.com/micro/go-micro/v2/logger"
)

func (h _default) ChangeAllServiceNodes(ctx context.Context, req *proto.Empty, resp *proto.Empty) (_ error) {
	err := h.consulAgent.ChangeAllServiceNodes()
	log.Infof("change all service nodes!, err: %v", err)
	return
}
