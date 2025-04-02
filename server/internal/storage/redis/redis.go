package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	ssov1 "github.com/memxire/protobuf/gen/go/sso"
	"github.com/redis/go-redis/v9"
)

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

func (r *Client) GetCache(ctx context.Context, userID int64) (
	*ssov1.GetUserResponse, error) {

	key := fmt.Sprintf("user:%d", userID)
	data, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, err
	}

	var userData *ssov1.GetUserResponse
	if err = json.Unmarshal([]byte(data), &userData); err != nil {
		return nil, err
	}

	return userData, nil
}

func (r *Client) SetCache(ctx context.Context, userID int64,
	userData *ssov1.GetUserResponse, ttl time.Duration) error {

	key := fmt.Sprintf("user:%d", userID)
	bytes, err := json.Marshal(userData)
	if err != nil {
		return err
	}

	return r.Client.Set(ctx, key, bytes, ttl).Err()
}
