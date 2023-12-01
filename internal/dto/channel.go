package dto

type CreateChannelReqDto struct {
	Name string `json:"name" binding:"required"`
}

type JoinChannelReqDto struct {
	ChannelID string `json:"channel_id" binding:"required"`
}
