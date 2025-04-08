package models

import "time"

type UserProfile struct {
	UserID    int64     `db:"user_id"`
	FirstName *string   `db:"first_name"` // Pointers to NULLable fields
	LastName  *string   `db:"last_name"`
	AvatarURL *string   `db:"avatar_url"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
