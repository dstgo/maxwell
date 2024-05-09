package mq

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log/slog"
	"time"
)

func NewStreamQueue(client *redis.Client) *StreamQueue {
	return &StreamQueue{redis: client}
}

// StreamQueue implement Queue interface by Redis Stream
type StreamQueue struct {
	redis *redis.Client
}

func (h *StreamQueue) Publish(ctx context.Context, topic string, msg any, maxLen int64) (id string, err error) {
	result, err := h.redis.XAdd(ctx, &redis.XAddArgs{
		Stream: topic,
		MaxLen: maxLen,
		Values: msg,
		ID:     "*",
	}).Result()

	return result, err
}

func (h *StreamQueue) Consume(ctx context.Context, topic, group, consumer string, batchSize int64, cb Consumer) error {
	stream := h.redis.XGroupCreateMkStream(ctx, topic, group, "0")
	if stream.Err() != nil && stream.Err().Error() != "BUSYGROUP Consumer Group name already exists" {
		return stream.Err()
	}

	go func() {
		slog.Debug(fmt.Sprintf("%s is running", consumer), slog.String("topic", topic), slog.String("group", group))
		defer func() {
			if err := recover(); err != nil {
				slog.Error("stream panic recovered", slog.Any("error", err))
			}
		}()

		for {

			// read the latest message
			if id, err := h.consume(ctx, topic, group, consumer, ">", batchSize, cb); err != nil {
				errorLog("stream consume > failed", err, id, topic, group, consumer)
			}

			// consume the messages that already received but not ack yet
			if id, err := h.consume(ctx, topic, group, consumer, "1", batchSize, cb); err != nil {
				errorLog("stream consume 1 failed", err, id, topic, group, consumer)
			}

			// clear dead messages in pending list
			if err := h.clearDead(ctx, topic, group, time.Minute*5, 10); err != nil {
				slog.Error("stream clear dead failed", slog.Any("error", err))
				return
			}

			time.Sleep(1)
		}
	}()
	return nil
}

func (h *StreamQueue) consume(ctx context.Context, topic, group, consumer, id string, batchSize int64, cb Consumer) (errorId string, err error) {
	// read from specified stream in specified group
	result, err := h.redis.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{topic, id},
		Count:    batchSize,
	}).Result()

	if err != nil {
		return "", err
	}

	for _, stream := range result {
		topic := stream.Stream
		for _, message := range stream.Messages {
			if err := cb(ctx, h.redis, topic, group, message.ID, message.Values); err != nil {
				return message.ID, err
			} else { // make sure message is consumed if callback executed successfully
				if err := h.redis.XAck(ctx, topic, group, message.ID).Err(); err != nil {
					return message.ID, err
				}

				// del it if ack ok
				if err := h.redis.XDel(ctx, topic, message.ID).Err(); err != nil {
					return message.ID, err
				}
			}
		}
	}

	return "", nil
}

// clear dead msg that idle timeout
func (h *StreamQueue) clearDead(ctx context.Context, topic, group string, idle time.Duration, count int64) error {
	pel, err := h.redis.XPendingExt(ctx, &redis.XPendingExtArgs{
		Stream: topic,
		Group:  group,
		Idle:   idle,
		Start:  "-",
		End:    "+",
		Count:  count,
	}).Result()

	if err != nil {
		return err
	}

	var ids []string
	for _, pending := range pel {
		ids = append(ids, pending.ID)
	}

	if len(ids) == 0 {
		return nil
	}

	/// delete msg
	if _, err := h.redis.XDel(ctx, topic, ids...).Result(); err != nil {
		return err
	}

	return nil
}

func errorLog(msg string, err error, id, topic, group, consumer string) {
	slog.Error(msg,
		slog.String("error", err.Error()),
		slog.String("msg-id", id),
		slog.String("topic", topic),
		slog.String("group", group),
		slog.String("consumer", consumer),
	)
}
