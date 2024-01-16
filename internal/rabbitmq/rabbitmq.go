package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type RabbitMQ struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewRabbitMQ(amqpURL string) (*RabbitMQ, error) {
	conn, err := amqp091.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
	}, nil
}
