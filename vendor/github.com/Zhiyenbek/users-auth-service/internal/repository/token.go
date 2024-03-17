package repository

import (
	"errors"
	"fmt"

	"github.com/Zhiyenbek/users-auth-service/internal/models"
	"github.com/go-redis/redis/v7"
)

type tokenRepository struct {
	client *redis.Client
}

func NewTokenRepository(client *redis.Client) TokenRepository {
	return &tokenRepository{
		client: client,
	}
}
func (r *tokenRepository) SetRTToken(token *models.Token) error {
	key := token.PublicID
	if err := r.client.Set(key, token.TokenValue, token.TTL).Err(); err != nil {
		return fmt.Errorf("%w could not set refresh token to redis for TokenValue : %s: %v", models.ErrInternalServer, token.TokenValue, err)
	}
	return nil
}

func (r *tokenRepository) UnsetRTToken(publicID string) error {
	key := publicID
	if err := r.client.Del(key).Err(); err != nil {
		return fmt.Errorf("%w could not delete refresh token to redis for TokenValue : %v", models.ErrInternalServer, err)
	}
	return nil
}

func (r *tokenRepository) GetToken(publicID string) (string, error) {
	key := publicID
	value := r.client.Get(key)
	TokenValue, err := value.Result()
	if err != nil || TokenValue == "" {
		if errors.Is(err, redis.Nil) {
			return "", fmt.Errorf("token does not exist in storage: %w", models.ErrTokenExpired)
		}
		return "", fmt.Errorf("%w could not retrieve refresh token from redis for TokenValue : %v", models.ErrInternalServer, err)
	}
	return TokenValue, nil
}
