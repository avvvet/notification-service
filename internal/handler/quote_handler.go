// internal/handler/quote_handler.go
package handler

import (
	"encoding/json"
	"log"

	"github.com/avvvet/notification-service/internal/email"
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

// consume email and send it
func (q *QuoteHandler) startEmailConsumer() {
	emailQueue := q.emailQueue

	// Consume email notifications from RabbitMQ
	_, err := q.rabbitMQ.Channel.QueueDeclare(
		emailQueue, // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatal("Error declaring email queue:", err)
	}

	msgs, err := q.rabbitMQ.Channel.Consume(
		emailQueue, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatal("Error consuming email queue:", err)
	}

	// Handle incoming email notifications
	for msg := range msgs {
		// Convert the message body to a notification
		notification := bytesToNotification(msg.Body)

		// Send email based on the notification
		emailContent := getEmailContent(notification)
		smtpConfig := email.GetSMTPConfig()
		email.SendEmail(emailContent, smtpConfig)
	}
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

	rabbitMQ, err := rabbitmq.NewRabbitMQ("amqp://rabbit:password@localhost:5672")
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

func notificationToBytes(notification *notification.Notification) []byte {
	// Implement logic to serialize the notification to a byte slice (e.g., JSON encoding)
	return []byte("notification_bytes_placeholder")
}

func bytesToNotification(body []byte) *notification.Notification {
	var notif notification.Notification
	err := json.Unmarshal(body, &notif)
	if err != nil {
		log.Println("Error decoding notification from bytes:", err)
		// Handle error (e.g., log, skip, etc.)
	}
	return &notif
}

// getEmailContent creates an Email struct based on the notification.
func getEmailContent(notif *notification.Notification) *email.Email {
	subject := "New Quote Request"
	body := "Dear " + notif.Recipient + ",\n\n" +
		"Thank you for your quote request. We will review it and get back to you shortly.\n\n" +
		"Best regards,\nThe Notification Service"

	return &email.Email{
		Subject: subject,
		Body:    body,
		To:      notif.Recipient,
	}
}
