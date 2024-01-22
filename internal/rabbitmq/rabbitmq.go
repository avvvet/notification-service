package rabbitmq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(amqpURL string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	log.Println("RabbitMQ connection established successfully.")

	return &RabbitMQ{
		conn:    conn,
		Channel: channel,
	}, nil
}

func (r *RabbitMQ) Publish(queueName string, message []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := r.Channel.PublishWithContext(
		ctx,
		"",        //exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)

	return err
}

func (r *RabbitMQ) Close() {
	r.Channel.Close()
	r.conn.Close()
}
