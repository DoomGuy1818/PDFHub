package kafka

import (
	"PDFHub/internal/bot/lib/e"
	"fmt"
	"io"
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
			"bootstrap.servers":  brokers,
			"group.id":           groupID,
			"compression.type":   "zstd",
			"linger.ms":          10,
			"batch.num.messages": 1000,
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

func (c *Client) SendMessage(topic string, value []byte) error {
	err := c.producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value:         value,
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

func (c *Client) SendStream(topic string, value io.ReadCloser, key string) error {
	pr, pw := io.Pipe()

	go func() {
		defer pr.Close()
		buf := make([]byte, 1024*1024)
		for {
			n, err := pr.Read(buf)
			if n > 0 {
				err := c.producer.Produce(
					&kafka.Message{
						TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
						Value:          buf[:n],
						Key:            []byte(key),
					}, nil,
				)
				if err != nil {
					fmt.Println("Ошибка продюса в Kafka:", err)
					break
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Ошибка чтения из Pipe:", err)
				break
			}
		}
	}()

	_, err := io.Copy(pw, value)
	pw.Close()
	if err != nil {
		return e.Wrap("Can't send message", err)
	}

	c.producer.Flush(5000)
	fmt.Println("Файл успешно отправлен в Kafka через поток")

	return nil
}
