package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationRes struct {
	ID         primitive.ObjectID `json:"id"`
	Title      string             `json:"title"`
	Body       string             `json:"body"`
	Status     int8               `json:"status"`
	IsRead     bool               `json:"is_read"`
	Created_At time.Time          `json:"created_at"`
}

type NotificationUpdateReq struct {
	Title  string `json:"title"`
	Body   string `json:"body"`
	Status int8   `json:"status"`
	IsRead bool   `json:"is_read"`
}
