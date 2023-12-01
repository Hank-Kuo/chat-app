package models

import "time"

type Message struct {
	ChannelID string    `json:"channel_id"`
	Bucket    int       `json:"_"`
	MessageID int64     `json:"message_id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type Reply struct {
	MessageID int64     `json:"message_id"`
	ReplyID   int64     `json:"reply_id"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
