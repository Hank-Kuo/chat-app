package message

import (
	"context"
	"fmt"

	"github.com/Hank-Kuo/chat-app/internal/models"
	"github.com/Hank-Kuo/chat-app/pkg/tracer"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/pkg/errors"
)

func (r *messageRepo) PublishMessage(ctx context.Context, message *models.Message) error {
	ctx, span := tracer.NewSpan(ctx, "messageRepo.PublishMessage", nil)
	defer span.End()

	topic := fmt.Sprintf("message_%s", message.ChannelID)
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("topic_message"),
		Value: []byte(fmt.Sprintf(`{
			"channel_id": "%s", "message_id": %d, "content": "%s", "user_id": "%s", "username": "%s", "created_at": "%s"
			}`, message.ChannelID, message.MessageID, message.Content, message.UserID,
			message.Username, message.CreatedAt.Format("2006-01-02T15:04:05.999999999Z"))),
	}

	if err := r.kafkaProducer.Produce(msg); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "messageRepo.PublishMessage")
	}

	return nil
}

func (r *messageRepo) PublishReply(ctx context.Context, reply *models.Reply) error {
	ctx, span := tracer.NewSpan(ctx, "messageRepo.PublishReply", nil)
	defer span.End()

	topic := fmt.Sprintf("message_%s", reply.ChannelID)
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte("topic_reply"),
		Value: []byte(fmt.Sprintf(`{
			"channel_id": "%s", "message_id": %d, 
			"reply_id": %d, "user_id": "%s", "username": "%s", 
			"created_at": "%s", "content": "%s"}`, reply.ChannelID,
			reply.MessageID, reply.ReplyID, reply.UserID, reply.Username, reply.CreatedAt.Format("2006-01-02T15:04:05.999999999Z"), reply.Content)),
	}

	if err := r.kafkaProducer.Produce(msg); err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "messageRepo.PublishReply")
	}

	return nil
}
