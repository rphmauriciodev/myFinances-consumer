package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Queue struct {
	client   *sqs.Client
	queueURL string
}

func NewQueue(ctx context.Context, queueURL string) (*Queue, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(cfg)

	return &Queue{
		client:   client,
		queueURL: queueURL,
	}, nil
}

func (q *Queue) ReceiveMessages(ctx context.Context, maxMessages int32) ([]types.Message, error) {
	output, err := q.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.queueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		return nil, err
	}

	return output.Messages, nil
}

func (q *Queue) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := q.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(q.queueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}
