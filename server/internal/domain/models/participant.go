package models

import "time"

type Participant struct {
	ParticipationID int64     `db:"participation_id"`
	ChatID          int64     `db:"chat_id"`
	UserID          int64     `db:"user_id"`
	JoinedAt        time.Time `db:"joined_at"`
}
