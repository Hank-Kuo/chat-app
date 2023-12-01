package models

import "time"

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	LoginTime time.Time `db:"login_time"`
}
