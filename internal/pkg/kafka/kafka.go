package kafka

import (
	"github.com/IBM/sarama"
	"github.com/pkg/errors"
)

type Producer struct {
	brokers      []string
	topic        string
	syncProducer sarama.SyncProducer
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()

	syncProducerConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner
	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll
	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1

	syncProducerConfig.Producer.CompressionLevel = sarama.CompressionLevelDefault

	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true

	syncProducerConfig.Producer.Compression = sarama.CompressionGZIP

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)

	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	return syncProducer, nil
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return nil, errors.Wrap(err, "error with sync kafka-producer")
	}

	return &Producer{
		brokers:      brokers,
		topic:        topic,
		syncProducer: syncProducer,
	}, nil
}

func (p *Producer) SendSyncMessage(msg []byte) (partition int32, offset int64, err error) {
	return p.syncProducer.SendMessage(p.buildMessage(msg))
}

func (p *Producer) Close() error {
	err := p.syncProducer.Close()
	if err != nil {
		return errors.Wrap(err, "kafka.Connector.Close")
	}

	return nil
}

func (p *Producer) buildMessage(msg []byte) *sarama.ProducerMessage {
	return &sarama.ProducerMessage{
		Topic:     p.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
	}
}
