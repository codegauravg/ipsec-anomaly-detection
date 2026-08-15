package main

import (
	"context"
	"log"
	"os"

	"github.com/codegauravg/ipsec-anomaly-worker/internal/aws"
	"github.com/codegauravg/ipsec-anomaly-worker/internal/cache"
	"github.com/codegauravg/ipsec-anomaly-worker/internal/db"
	"github.com/codegauravg/ipsec-anomaly-worker/internal/models"
	"github.com/codegauravg/ipsec-anomaly-worker/internal/processor"
)

func main() {
	ctx := context.Background()

	// 1. Load Config
	cfg := models.Config{
		QueueURL:    os.Getenv("SQS_QUEUE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	// 2. Init Clients
	sqsClient, err := aws.NewSQSClient(ctx, cfg.QueueURL)
	if err != nil {
		log.Fatalf("Failed to init SQS: %v", err)
	}

	redisCache := cache.NewCache(cfg.RedisURL)

	dbStore, err := db.NewStore(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer dbStore.Close()

	// 3. Start Processor
	proc := processor.New(sqsClient, redisCache, dbStore)
	proc.Start(ctx)
}
