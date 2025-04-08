package user

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go-chat/internal/domain/models"
	"go-chat/internal/lib/logger/sl"
	"go-chat/internal/storage"
	"go-chat/internal/storage/redis"

	ssov1 "github.com/memxire/protobuf/gen/go/sso"
)

// Service Errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrProfileNotFound = storage.ErrNotFound
	ErrUsernameLookup  = errors.New("failed to find user by username")
	ErrProfileUpdate   = errors.New("failed to update profile")
	ErrCaching         = errors.New("caching operation failed")
	ErrInternal        = errors.New("internal service error")
	ErrSSODataMissing  = errors.New("required SSO data is missing")
)

// UserService defines methods for working with user data
type UserService interface {
	GetFullUserProfile(ctx context.Context, userID int64) (
		*FullUserProfileDTO, error)
	UpdateUserProfile(ctx context.Context, userID int64,
		req *UpdateProfileRequestDTO) (*FullUserProfileDTO, error)
	// TODO: implement for public profile view
	// GetPublicProfileByUsername(ctx context.Context, username string) (
	// *PublicUserProfileDTO, error)
}

// userService implements the UserService interface
type userService struct {
	log         *slog.Logger
	profileRepo storage.UserProfileRepository
	ssoClient   SSOUserClient
	redisClient RedisClient
	ssoCacheTTL time.Duration
	pgCacheTTL  time.Duration
}

type SSOUserClient interface {
	GetUser(ctx context.Context, userID int64) (*ssov1.GetUserResponse, error)
	// TODO: this method must be added to gRPC contract and client
	// GetUserIDByUsername(ctx context.Context, username string) (int64, error)
}
type RedisClient interface {
	GetSSOUserCache(ctx context.Context, userID int64) (*ssov1.GetUserResponse, error)
	SetSSOUserCache(ctx context.Context, userID int64, userData *ssov1.GetUserResponse, ttl time.Duration) error
	DelSSOUserCache(ctx context.Context, userID int64) error
	GetUserProfileCache(ctx context.Context, userID int64) (*models.UserProfile, error)
	SetUserProfileCache(ctx context.Context, userID int64, profile *models.UserProfile, ttl time.Duration) error
	DelUserProfileCache(ctx context.Context, userID int64) error
}

// NewUserService creates a new instance of UserService
func NewUserService(
	log *slog.Logger,
	profileRepo storage.UserProfileRepository,
	ssoClient SSOUserClient,
	redisClient RedisClient,
	ssoCacheTTL time.Duration,
	pgCacheTTL time.Duration,
) UserService {
	return &userService{
		log:         log,
		profileRepo: profileRepo,
		ssoClient:   ssoClient,
		redisClient: redisClient,
		ssoCacheTTL: ssoCacheTTL,
		pgCacheTTL:  pgCacheTTL,
	}
}

// getSSOUserData gets data from SSO with caching
func (s *userService) getSSOUserData(ctx context.Context, userID int64) (
	*ssov1.GetUserResponse, error) {
	const op = "service.user.getSSOUserData"
	log := s.log.With(slog.String("op", op), slog.Int64("user_id", userID))

	// Check Redis cache
	ssoData, err := s.redisClient.GetSSOUserCache(ctx, userID)
	if err == nil {
		log.Debug("cache hit for sso user data")
		return ssoData, nil
	}
	if !errors.Is(err, redis.ErrCacheMiss) {
		log.Error("failed to get sso user data from cache", sl.Err(err))
	} else {
		log.Debug("cache miss for sso user data")
	}

	// Cache Miss - making a request to the SSO gRPC service
	ssoResp, err := s.ssoClient.GetUser(ctx, userID)
	if err != nil {
		log.Error("failed to get user data from sso service", sl.Err(err))
		return nil, fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	// Cache the result
	errCache := s.redisClient.SetSSOUserCache(ctx, userID, ssoResp, s.ssoCacheTTL)
	if errCache != nil {
		log.Error("failed to set sso user data to cache", sl.Err(errCache))
	}

	return ssoResp, nil
}

// getPostgresProfile gets profile from PG with caching
func (s *userService) getPostgresProfile(ctx context.Context, userID int64) (
	*models.UserProfile, error) {
	const op = "service.user.getPostgresProfile"
	log := s.log.With(slog.String("op", op), slog.Int64("user_id", userID))

	// Check Redis cache for PG profile
	profile, err := s.redisClient.GetUserProfileCache(ctx, userID)
	if err == nil {
		log.Debug("cache hit for pg profile data")
		return profile, nil
	}
	if !errors.Is(err, redis.ErrCacheMiss) {
		log.Error("failed to get pg profile from cache", sl.Err(err))
	} else {
		log.Debug("cache miss for pg profile data")
	}

	// Cache Miss - Making a request to the PG repository
	profile, err = s.profileRepo.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			log.Debug("profile not found in postgres")
			return nil, ErrProfileNotFound
		}
		log.Error("failed to get profile from repository", sl.Err(err))
		return nil, ErrInternal
	}

	// Cache the result
	errCache := s.redisClient.SetUserProfileCache(ctx, userID, profile,
		s.pgCacheTTL)
	if errCache != nil {
		log.Error("failed to set pg profile data to cache", sl.Err(errCache))
	}

	return profile, nil
}

// GetFullUserProfile gets the combined user data
func (s *userService) GetFullUserProfile(ctx context.Context, userID int64) (
	*FullUserProfileDTO, error) {
	const op = "service.user.GetFullUserProfile"
	log := s.log.With(slog.String("op", op), slog.Int64("user_id", userID))

	ssoData, errSso := s.getSSOUserData(ctx, userID)
	if errSso != nil {
		log.Error("failed to get required sso data", sl.Err(errSso))
		if errors.Is(errSso, ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, ErrSSODataMissing
	}

	pgProfile, errPg := s.getPostgresProfile(ctx, userID)
	if errPg != nil && !errors.Is(errPg, ErrProfileNotFound) {
		log.Error("failed to get postgres profile data, proceeding without it",
			sl.Err(errPg))
		pgProfile = nil
	}

	fullProfile := &FullUserProfileDTO{
		UserID:   userID,
		Email:    ssoData.GetEmail(),
		Username: ssoData.GetUsername(),
	}
	if pgProfile != nil {
		fullProfile.FirstName = pgProfile.FirstName
		fullProfile.LastName = pgProfile.LastName
		fullProfile.AvatarURL = pgProfile.AvatarURL
		fullProfile.CreatedAt = &pgProfile.CreatedAt
		fullProfile.UpdatedAt = &pgProfile.UpdatedAt
	}

	return fullProfile, nil
}

// UpdateUserProfile updates the profile and invalidates the cache
func (s *userService) UpdateUserProfile(ctx context.Context, userID int64,
	req *UpdateProfileRequestDTO) (*FullUserProfileDTO, error) {
	const op = "service.user.UpdateUserProfile"
	log := s.log.With(slog.String("op", op), slog.Int64("user_id", userID))

	// TODO: Validation of data from req (string length, URL format, etc.)

	profileDataToSave := &models.UserProfile{
		UserID:    userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		AvatarURL: req.AvatarURL,
	}

	// Call repository for saving
	err := s.profileRepo.Save(ctx, profileDataToSave)
	if err != nil {
		log.Error("failed to save profile via repository", sl.Err(err))
		return nil, ErrProfileUpdate
	}

	// PG profile cache invalidation
	err = s.redisClient.DelUserProfileCache(ctx, userID)
	if err != nil {
		log.Error("failed to delete pg profile cache after update", sl.Err(err))
	} else {
		log.Info("pg profile cache invalidated after update")
	}

	// Receive current full data for the answer after the update
	fullProfileDTO, err := s.GetFullUserProfile(ctx, userID)
	if err != nil {
		log.Error("failed to get full user profile after update", sl.Err(err),
			slog.Int64("user_id", userID))
		return nil, fmt.Errorf("%s: failed to retrieve full profile post-update: %w", op, err)
	}

	return fullProfileDTO, nil
}

// // GetPublicProfileByUsername gets the public profile
// func (s *userService) GetPublicProfileByUsername(ctx context.Context, username string) (*PublicUserProfileDTO, error) {
// 	const op = "service.user.GetPublicProfileByUsername"
// 	log := s.log.With(slog.String("op", op), slog.String("username", username))

// 	// Find userID by username via SSO Client
// 	targetUserID, err := s.ssoClient.GetUserIDByUsername(ctx, username)
// 	if err != nil {
// 		log.Warn("failed to get user id by username from sso", sl.Err(err))
// 		return nil, ErrUserNotFound
// 	}
// 	log = log.With(slog.Int64("target_user_id", targetUserID))

// 	// Get PG profile data with cache
// 	pgProfile, errPg := s.getPostgresProfile(ctx, targetUserID)
// 	if errPg != nil && !errors.Is(errPg, ErrProfileNotFound) {
// 		log.Error("failed to get postgres profile data for public view",
// 			sl.Err(errPg))
// 		return nil, ErrInternal
// 	}

// 	publicDTO := &PublicUserProfileDTO{
// 		UserID:   targetUserID,
// 		Username: username,
// 	}
// 	if pgProfile != nil {
// 		publicDTO.FirstName = pgProfile.FirstName
// 		publicDTO.LastName = pgProfile.LastName
// 		publicDTO.AvatarURL = pgProfile.AvatarURL
// 	}

// 	return publicDTO, nil
// }
