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
	redisMessages, err := w.redis.GetMessages(ctx)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbMessages, err := w.store.GetSendingMessages(ctx)
	if err != nil {
		return err
	}

	redisMap := make(map[string]struct{}, len(redisMessages))
	for _, m := range redisMessages {
		redisMap[m.MessageID] = struct{}{}
	}

	for _, dbMsg := range dbMessages {
		if _, ok := redisMap[dbMsg.ID]; ok {
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
