package rabbit

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

const (
	_logsQueueName      = "logs"
	_messageContentType = "text/plain"
)

type RabbitMQConfig struct {
	URL string `default:"amqp://guest:guest@amqp/"`
}

func ConnectToRabbitMQ(rabbitURL string) (
	*amqp.Connection,
	*amqp.Channel,
	amqp.Queue,
	error,
) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	q, err := ch.QueueDeclare(
		_logsQueueName,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, nil, amqp.Queue{}, err
	}

	return conn, ch, q, nil
}

type rabbitMQWriter struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
	Context context.Context
}

func (w *rabbitMQWriter) Write(p []byte) (
	n int,
	err error,
) {
	err = w.Channel.PublishWithContext(
		w.Context,
		"",
		w.Queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: _messageContentType,
			Body:        p,
		},
	)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func NewLogger(
	ctx context.Context,
	channel *amqp.Channel,
	queue amqp.Queue,
) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetOutput(&rabbitMQWriter{Channel: channel, Queue: queue, Context: ctx})

	return logger
}

func NewConsumer(channel *amqp.Channel, queue amqp.Queue) (func(), error) {
	messages, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	consumerWork := func() {
		log.Println("Starting to consume messages")
		for message := range messages {
			log.Print(string(message.Body))
		}
	}

	return consumerWork, nil
}
