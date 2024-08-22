package kafka

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/sethvargo/go-retry"
	"github.com/stretchr/testify/require"
)

const (
	AttemptCount = 3
	AttemptDelay = 1 * time.Second
)

var (
	brokers    = []string{"127.0.0.1:9091", "127.0.0.1:9092", "127.0.0.1:9093"}
	kafkaTopic = "events_test"
	kafkaGroup = "homework-one_test"

	bytesToSend = []byte{1, 2, 0, 5, 7}
	isReceived  = false
)

type TestReader struct{}

func (e *TestReader) Handle(msg []byte) error {
	isReceived = true

	if len(msg) != len(bytesToSend) {
		log.Fatal("error")
	}

	for i := 0; i < len(msg); i++ {
		if msg[i] != bytesToSend[i] {
			log.Fatal("error")
		}
	}

	return nil
}

func TestKafka(t *testing.T) {
	producer, err := NewProducer(brokers, kafkaTopic)
	require.NoError(t, err)

	testReader := TestReader{}

	consumer, err := NewConsumerGroup(brokers, kafkaGroup, &testReader)
	require.NoError(t, err)

	_, _, err = producer.SendSyncMessage(bytesToSend)
	require.NoError(t, err)

	go consumer.Run(kafkaTopic)

	b := retry.WithMaxRetries(
		AttemptCount,
		retry.NewExponential(AttemptDelay),
	)

	err = retry.Do(
		context.Background(),
		b,
		func(ctx context.Context) error {
			if isReceived != true {
				return retry.RetryableError(errors.New("repeat"))
			}

			return nil
		},
	)

	require.NoError(t, err)
	require.Equal(t, true, isReceived)
}
