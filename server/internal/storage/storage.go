package storage

import (
	"context"
	"errors"
	"go-chat/internal/domain/models"
)

var ErrNotFound = errors.New("not found")

// UserProfileRepository defines methods for working with user profiles
type UserProfileRepository interface {
	// Get returns profile by user id. Returns ErrNotFound, if profile not found.
	Get(ctx context.Context, userID int64) (*models.UserProfile, error)
	// Save creates or updates user profile.
	Save(ctx context.Context, profile *models.UserProfile) error
}

// ChatRepository defines methods for working with chats and participants
type ChatRepository interface {
	// CreateChat creates a new chat and adds its creator as a member.
	// Returns the ID of the created chat.
	CreateChat(ctx context.Context, chatType string, creatorUserID int64,
		participantUserIDs []int64) (int64, error)
	// GetUserChats returns a list of chats (with the latest message?) the
	// user is a participant in.
	GetUserChats(ctx context.Context, userID int64, limit, offset int) (
		[]models.Chat, error)
	// GetChatParticipants returns a list of chat participant IDs.
	GetChatParticipants(ctx context.Context, chatID int64) ([]int64, error)
	// IsUserInChat checks if the user is a member of the chat.
	IsUserInChat(ctx context.Context, userID int64, chatID int64) (bool, error)
	// AddUserToChat adds a user to the chat.
	AddUserToChat(ctx context.Context, userID int64, chatID int64) error
	// TODO: Add methods to delete a member, get chat info, etc.
}

// MessageRepository defines methods for working with messages
type MessageRepository interface {
	// Save saves a new message in the chat. Returns the message ID.
	Save(ctx context.Context, msg *models.Message) (int64, error)
	// GetMessages returns a list of messages for the chat with pagination.
	GetMessages(ctx context.Context, chatID int64, limit, offset int) (
		[]models.Message, error)
	// Update updates message text.
	Update(ctx context.Context, messageID int64, userID int64,
		newContent string) error
	// Delete deletes message.
	Delete(ctx context.Context, messageID int64, userID int64) error
}
