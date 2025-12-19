// Package service provides business logic for message sending and updating operations.
package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"insider/internal/cache"
	"log/slog"
	"sync"
	"time"

	"insider/internal/repository"
	"insider/internal/webhook"
)

var (
	messageLimit = 2
)

type MessageService struct {
	Store         MessageStore
	Cache         CacheRepository
	WebhookClient WebhookClient

	running bool
	mu      sync.Mutex
}

type MessageStore interface {
	Fetch(ctx context.Context, limit int) ([]repository.Message, error)
	GetSendingMessages(ctx context.Context) ([]repository.Message, error)
	UpdateStatus(ctx context.Context, id string, status repository.MessageStatus) (*sql.Tx, error)
}

type CacheRepository interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	GetMessages(ctx context.Context) ([]cache.SentMessage, error)
	Del(ctx context.Context, key string) error
}

type WebhookClient interface {
	Send(phone, content string) (*webhook.MessageResponse, error)
}

func NewMessageService(store MessageStore, cacheClient CacheRepository, webhookClient WebhookClient) *MessageService {
	return &MessageService{
		Store:         store,
		Cache:         cacheClient,
		WebhookClient: webhookClient,
	}
}

func (s *MessageService) Fetch(ctx context.Context) ([]repository.Message, error) {
	return s.Store.Fetch(ctx, messageLimit)
}

func (s *MessageService) process() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	messages, err := s.Fetch(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch messages: %w", err)
	}

	for _, msg := range messages {
		if msg.Status == repository.StatusSent {
			continue
		}

		txS, err := s.Store.UpdateStatus(ctx, msg.ID, repository.StatusSending)
		if err != nil {
			return fmt.Errorf("failed to update status to StatusSending: %w", err)
		}

		response, err := s.WebhookClient.Send(msg.PhoneNumber, msg.Content)
		if err != nil {
			slog.Error("failed to send message to webhook", slog.Any("error", err))
			err = txS.Rollback()
			if err != nil {
				return fmt.Errorf("failed to StatusSending rollback transaction: %w", err)
			}
			continue
		}
		err = txS.Commit()
		if err != nil {
			slog.Error("failed to StatusSending commit transaction", slog.Any("error", err))
		}

		tx, err := s.Store.UpdateStatus(ctx, msg.ID, repository.StatusSent)
		if err != nil {
			return fmt.Errorf("failed to update status to StatusSent: %w", err)
		}

		sc := cache.SentMessage{
			MessageID:         msg.ID,
			PhoneNumber:       msg.PhoneNumber,
			Content:           msg.Content,
			Status:            repository.StatusSent,
			CreatedAt:         msg.CreatedAt,
			ResponseMessageID: response.MessageID,
			SentAt:            time.Now(),
		}

		data, err := json.Marshal(sc)
		if err != nil {
			return fmt.Errorf("failed to marshal message: %w", err)
		}

		if err = s.Cache.Set(ctx, "sent_message:"+msg.ID, data, 0); err != nil {
			// Cache write failures are tolerated.
			// The worker periodically syncs sending state from the database, ensuring consistency.
			slog.Error("failed to cache message", slog.Any("error", err))
			err = tx.Rollback()
			if err != nil {
				return fmt.Errorf("failed to StatusSent rollback transaction: %w", err)
			}
		}
		err = tx.Commit()
		if err != nil {
			return fmt.Errorf("failed to StatusSent commit transaction: %w", err)
		}
		slog.Info(
			"message sent successfully",
			slog.String("message_id", msg.ID),
			slog.String("response_message_id", response.MessageID),
		)
	}
	return nil
}

func (s *MessageService) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			slog.Info("Worker Shutting Down...")
			return
		case <-ticker.C:
			if s.running {
				err := s.process()
				if err != nil {
					slog.Error("failed to process message", slog.Any("error", err))
					return
				}
			}
		}
	}
}

func (s *MessageService) Start() error {
	s.mu.Lock()
	s.running = true
	s.mu.Unlock()
	slog.Info("Worker Started...")
	return nil
}

func (s *MessageService) Stop() error {
	s.mu.Lock()
	s.running = false
	s.mu.Unlock()
	slog.Info("Worker Stopped...")
	return nil
}
