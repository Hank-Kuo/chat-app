package dto

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

type GetMessageReqDto struct {
	ChannelID string `json:"channel_id" binding:"required"`
}

type GetReplyReqDto struct {
	MessageID int64 `json:"message_id" binding:"required"`
}
