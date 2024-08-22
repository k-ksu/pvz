package events

import (
	"encoding/json"
	"fmt"

	"HomeWork_1/internal/model"
)

type Reader struct {
	KafkaEnable bool
}

func NewReader(KafkaEnable bool) *Reader {
	return &Reader{
		KafkaEnable: KafkaEnable,
	}
}

func (e *Reader) Handle(msg []byte) error {

	var event model.EventMessage
	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	if e.KafkaEnable {
		fmt.Printf("New event: command %v was called with: \nArgs %v\ntime %v\n\n", event.Method, event.Args, event.TimeStamp)
	}

	return nil
}
