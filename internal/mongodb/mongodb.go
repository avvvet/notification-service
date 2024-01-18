package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/avvvet/notification-service/internal/notification"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDB(mongoURI, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	con, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	return &MongoDB{
		client:   con,
		database: con.Database(dbName),
	}, nil
}

func (m *MongoDB) StoreNotification(notification *notification.Notification) error {
	_, err := m.database.Collection("notifications").InsertOne(context.TODO(), notification)
	if err != nil {
		log.Println("Error storing notification in MongoDB:", err)
		return err
	}
	return nil
}
