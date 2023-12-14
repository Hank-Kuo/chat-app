package dto

import (
	"chat-app/internal/models"
)

type CreateMessageReqDto struct {
	ChannelID string `json:"channel_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type CreateReplyReqDto struct {
	MessageID int64  `json:"message_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type GetMessageQueryDto struct {
	ChannelID string `form:"channel_id" binding:"required"`
	Cursor    string `form:"cursor"`
}

type GetMessageResDto struct {
	Messages   []*models.Message `json:"messages"`
	NextCursor string            `json:"next_cursor"`
}
type GetReplyQueryDto struct {
	MessageID int64  `form:"message_id" binding:"required"`
	Cursor    string `form:"cursor"`
}

type GetReplyResDto struct {
	Replies    []*models.Reply `json:"replies"`
	NextCursor string          `json:"next_cursor"`
}
