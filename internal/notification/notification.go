package notification

import "time"

type Notification struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Recipient string    `json:"recipient"`
	Timestamp time.Time `json:"tiimestamp"`
}

func NewNotification(message, recipient string) *Notification {
	return &Notification{
		ID:        generateID(),
		Message:   message,
		Recipient: recipient,
		Timestamp: time.Now(),
	}
}

func generateID() string {
	return "unique_id"
}
