package models

import "time"

type Chat struct {
	ChatID    int64     `db:"chat_id"`
	ChatType  string    `db:"chat_type"` // "direct" or "group"
	CreatedAt time.Time `db:"created_at"`
	// Name *string `db:"name"` // TODO: add for group chats
}
