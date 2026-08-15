package processor

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/codegauravg/ipsec-anomaly-worker/internal/aws"
	"github.com/codegauravg/ipsec-anomaly-worker/internal/cache"
	"github.com/codegauravg/ipsec-anomaly-worker/internal/db"
	"github.com/codegauravg/ipsec-anomaly-worker/internal/models"
)

const (
	Threshold     = 50
	WindowSeconds = 60 * time.Second
)

type Processor struct {
	sqs   *aws.SQSClient
	cache *cache.Cache
	db    *db.Store
}

func New(sqs *aws.SQSClient, cache *cache.Cache, db *db.Store) *Processor {
	return &Processor{sqs: sqs, cache: cache, db: db}
}

func (p *Processor) Start(ctx context.Context) {
	log.Println("Starting Anomaly Detection Processor...")

	for {
		messages, err := p.sqs.ReceiveMessages(ctx, 10, 20)
		if err != nil {
			log.Printf("Error fetching messages: %v", err)
			time.Sleep(5 * time.Second) // Backoff
			continue
		}

		if len(messages) == 0 {
			continue
		}

		var wg sync.WaitGroup
		for _, msg := range messages {
			wg.Add(1)
			go func(m types.Message) {
				defer wg.Done()
				p.handleMessage(ctx, m)
			}(msg)
		}

		wg.Wait()
	}
}

func (p *Processor) handleMessage(ctx context.Context, msg types.Message) {
	var ipsecLog models.IPSecLog

	if err := json.Unmarshal([]byte(*msg.Body), &ipsecLog); err != nil {
		log.Printf("Invalid JSON: %v", err)
		return
	}

	if ipsecLog.Action != "IKE_NEGOTIATION_FAILED" {
		p.sqs.DeleteMessage(ctx, msg.ReceiptHandle)
		return
	}

	redisKey := "ipsec:anomaly:" + ipsecLog.SourceIP

	count, err := p.cache.IncrementAndExpire(ctx, redisKey, WindowSeconds)
	if err != nil {
		log.Printf("Cache error for %s: %v", ipsecLog.SourceIP, err)
		return
	}

	if count == Threshold {
		log.Printf("[ALERT] Threshold breached for IP: %s", ipsecLog.SourceIP)

		err := p.db.InsertAnomaly(ctx, ipsecLog.SourceIP, "BRUTE_FORCE_IKE", *msg.Body, 8)
		if err != nil {
			log.Printf("DB insert failed: %v", err)
			return // Don't delete from SQS, allow retry
		}
	}

	p.sqs.DeleteMessage(ctx, msg.ReceiptHandle)
}
