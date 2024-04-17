package mq

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

// Consumer will be called while received messages
type Consumer func(ctx context.Context, client *redis.Client, topic, group string, id string, value any) error

// Queue define a set of methods that message queue handler should implement
type Queue interface {
	// Publish publishes a message to the specified topic
	Publish(ctx context.Context, topic string, msg any, maxLen int64) (id string, err error)
	// Consume register a consumer with callback
	Consume(ctx context.Context, topic, group, consumer string, batchSize int64, cb Consumer) error
}
