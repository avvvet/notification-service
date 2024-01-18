package main

import (
	"log"

	"github.com/avvvet/notification-service/internal/handler"
)

func main() {
	// Initialize the QuoteHandler
	quoteHandler := handler.NewQuoteHandler()

	log.Println("Notification Service starting...")

	// Start the email consumer
	quoteHandler.Start()

	log.Println("Notification Service is now running.")

	// Placeholder: Simulate sending a quote request
	quoteRequest := &handler.QuoteRequest{
		CustomerName: "John Doe",
		// Include other fields relevant to the quote request
	}

	// Handle the quote request
	quoteHandler.HandleQuoteRequest(quoteRequest)

	// Keep the application running (use a signal handler or other mechanism in a real application)
	select {}
}
