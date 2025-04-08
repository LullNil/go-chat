package user

import "time"

// FullUserProfileDTO combines data from SSO and local profile
type FullUserProfileDTO struct {
	UserID    int64      `json:"user_id"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	FirstName *string    `json:"first_name,omitempty"`
	LastName  *string    `json:"last_name,omitempty"`
	AvatarURL *string    `json:"avatar_url,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// PublicUserProfileDTO contains only public data
type PublicUserProfileDTO struct {
	UserID    int64   `json:"user_id"`
	Username  string  `json:"username"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	AvatarURL *string `json:"avatar_url,omitempty"`
}

// UpdateProfileRequestDTO structure for profile update request received by handler
type UpdateProfileRequestDTO struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	AvatarURL *string `json:"avatar_url"`
}
