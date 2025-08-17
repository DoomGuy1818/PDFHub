package msg_broker

import (
	"time"
)

type Consumer interface {
	ReadMessage(topic string, t time.Duration) (string, error)
}

type Publisher interface {
	SendMessage(topic string, value string) error
}
