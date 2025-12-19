// Package cache provides a Redis client wrapper for caching functionality.
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"insider/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

type Repository struct {
	redis *redis.Client
}

func New(conf Config) *Repository {
	return &Repository{
		redis: redis.NewClient(&redis.Options{
			Addr:     conf.Addr,
			Password: conf.Password,
			DB:       conf.DB,
		}),
	}
}

func (r *Repository) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.redis.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

func (r *Repository) Exists(ctx context.Context, key string) (bool, error) {
	_, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.redis.Set(ctx, key, value, expiration).Err()
}

type SentMessage struct {
	MessageID         string                   `json:"message_id"`
	ResponseMessageID string                   `json:"response_message_id,omitempty"`
	PhoneNumber       string                   `json:"phone_number"`
	Content           string                   `json:"content"`
	Status            repository.MessageStatus `json:"status"`
	CreatedAt         time.Time                `json:"created_at"`
	SentAt            time.Time                `json:"sent_at"`
}

func (r *Repository) GetMessages(ctx context.Context) ([]SentMessage, error) {
	var result []SentMessage

	iter := r.redis.Scan(ctx, 0, "sent_message:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		val, err := r.redis.Get(ctx, key).Result()
		if err != nil {
			// The key may have expired
			if err == redis.Nil {
				continue
			}
			return nil, err
		}

		var msg SentMessage
		if err := json.Unmarshal([]byte(val), &msg); err != nil {
			return nil, err
		}

		result = append(result, msg)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *Repository) Del(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}
