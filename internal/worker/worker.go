package worker

import (
	"context"
	"encoding/json"
	"insider/internal/cache"
	"insider/internal/repository"
	"log/slog"
	"time"
)

type Worker struct {
	redis *cache.Repository
	store *repository.PostgreStore
}

func NewWorker(r *cache.Repository, s *repository.PostgreStore) *Worker {
	return &Worker{
		redis: r,
		store: s,
	}
}

func (w *Worker) Run(ctx context.Context) error {
	dbMessages, err := w.store.GetSendingMessages(ctx)
	if err != nil {
		return err
	}

	for _, dbMsg := range dbMessages {
		exists, err := w.redis.Exists(ctx, "sent_message:"+dbMsg.ID)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		msg := cache.SentMessage{
			MessageID:   dbMsg.ID,
			PhoneNumber: dbMsg.PhoneNumber,
			Content:     dbMsg.Content,
			Status:      dbMsg.Status,
			CreatedAt:   dbMsg.CreatedAt,
			SentAt:      time.Now(),
		}
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		if err := w.redis.Set(ctx, "sent_message:"+dbMsg.ID, data, 0); err != nil {
			return err
		}

		tx, err := w.store.UpdateStatus(ctx, dbMsg.ID, repository.StatusSent)
		if err != nil {
			return err
		}
		err = tx.Commit()
		if err != nil {
			slog.Error("failed to worker commit StatusSent", slog.Any("error", err))
		}
	}
	return nil
}
