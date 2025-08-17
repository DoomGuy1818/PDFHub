package client

import (
	"PDFHub/internal/bot/lib/e"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Client struct {
	producer *kafka.Producer
	consumer *kafka.Consumer
}

func New(brokers string, groupID string) (*Client, error) {
	p, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": brokers,
		},
	)

	if err != nil {
		return nil, e.Wrap("Can't create producer", err)
	}

	c, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers": brokers,
			"group.id":          groupID,
		},
	)

	if err != nil {
		return nil, e.Wrap("Can't create consumer", err)
	}

	return &Client{
		producer: p,
		consumer: c,
	}, nil
}
func (c *Client) ReadMessage(topic string, time time.Duration) (string, error) {

	err := c.consumer.Subscribe(
		topic, func(c *kafka.Consumer, e kafka.Event) error {
			return nil
		},
	)

	if err != nil {
		return "", e.Wrap("Can't subscribe to topic", err)
	}

	msg, err := c.consumer.ReadMessage(time)

	if err != nil {
		return "", e.Wrap("Can't read message", err)
	}

	return string(msg.Value), nil
}

func (c *Client) SendMessage(topic string, value string) error {
	err := c.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value:         []byte(value),
			TimestampType: kafka.TimestampCreateTime,
			Timestamp:     time.Now(),
		},
		//Тут можно указать канал, в который придёт TopicPartition либо с ошибкой либо без таким образом получим инфу об ошибке
		nil,
	)

	if err != nil {
		return e.Wrap("Can't send message", err)
	}

	return nil
}
