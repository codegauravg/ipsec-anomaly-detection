package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSClient struct {
	client   *sqs.Client
	queueURL string
}

func NewSQSClient(ctx context.Context, queueURL string) (*SQSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		return nil, err
	}

	return &SQSClient{
		client:   sqs.NewFromConfig(cfg),
		queueURL: queueURL,
	}, nil
}

func (s *SQSClient) ReceiveMessages(ctx context.Context, maxMessages int32, waitTime int32) ([]types.Message, error) {
	msgResult, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTime,
	})

	if err != nil {
		return nil, err
	}

	return msgResult.Messages, nil
}

func (s *SQSClient) DeleteMessage(ctx context.Context, receiptHandle *string) {
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueURL),
		ReceiptHandle: receiptHandle,
	})
	if err != nil {
		log.Printf("Failed to delete message: %v", err)
	}
}
