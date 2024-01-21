package message

import (
	"chat-app/internal/models"
	"chat-app/pkg/tracer"
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func (r *messageRepo) CreateMessage(ctx context.Context, message *models.Message) error {
	ctx, span := tracer.NewSpan(ctx, "MessageRepo.CreateMessage", nil)
	defer span.End()

	sqlQuery := `INSERT INTO message (channel_id, bucket, message_id, 
							content, user_id, username, created_at) 
					VALUES (?, ?, ?, ?, ?, ?, ?)`

	err := r.session.Query(sqlQuery, message.ChannelID, message.Bucket,
		message.MessageID, message.Content, message.UserID, message.Username, message.CreatedAt).WithContext(ctx).Exec()

	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "MessageRepo.CreateMessage")
	}
	return nil
}
func (r *messageRepo) CreateReply(ctx context.Context, reply *models.Reply) error {
	ctx, span := tracer.NewSpan(ctx, "MessageRepo.CreateReply", nil)
	defer span.End()

	sqlQuery := `INSERT INTO reply (message_id, reply_id,
							content, user_id, username, created_at) 
					VALUES (?, ?, ?, ?, ?, ?)`

	err := r.session.Query(sqlQuery, reply.MessageID, reply.ReplyID,
		reply.Content, reply.UserID, reply.Username, reply.CreatedAt).WithContext(ctx).Exec()

	if err != nil {
		tracer.AddSpanError(span, err)
		return errors.Wrap(err, "MessageRepo.CreateReply")
	}

	return nil
}

func (r *messageRepo) GetMessage(ctx context.Context, channelID string, cursor int64, limit int) ([]*models.Message, error) {
	ctx, span := tracer.NewSpan(ctx, "MessageRepo.GetMessage", nil)
	defer span.End()

	var filterQuery string
	var args []interface{}

	if cursor > 0 {
		filterQuery = "AND message_id <= ?"
		args = []interface{}{channelID, cursor}
	} else {
		args = []interface{}{channelID}
	}

	sqlQuery := fmt.Sprintf(`
		SELECT channel_id, bucket, message_id, 
		content, user_id, username, created_at 
		FROM message WHERE channel_id = ? %s LIMIT %d ALLOW FILTERING`, filterQuery, limit)

	scanner := r.session.Query(sqlQuery, args...).WithContext(ctx).Iter().Scanner()
	messages := []*models.Message{}

	for scanner.Next() {
		var message models.Message
		err := scanner.Scan(&message.ChannelID, &message.Bucket, &message.MessageID, &message.Content, &message.UserID, &message.Username, &message.CreatedAt)
		if err != nil {
			tracer.AddSpanError(span, err)
			return nil, errors.Wrap(err, "MessageRepo.GetMessage")
		}
		messages = append(messages, &message)
	}

	if err := scanner.Err(); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "MessageRepo.Get")
	}

	return messages, nil
}

func (r *messageRepo) GetReply(ctx context.Context, messageID int64, cursor int64, limit int) ([]*models.Reply, error) {
	ctx, span := tracer.NewSpan(ctx, "MessageRepo.GetReplyet", nil)
	defer span.End()

	var filterQuery string
	var args []interface{}

	if cursor > 0 {
		filterQuery = "AND reply_id <= ?"
		args = []interface{}{messageID, cursor}
	} else {
		args = []interface{}{messageID}
	}

	sqlQuery := fmt.Sprintf(`
		SELECT message_id, reply_id,
		content, user_id, username, created_at 
		FROM reply WHERE message_id = ? %s LIMIT %d ALLOW FILTERING`, filterQuery, limit)

	// sqlQuery := `
	// 	SELECT message_id, reply_id,
	// 	content, user_id, username, created_at
	// 	FROM reply WHERE message_id = ? LIMIT 10 ALLOW FILTERING
	// `

	scanner := r.session.Query(sqlQuery, args...).WithContext(ctx).Iter().Scanner()
	// scanner := r.session.Query(sqlQuery, messageID).WithContext(ctx).Iter().Scanner()
	replies := []*models.Reply{}

	for scanner.Next() {
		var reply models.Reply
		err := scanner.Scan(&reply.MessageID, &reply.ReplyID, &reply.Content, &reply.UserID, &reply.Username, &reply.CreatedAt)
		if err != nil {
			tracer.AddSpanError(span, err)
			return nil, errors.Wrap(err, "MessageRepo.GetReply")
		}
		replies = append(replies, &reply)
	}

	if err := scanner.Err(); err != nil {
		tracer.AddSpanError(span, err)
		return nil, errors.Wrap(err, "MessageRepo.GetReply")
	}

	return replies, nil
}
