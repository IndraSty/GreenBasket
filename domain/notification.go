package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Notification struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	Email      string             `json:"email"`
	Title      string             `json:"title"`
	Body       string             `json:"body"`
	Status     int8               `json:"status"`
	IsRead     bool               `json:"is_read"`
	Created_At time.Time          `json:"created_at"`
}

type NotificationRepository interface {
	FindByUser(ctx context.Context, userID string) ([]Notification, error)
	Insert(ctx context.Context, notification *Notification) error
	Update(ctx context.Context, id primitive.ObjectID, userID string, notification *dto.NotificationUpdateReq) error
}

type NotificationService interface {
	FindByUser(ctx context.Context, userId string) ([]dto.NotificationRes, error)
	Insert(ctx context.Context, userId string, code string, data map[string]string) error
}
