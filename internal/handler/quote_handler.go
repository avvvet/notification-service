// internal/handler/quote_handler.go
package handler

import (
	"encoding/json"
	"log"
	"os"

	"github.com/avvvet/notification-service/internal/email"
	"github.com/avvvet/notification-service/internal/mongodb"
	"github.com/avvvet/notification-service/internal/notification"
	"github.com/avvvet/notification-service/internal/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
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
		true,       // durable
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
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Fatal("Error consuming email queue:", err)
	}

	// Handle incoming email notifications
	go func() {
		for msg := range msgs {
			q.handleConsumeMsg(&msg)
		}
	}()

}

// HandleQuoteRequest handles the quote request and triggers the notification workflow.
func (q *QuoteHandler) HandleQuoteRequest(quoteRequest *QuoteRequest) {
	// Generate a notification for the quote request
	notification := notification.NewNotification("Quote Request Received", quoteRequest.CustomerName)

	//convert struct to byte
	bytes, err := notificationToBytes(notification)
	if err != nil {
		log.Printf("Error encoding notification to JSON: %v", err)
	}

	// Publish the notification to the email queue
	err = q.rabbitMQ.Publish(q.emailQueue, bytes)
	if err != nil {
		log.Println("Error publishing notification to RabbitMQ:", err)
		// Handle error (e.g., retry, log, etc.)
	}

	// Store the notification in MongoDB for persistence
	// err = q.mongoDB.StoreNotification(notification)
	// if err != nil {
	// 	log.Println("Error storing notification in MongoDB:", err)
	// 	// Handle error (e.g., retry, log, etc.)
	// }
}

func NewQuoteHandler() *QuoteHandler {

	rabbitMQ, err := rabbitmq.NewRabbitMQ("amqp://rabbit:password@localhost:5672")
	if err != nil {
		log.Fatal("Error initializing RabbitMQ:", err)
	}

	// Declare a queue (optional, depends on your scenario)
	queueName := "quote_email_queue"

	_, err = rabbitMQ.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	log.Println("Queue declared successfully.")

	mongoDB, err := mongodb.NewMongoDB("mongodb://localhost:27017/", "notificationDB")
	if err != nil {
		log.Fatal("Error initializing MongoDB:", err)
	}

	return &QuoteHandler{
		rabbitMQ:   rabbitMQ,
		mongoDB:    mongoDB,
		emailQueue: queueName,
	}
}

// Start starts the quote handler.
func (q *QuoteHandler) Start() {
	go q.startEmailConsumer()
}

func (q *QuoteHandler) handleConsumeMsg(msg *amqp.Delivery) {
	// Convert the message body to a notification
	//notification := bytesToNotification(msg.Body)

	// Logging to a file
	q.logToFile(msg.Body)

	// Acknowledge the message
	msg.Ack(false)

	// Send email based on the notification
	// emailContent := getEmailContent(notification)
	// smtpConfig := email.GetSMTPConfig()
	// email.SendEmail(emailContent, smtpConfig)
}

// logToFile appends processed messages to a log file.
func (q *QuoteHandler) logToFile(message []byte) {
	// Open or create a log file for appending
	file, err := os.OpenFile("message_consume_log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return
	}
	defer file.Close()

	// Log the message to the file
	logger := log.New(file, "", log.LstdFlags)
	logger.Printf("Message processed: %s", message)
}

func notificationToBytes(notification *notification.Notification) ([]byte, error) {
	bytes, err := json.Marshal(notification)
	if err != nil {
		return nil, err
	}

	return bytes, nil
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
