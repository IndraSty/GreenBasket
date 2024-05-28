package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/mongo"
)

type Contact struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
}

type ContactRepository interface {
	AddStoreContact(ctx context.Context, email string, storeID string, contact Contact, updateAt time.Time) (*mongo.UpdateResult, error)
	GetStoreContact(ctx context.Context, email string, storeID string) (*Contact, error)
	UpdateStoreContact(ctx context.Context, email string, storeID string, update Contact, updateAt time.Time) (*mongo.UpdateResult, error)
	RemoveStoreContact(ctx context.Context, email string, storeID string, updateAt time.Time) (*mongo.UpdateResult, error)
}

type ContactService interface {
	AddStoreContact(ctx context.Context, email string, storeID string, req Contact) (*dto.ContactRes, error)
	GetStoreContact(ctx context.Context, email string, storeID string) (*Contact, error)
	UpdateStoreContact(ctx context.Context, email string, storeID string, req Contact) (*dto.ContactRes, error)
	RemoveStoreContact(ctx context.Context, email string, storeID string) (*dto.ContactDelRes, error)
}
