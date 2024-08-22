package model

import "time"

type EventMessage struct {
	Method    string
	Args      []string
	TimeStamp time.Time
}
