package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go-chat/internal/domain/models"

	ssov1 "github.com/memxire/protobuf/gen/go/sso"
	"github.com/redis/go-redis/v9"
)

// Keys for redis cache errors
var ErrCacheMiss = errors.New("cache: key not found")
var ErrCacheDataCorrupted = errors.New("cache: corrupted data")

type Client struct {
	Client *redis.Client
}

// New returns a new Redis client.
func New(addr string, password string, db int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &Client{Client: rdb}
}

// getData gets the data by key and deserializes it into a new instance of type T.
func getData[T any](r *Client, ctx context.Context, key string) (*T, error) {
	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}
		return nil, fmt.Errorf("redis GET error for key %s: %w", key, err)
	}

	var target T
	if err = json.Unmarshal([]byte(data), &target); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCacheDataCorrupted, err)
	}
	return &target, nil
}

// setDataJSON serializes data (of type T) to JSON and stores by key with TTL.
func setData[T any](r *Client, ctx context.Context, key string, data *T,
	ttl time.Duration) error {
	if data == nil {
		return errors.New("cannot cache nil generic data")
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal generic data for cache (key: %s): %w", key, err)
	}

	err = r.Client.Set(ctx, key, bytes, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis SET error for key %s: %w", key, err)
	}
	return nil
}

func (r *Client) delData(ctx context.Context, key string) error {
	err := r.Client.Del(ctx, key).Err()
	if err != nil && !errors.Is(err, redis.Nil) {
		return fmt.Errorf("redis DEL error for key %s: %w", key, err)
	}
	return nil
}

// Public methods for caching SSO data
func ssoUserCacheKey(userID int64) string {
	return fmt.Sprintf("user_sso:%d", userID)
}

// GetSSOUserCache gets the cached user data from SSO.
// Returns ErrCacheMiss if the key is not found.
func (r *Client) GetSSOUserCache(ctx context.Context, userID int64) (
	*ssov1.GetUserResponse, error) {
	return getData[ssov1.GetUserResponse](r, ctx, ssoUserCacheKey(userID))
}

// SetSSOUserCache saves user data from SSO to cache.
func (r *Client) SetSSOUserCache(ctx context.Context, userID int64,
	userData *ssov1.GetUserResponse, ttl time.Duration) error {
	return setData(r, ctx, ssoUserCacheKey(userID), userData, ttl)
}

// DelSSOUserCache deletes the SSO user data cache.
func (r *Client) DelSSOUserCache(ctx context.Context, userID int64) error {
	return r.delData(ctx, ssoUserCacheKey(userID))
}

// Public methods for Caching PG Profile Data
func pgProfileCacheKey(userID int64) string {
	return fmt.Sprintf("user_profile:%d", userID)
}

// GetUserProfile Cache gets the cached user profile from PG.
// Returns ErrCacheMiss if the key is not found.
func (r *Client) GetUserProfileCache(ctx context.Context, userID int64) (
	*models.UserProfile, error) {
	return getData[models.UserProfile](r, ctx, pgProfileCacheKey(userID))
}

// SetUserProfileCache saves the user profile from PG to cache.
func (r *Client) SetUserProfileCache(ctx context.Context, userID int64,
	profile *models.UserProfile, ttl time.Duration) error {
	return setData(r, ctx, pgProfileCacheKey(userID), profile, ttl)
}

// DelUserProfile Cache deletes the PG user profile cache.
func (r *Client) DelUserProfileCache(ctx context.Context, userID int64) error {
	return r.delData(ctx, pgProfileCacheKey(userID))
}
