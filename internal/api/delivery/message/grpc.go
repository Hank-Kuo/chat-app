package message

import (
	"context"

	messageSrv "github.com/Hank-Kuo/chat-app/internal/api/service/message"
	messagePb "github.com/Hank-Kuo/chat-app/pb/message"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"
)

type grpcHandler struct {
	messageSrv messageSrv.Service
	logger     logger.Logger
}

func NewGrpcHandler(messageSrv messageSrv.Service, logger logger.Logger) *grpcHandler {
	return &grpcHandler{messageSrv: messageSrv, logger: logger}
}

func (h *grpcHandler) ReceiveMessage(ctx context.Context, request *messagePb.ReceiveMessageRequest) (*messagePb.ReceiveMessageResponse, error) {

	_, span := tracer.NewSpan(ctx, "MessageGrpcHandler.ReceiveMessage", nil)
	defer span.End()

	return &messagePb.ReceiveMessageResponse{}, nil

}

func (h *grpcHandler) ReceiveReply(ctx context.Context, request *messagePb.ReceiveReplyRequest) (*messagePb.ReceiveReplyResponse, error) {

	_, span := tracer.NewSpan(ctx, "MessageGrpcHandler.ReceiveReply", nil)
	defer span.End()

	return &messagePb.ReceiveReplyResponse{}, nil
}
