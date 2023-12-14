package message

import (
	"context"
	"fmt"
	"time"

	"chat-app/config"
	messageRepo "chat-app/internal/api/repository/message"
	"chat-app/internal/dto"
	"chat-app/internal/models"
	"chat-app/pkg/logger"
	"chat-app/pkg/tracer"
	"chat-app/pkg/utils"

	"chat-app/pkg/customError"
	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
)

type Service interface {
	CreateMessage(ctx context.Context, message *dto.CreateMessageReqDto) (*models.Message, error)
	CreateReply(ctx context.Context, reply *dto.CreateReplyReqDto) error
	GetMessage(ctx context.Context, m *dto.GetMessageQueryDto) (*dto.GetMessageResDto, error)
	GetReply(ctx context.Context, messageID int64) ([]*models.Reply, error)
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

func (srv *messageSrv) CreateMessage(ctx context.Context, message *dto.CreateMessageReqDto) (*models.Message, error) {
	c, span := tracer.NewSpan(ctx, "MessageService.CreateMessage", nil)
	defer span.End()

	id := srv.snowflake.Generate().Int64()
	m := &models.Message{
		ChannelID: message.ChannelID,
		Bucket:    utils.MakeBucket(id),
		MessageID: id,
		Content:   message.Content,
		UserID:    message.UserID,
		Username:  message.Username,
		CreatedAt: time.Now().In(srv.cfg.Server.Location),
	}
	if err := srv.messageRepo.CreateMessage(c, m); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "MessageService.CreateMessage")
	}

	return m, nil
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

type Cursor struct {
	NextCursor string
}

func (srv *messageSrv) GetMessage(ctx context.Context, m *dto.GetMessageQueryDto) (*dto.GetMessageResDto, error) {
	c, span := tracer.NewSpan(ctx, "MessageService.GetMessage", nil)
	defer span.End()
	cursor := int64(-1)

	if len(m.Cursor) > 0 {
		_cursor, err := utils.DecodeCursor("message_id", m.Cursor)
		cursor = _cursor
		if err != nil {
			return nil, customError.ErrBadQueryParams
		}
	}

	limit := 11
	messages, err := srv.messageRepo.GetMessage(c, m.ChannelID, cursor, limit)

	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "MessageService.GetMessage")
	}
	// adjust cursor
	var nextCursor string
	if len(messages) >= limit {
		nextCursor = utils.EncodeCursor("message_id", messages[limit-1].MessageID)
		messages = messages[:limit-1]
	}

	return &dto.GetMessageResDto{
		Messages:   messages,
		NextCursor: nextCursor,
	}, nil
}

func (srv *messageSrv) GetReply(ctx context.Context, messageID int64) ([]*models.Reply, error) {
	c, span := tracer.NewSpan(ctx, "MessageService.GetReply", nil)
	defer span.End()

	reply, err := srv.messageRepo.GetReply(c, messageID)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "MessageService.GetReply")
	}

	return reply, nil
}
