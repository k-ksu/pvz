package events

import (
	"encoding/json"
	"fmt"

	"HomeWork_1/internal/model"
)

type producer interface {
	SendSyncMessage(msg []byte) (partition int32, offset int64, err error)
}

type Sender struct {
	producer    producer
	KafkaEnable bool
}

func NewSender(KafkaEnable bool, producer producer) *Sender {
	return &Sender{
		producer:    producer,
		KafkaEnable: KafkaEnable,
	}
}

func (s *Sender) SendMessage(event *model.EventMessage) error {
	msg, err := json.Marshal(event)
	if err != nil {
		fmt.Println("Send message marshal error", err)
		return err
	}

	_, _, err = s.producer.SendSyncMessage(msg)
	if err != nil {
		return err
	}

	if !s.KafkaEnable {
		fmt.Printf("New event was sent to kafka: command %v was called with: \nArgs %v\ntime %v\n\n", event.Method, event.Args, event.TimeStamp)
	}

	return nil
}
