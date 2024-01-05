package databus_kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/art-es/blog/internal/auth/dto"
)

type Client struct {
	activationEmailWriter *kafka.Writer
}

func New(kafkaURL string) *Client {
	return &Client{
		activationEmailWriter: &kafka.Writer{
			Addr:     kafka.TCP(kafkaURL),
			Topic:    "auth.activation_codes",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (c *Client) ProduceActivationEmail(ctx context.Context, msg *dto.UserActivationEmailMessage) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}

	err = c.activationEmailWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte("send_email"),
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("write message to kafka error: %w", err)
	}

	return nil
}
