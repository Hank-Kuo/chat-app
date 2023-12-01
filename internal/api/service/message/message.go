package message

import (
	"context"
	"time"

	"chat-app/config"
	messageRepo "chat-app/internal/api/repository/message"
	"chat-app/internal/dto"
	"chat-app/internal/models"
	"chat-app/pkg/logger"
	"chat-app/pkg/tracer"

	// "chat-app/pkg/customError"
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
)

type Service interface {
	CreateMessage(ctx context.Context, message *dto.CreateMessageReqDto) error
	CreateReply(ctx context.Context, reply *dto.CreateReplyReqDto) error
	GetMessage(ctx context.Context, channelID string) ([]*models.Message, error)
}

type messageSrv struct {
	cfg         *config.Config
	messageRepo messageRepo.Repository
	snowflake   *snowflake.Node
	logger      logger.Logger
}

func NewService(cfg *config.Config, messageRepo messageRepo.Repository, node *snowflake.Node, logger logger.Logger) Service {
	return &messageSrv{
		cfg:         cfg,
		messageRepo: messageRepo,
		snowflake:   node,
		logger:      logger,
	}
}

func (srv *messageSrv) CreateMessage(ctx context.Context, message *dto.CreateMessageReqDto) error {
	c, span := tracer.NewSpan(ctx, "MessageService.CreateMessage", nil)
	defer span.End()

	id := srv.snowflake.Generate().Int64()

	if err := srv.messageRepo.CreateMessage(c, &models.Message{
		ChannelID: message.ChannelID,
		Bucket:    3,
		MessageID: id,
		Content:   message.Content,
		UserID:    message.UserID,
		Username:  message.Username,
		CreatedAt: time.Now().In(srv.cfg.Server.Location),
	}); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "MessageService.CreateMessage")
	}

	return nil
}

func (srv *messageSrv) CreateReply(ctx context.Context, reply *dto.CreateReplyReqDto) error {
	c, span := tracer.NewSpan(ctx, "MessageService.CreateReply", nil)
	defer span.End()

	id := srv.snowflake.Generate().Int64()

	if err := srv.messageRepo.CreateReply(c, &models.Reply{
		MessageID: reply.MessageID,
		ReplyID:   id,
		Content:   reply.Content,
		UserID:    reply.UserID,
		Username:  reply.Username,
		CreatedAt: time.Now().In(srv.cfg.Server.Location),
	}); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "MessageService.CreateReply")
	}

	return nil
}

func (srv *messageSrv) GetMessage(ctx context.Context, channelID string) ([]*models.Message, error) {
	c, span := tracer.NewSpan(ctx, "MessageService.Get", nil)
	defer span.End()

	messages, err := srv.messageRepo.GetMessage(c, channelID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "MessageService.Get")
	}

	return messages, nil
}

// func (srv *messageSrv) Disconnect(ctx context.Context) error {
// 	ctx, span := tracer.NewSpan(ctx, "MessageService.SendMessage", nil)
// 	defer span.End()
// 	// {"connect_1": group: 1, session_at: 2023/11/1}

// 	_, ok := srv.clients.Channel["2313"]
// 	if ok {

// 	}
// 	fmt.Println(ok)

// 	return nil
// }

// func (srv *messageSrv) ReceivedMessage(ctx context.Context) error {
// 	ctx, span := tracer.NewSpan(ctx, "MessageService.ReceivedMessage", nil)
// 	defer span.End()

// 	_, ok := srv.clients.Channel["2313"]
// 	if ok {

// 	}
// 	fmt.Println(ok)

// 	return nil
// }
