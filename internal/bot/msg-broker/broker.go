package msg_broker

import (
	"io"
	"time"
)

type Consumer interface {
	ReadMessage(topic string, t time.Duration) (string, error)
}

type Publisher interface {
	SendMessage(topic string, value []byte) error
	SendStream(topic string, value io.ReadCloser, key string) error
}
