package message

import (
	"context"

	messagePb "github.com/Hank-Kuo/chat-app/pb/message"
	"github.com/Hank-Kuo/chat-app/pkg/logger"
	"github.com/Hank-Kuo/chat-app/pkg/manager"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"
)

type grpcHandler struct {
	manager *manager.ClientManager
	logger  logger.Logger
}

func NewGrpcHandler(manager *manager.ClientManager, logger logger.Logger) *grpcHandler {
	return &grpcHandler{manager: manager, logger: logger}
}

func (h *grpcHandler) MessageReceived(ctx context.Context, request *messagePb.ReceiveMessageRequest) (*messagePb.ReceiveMessageResponse, error) {

	_, span := tracer.NewSpan(ctx, "MessageGrpcHandler.ReceiveMessage", nil)
	defer span.End()

	h.manager.ToClientChan <- manager.ToClientInfo{OriginClientId: request.OriginClientID, ClientId: request.ClientId, InstanceId: request.InstanceId, Data: request.Message}

	return &messagePb.ReceiveMessageResponse{}, nil

}

func (h *grpcHandler) ReplyReceived(ctx context.Context, request *messagePb.ReceiveReplyRequest) (*messagePb.ReceiveReplyResponse, error) {

	_, span := tracer.NewSpan(ctx, "MessageGrpcHandler.ReceiveReply", nil)
	defer span.End()

	h.manager.ToClientChan <- manager.ToClientInfo{OriginClientId: request.OriginClientID, ClientId: request.ClientId, InstanceId: request.InstanceId, Data: request.Reply}

	return &messagePb.ReceiveReplyResponse{}, nil
}
