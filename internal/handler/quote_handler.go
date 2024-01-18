// internal/handler/quote_handler.go
package handler

import (
	"log"

	"github.com/avvvet/notification-service/internal/mongodb"
	"github.com/avvvet/notification-service/internal/notification"
	"github.com/avvvet/notification-service/internal/rabbitmq"
)

type QuoteRequest struct {
	CustomerName string `json:"customer_name"`
	// Include other fields relevant to the quote request
}

type QuoteHandler struct {
	rabbitMQ   *rabbitmq.RabbitMQ
	mongoDB    *mongodb.MongoDB
	emailQueue string
}

// startEmailConsumer starts the email notification consumer.
func (q *QuoteHandler) startEmailConsumer() {
	// Implement logic to consume email notifications from RabbitMQ
	// and send emails using the email package
	// This could include using the email package to send emails
	// based on the notifications received from the RabbitMQ queue.
}

// HandleQuoteRequest handles the quote request and triggers the notification workflow.
func (q *QuoteHandler) HandleQuoteRequest(quoteRequest *QuoteRequest) {
	// Generate a notification for the quote request
	notification := notification.NewNotification("Quote Request Received", quoteRequest.CustomerName)

	// Publish the notification to the email queue
	err := q.rabbitMQ.Publish(q.emailQueue, notificationToBytes(notification))
	if err != nil {
		log.Println("Error publishing notification to RabbitMQ:", err)
		// Handle error (e.g., retry, log, etc.)
	}

	// Store the notification in MongoDB for persistence
	err = q.mongoDB.StoreNotification(notification)
	if err != nil {
		log.Println("Error storing notification in MongoDB:", err)
		// Handle error (e.g., retry, log, etc.)
	}
}

func NewQuoteHandler() *QuoteHandler {
	rabbitMQ, err := rabbitmq.NewRabbitMQ("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Error initializing RabbitMQ:", err)
	}

	mongoDB, err := mongodb.NewMongoDB("mongodb://localhost:27017/", "notificationDB")
	if err != nil {
		log.Fatal("Error initializing MongoDB:", err)
	}

	return &QuoteHandler{
		rabbitMQ:   rabbitMQ,
		mongoDB:    mongoDB,
		emailQueue: "quote_email_queue",
	}
}

// Start starts the quote handler.
func (q *QuoteHandler) Start() {
	go q.startEmailConsumer()
}

// notificationToBytes converts a notification to a byte slice.
func notificationToBytes(notification *notification.Notification) []byte {
	// Implement logic to serialize the notification to a byte slice (e.g., JSON encoding)
	return []byte("notification_bytes_placeholder")
}
