package models

import "time"

type Message struct {
	MessageID    int64     `db:"message_id"`
	ChatID       int64     `db:"chat_id"`
	SenderUserID int64     `db:"sender_user_id"`
	Content      string    `db:"content"`
	SentAt       time.Time `db:"sent_at"`
	// MessageType  string    `db:"message_type"`
}
