package main

import (
	"context"
	"insider/internal/worker"

	_ "github.com/lib/pq"

	"insider/api"
	_ "insider/docs"
	"insider/internal/cache"
	"insider/internal/repository"
	"insider/internal/service"
	"insider/internal/webhook"
	"log/slog"
	"net/http"
	"os"
)

//const (
//	redisAddr  = "127.0.0.1:6379"
//	postgreURL = "user=postgres password=1234 dbname=mydb sslmode=disable"
//	webhookURL = "https://webhook.site/d90c7e79-8f22-4876-9f0b-5d9f2b20fb32"
//)

var (
	redisAddr  = os.Getenv("REDIS_ADDR")
	postgreURL = os.Getenv("POSTGRES_URL")
	webhookURL = os.Getenv("WEBHOOK_URL")
)

// @title			Message Sender API
// @version			1.0
// @description     This is a sample server for automatically sending messages.
// @host			localhost:8080
// @BasePath		/
func main() {
	redisClient := cache.New(cache.Config{
		Addr:     redisAddr,
		Password: "",
		DB:       0,
	})
	err := redisClient.Connect()
	if err != nil {
		slog.Error("failed to connect to redis", slog.Any("error", err))
		os.Exit(1)
	}

	db, err := repository.New(postgreURL)
	if err != nil {
		slog.Error("failed to connect to postgres", slog.Any("error", err))
		os.Exit(1)
	}

	webhookClient := webhook.New(webhookURL)
	msgSvc := service.NewMessageService(db, redisClient, webhookClient)
	w := worker.NewWorker(redisClient, db)

	ctx := context.Background()

	go func() {
		for {
			err = w.Run(ctx)
			if err != nil {
				slog.Error("worker stopped with error", slog.Any("error", err))
			}
		}
	}()

	go func() {
		msgSvc.Run(ctx)
	}()

	a := api.New(msgSvc)
	router := a.SetupRoutes()
	slog.Info("http server started", slog.String("addr", ":8080"))
	if err := http.ListenAndServe(":8080", router); err != nil {
		slog.Error("http server failed", slog.Any("error", err))
		os.Exit(1)
	}
}
