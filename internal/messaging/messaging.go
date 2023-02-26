package messaging

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"os"
)

type messageProducer = sarama.SyncProducer

type MessageQueue struct {
	Producer messageProducer
}

func NewMessageQueue() (*MessageQueue, error) {
	if os.Getenv("KAFKA_ACTIVE") != "" {
		brokersUrl := os.Getenv("BROKERS_URL")
		conn, err := connectProducer([]string{brokersUrl})
		if err != nil {
			return nil, err
		}
		return &MessageQueue{Producer: conn}, nil
	}
	return &MessageQueue{}, nil
}

func connectProducer(brokersUrl []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	// NewSyncProducer creates a new SyncProducer using the given broker addresses and configuration.
	conn, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create a sync producer: %v", err)
	}
	return conn, nil
}

func (mq *MessageQueue) PushToQueue(ctx context.Context, topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := mq.Producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to push to queue: %v", err)
	}
	log.Debugf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
