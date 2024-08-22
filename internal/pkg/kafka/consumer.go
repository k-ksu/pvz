package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type messageHandler interface {
	Handle(msg []byte) error
}

type ConsumerGroup struct {
	brokers       []string
	consumerGroup sarama.ConsumerGroup
	ready         chan bool
	msgHandler    messageHandler
}

func newConsumerGroup(brokers []string, group string) (sarama.ConsumerGroup, error) {
	log.Println("Starting a new Sarama consumer")

	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	client, err := sarama.NewConsumerGroup(brokers, group, config)

	if err != nil {
		return nil, errors.Wrap(err, "error creating consumer group client")
	}

	return client, nil
}

func NewConsumerGroup(brokers []string, group string, msgHandler messageHandler) (*ConsumerGroup, error) {
	client, err := newConsumerGroup(brokers, group)

	if err != nil {
		return nil, err
	}

	return &ConsumerGroup{
		brokers:       brokers,
		consumerGroup: client,
		ready:         make(chan bool),
		msgHandler:    msgHandler,
	}, nil
}

func (c *ConsumerGroup) Ready() <-chan bool {
	return c.ready
}

func (c *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (c *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			err := c.msgHandler.Handle(message.Value)
			if err != nil {
				fmt.Println("Consumer group error", err)
			}

			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func (c *ConsumerGroup) Run(topic string) {
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := c.consumerGroup.Consume(ctx, []string{topic}, c); err != nil {
				log.Fatalf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-c.Ready()
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

loop:
	for {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			break loop
		case <-sigterm:
			log.Println("terminating: via signal")
			break loop
		}
	}

	cancel()
	wg.Wait()

	if err := c.consumerGroup.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
	log.Println("terminated")
}
