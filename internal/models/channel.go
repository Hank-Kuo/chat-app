package models

import "time"

type Channel struct {
	ID        string    `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	UserID    string    `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type UserToChannel struct {
	UserID    string    `db:"user_id"`
	ChannelID string    `db:"channel_id"`
	CreatedAt time.Time `db:"created_at"`
}
