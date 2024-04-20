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
	"chat-app/pkg/utils"

	"chat-app/pkg/customError"

	"github.com/bwmarrin/snowflake"
	"github.com/pkg/errors"
)

const LIMIT = 11

type Service interface {
	CreateMessage(ctx context.Context, message *dto.CreateMessageReqDto) (*models.Message, error)
	CreateReply(ctx context.Context, reply *dto.CreateReplyReqDto) (*models.Reply, error)
	GetMessage(ctx context.Context, m *dto.GetMessageQueryDto) (*dto.GetMessageResDto, error)
	GetReply(ctx context.Context, r *dto.GetReplyQueryDto) (*dto.GetReplyResDto, error)
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

func (srv *messageSrv) CreateReply(ctx context.Context, reply *dto.CreateReplyReqDto) (*models.Reply, error) {
	c, span := tracer.NewSpan(ctx, "MessageService.CreateReply", nil)
	defer span.End()

	id := srv.snowflake.Generate().Int64()
	r := &models.Reply{
		MessageID: reply.MessageID,
		ReplyID:   id,
		Content:   reply.Content,
		UserID:    reply.UserID,
		Username:  reply.Username,
		CreatedAt: time.Now().In(srv.cfg.Server.Location),
	}
	if err := srv.messageRepo.CreateReply(c, r); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "MessageService.CreateReply")
	}

	return r, nil
}

func (srv *messageSrv) GetMessage(ctx context.Context, m *dto.GetMessageQueryDto) (*dto.GetMessageResDto, error) {
	c, span := tracer.NewSpan(ctx, "MessageService.GetMessage", nil)
	defer span.End()

	// parse cursor
	cursor := int64(-1)
	if len(m.Cursor) > 0 {
		_cursor, err := utils.DecodeCursor("message_id", m.Cursor)
		if err != nil {
			return nil, customError.ErrBadQueryParams
		}
		cursor = _cursor
	}

	messages, err := srv.messageRepo.GetMessage(c, m.ChannelID, cursor, LIMIT)

	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "MessageService.GetMessage")
	}

	// update message & add nextCursor
	var nextCursor string
	if len(messages) >= LIMIT {
		nextCursor = utils.EncodeCursor("message_id", messages[LIMIT-1].MessageID)
		messages = messages[:LIMIT-1]
	}

	return &dto.GetMessageResDto{
		Messages:   messages,
		NextCursor: nextCursor,
	}, nil
}

func (srv *messageSrv) GetReply(ctx context.Context, r *dto.GetReplyQueryDto) (*dto.GetReplyResDto, error) {
	c, span := tracer.NewSpan(ctx, "MessageService.GetReply", nil)
	defer span.End()

	cursor := int64(-1)

	if len(r.Cursor) > 0 {
		_cursor, err := utils.DecodeCursor("reply_id", r.Cursor)
		if err != nil {
			return nil, customError.ErrBadQueryParams
		}
		cursor = _cursor
	}

	reply, err := srv.messageRepo.GetReply(c, r.MessageID, cursor, LIMIT)
	if err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "MessageService.GetReply")
	}

	// update message & add nextCursor
	var nextCursor string
	if len(reply) >= LIMIT {
		nextCursor = utils.EncodeCursor("reply_id", reply[LIMIT-1].ReplyID)
		reply = reply[:LIMIT-1]
	}

	return &dto.GetReplyResDto{
		Replies:    reply,
		NextCursor: nextCursor,
	}, nil
}

func (srv *messageSrv) CreateMessage1(ctx context.Context, message *dto.CreateMessageReqDto) {
}
