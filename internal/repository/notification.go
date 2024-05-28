package repository

import (
	"context"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type notificationRepository struct {
	Collection *mongo.Collection
}

func NewNotificationRepository(client *mongo.Client) domain.NotificationRepository {
	return &notificationRepository{
		Collection: db.OpenCollection(client, "Notifications"),
	}
}

// FindByUser implements domain.NotificationRepository.
func (repo *notificationRepository) FindByUser(ctx context.Context, userID string) ([]domain.Notification, error) {
	var notifications []domain.Notification
	filter := bson.M{"user_id": userID}
	cur, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var notification domain.Notification
		err := cur.Decode(&notification)
		if err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

// Insert implements domain.NotificationRepository.
func (repo *notificationRepository) Insert(ctx context.Context, notification *domain.Notification) error {
	_, err := repo.Collection.InsertOne(ctx, notification)
	if err != nil {
		return err
	}

	return nil
}

// Update implements domain.NotificationRepository.
func (repo *notificationRepository) Update(ctx context.Context, id primitive.ObjectID, userID string, notification *dto.NotificationUpdateReq) error {
	filter := bson.M{"_id": id, "user_id": userID}
	_, err := repo.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: notification}})
	if err != nil {
		return err
	}

	return nil
}
