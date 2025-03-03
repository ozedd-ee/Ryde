package data

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type TokenStore struct {
	RedisClient *redis.Client
}

func NewTokenStore(redisClient *redis.Client) *TokenStore {
	return &TokenStore{
		RedisClient: redisClient,
	}
}

func (s *TokenStore) UpdateFCMToken(ctx context.Context, ownerID, token string) error {
	_, err := s.RedisClient.Set(ctx, ownerID, token, 0).Result()
	if err != nil {
		return fmt.Errorf("failed to update FCM token for %v", ownerID)
	}
	return nil
}

func (s *TokenStore) GetFCMToken(ctx context.Context, ownerID string) (string, error) {
	token, err := s.RedisClient.Get(ctx, ownerID).Result()
	if err != nil {
		return "", fmt.Errorf("no FCM token stored for %v", ownerID)
	}
	return token, nil
}
