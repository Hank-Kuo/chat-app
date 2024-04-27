package message

import (
	"context"
	"fmt"

	"github.com/Hank-Kuo/chat-app/internal/models"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func (r *messageRepo) PublishMessage(ctx context.Context, userID string, message *models.Message) error {
	ctx, span := tracer.NewSpan(ctx, "messageRepo.PublishMessage", nil)
	defer span.End()

	msg := kafka.Message{
		Value: []byte(fmt.Sprintf(`{
			"received_user_id": "%s", "channel_id": "%s", "message_id": "%s", 
			"user_id": "%s", "username": "%s", "created_at": "%s", "content": "%s"
			}`, userID, message.ChannelID, message.MessageID, message.UserID,
			message.Username, message.CreatedAt, message.Content)),
	}

	if err := r.kafkaMessageWriter.WriteMessages(ctx, msg); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "messageRepo.PublishMessage")
	}

	return nil
}

func (r *messageRepo) PublishReply(ctx context.Context, userID string, reply *models.Reply) error {
	ctx, span := tracer.NewSpan(ctx, "messageRepo.PublishReply", nil)
	defer span.End()

	msg := kafka.Message{
		Value: []byte(fmt.Sprintf(`{
			"received_user_id": "%s", "channel_id": "%s", "message_id": "%s", 
			"reply_id": "%s", "user_id": "%s", "username": "%s", 
			"created_at": "%s", "content": "%s"}`, userID, reply.ChannelID,
			reply.MessageID, reply.UserID, reply.Username, reply.CreatedAt, reply.Content)),
	}

	if err := r.kafkaReplyWriter.WriteMessages(ctx, msg); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "messageRepo.PublishReply")
	}

	return nil
}
