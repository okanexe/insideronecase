// Package repository provides database access and repository implementations.
package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type PostgreStore struct {
	db *sql.DB
}

type MessageStatus string

const (
	StatusPending MessageStatus = "pending"
	StatusSending MessageStatus = "sending"
	StatusSent    MessageStatus = "sent"
)

type Message struct {
	ID          string        `json:"id"`
	PhoneNumber string        `json:"phone_number"`
	Content     string        `json:"content"`
	Status      MessageStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	SentAt      *time.Time    `json:"sent_at"`
}

func New(databaseURL string) (*PostgreStore, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &PostgreStore{db: db}, nil
}

func (p *PostgreStore) Fetch(ctx context.Context, limit int) ([]Message, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, phone_number, content, status, created_at, sent_at FROM messages WHERE status=$1 LIMIT $2", StatusPending, limit)
	if err != nil {
		return nil, fmt.Errorf("error fetching messages query: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.PhoneNumber, &msg.Content, &msg.Status, &msg.CreatedAt, &msg.SentAt); err != nil {
			return nil, fmt.Errorf("error Fetch scanning row: %w", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (p *PostgreStore) GetSendingMessages(ctx context.Context) ([]Message, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, phone_number, content, status, created_at, sent_at FROM messages WHERE status='sending'")
	if err != nil {
		return nil, fmt.Errorf("error get sent messages query: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.PhoneNumber, &msg.Content, &msg.Status, &msg.CreatedAt, &msg.SentAt); err != nil {
			return nil, fmt.Errorf("error GetSentMessages scanning row: %w", err)
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (p *PostgreStore) UpdateStatus(ctx context.Context, id string, status MessageStatus) (*sql.Tx, error) {
	tx, err := p.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error UpdateStatus beginning transaction: %w", err)
	}
	res, err := tx.ExecContext(ctx, `
		UPDATE messages
		SET status = $1,
		    sent_at = $2
		WHERE id = $3
	`, status, time.Now(), id)
	if err != nil {
		return nil, fmt.Errorf("error UpdateStatus exec context: %w", err)
	}

	i, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error UpdateStatus get rows affected: %w", err)
	}
	if i == 0 {
		return nil, fmt.Errorf("no rows affected by UpdateStatus update")
	}
	return tx, nil
}
