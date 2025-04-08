package postgres

import (
	"context"
	"errors"
	"fmt"
	"go-chat/internal/domain/models"
	"go-chat/internal/storage"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

// userProfileRepo implements the UserProfileRepository interface
type userProfileRepo struct {
	pool *pgxpool.Pool
}

// NewUserProfileRepository creates a new instance of the profile repository
func NewUserProfileRepository(pool *pgxpool.Pool) storage.UserProfileRepository {
	return &userProfileRepo{pool: pool}
}

// Get gets the user profile by ID.
func (r *userProfileRepo) Get(ctx context.Context, userID int64) (
	*models.UserProfile, error) {
	const op = "storage.postgres.userProfileRepo.Get"

	query := `SELECT user_id, first_name, last_name, avatar_url, created_at, updated_at
		FROM users_profiles WHERE user_id = $1`

	profile := &models.UserProfile{}

	err := r.pool.QueryRow(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.AvatarURL,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, storage.ErrNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return profile, nil
}

// Save saves or updates user profile by ID.
func (r *userProfileRepo) Save(ctx context.Context,
	profile *models.UserProfile) error {
	const op = "storage.postgres.userProfileRepo.Save"

	query := `
		INSERT INTO users_profiles (user_id, first_name, last_name, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (user_id) DO UPDATE SET
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			avatar_url = EXCLUDED.avatar_url,
			updated_at = NOW()`

	_, err := r.pool.Exec(ctx, query,
		profile.UserID,
		profile.FirstName,
		profile.LastName,
		profile.AvatarURL,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
