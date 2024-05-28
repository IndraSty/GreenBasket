package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	ID              primitive.ObjectID `bson:"_id"`
	Name            string             `json:"name" valid:"required,min=2,max=100" bson:"name"`
	Description     string             `json:"description" valid:"required" bson:"description"`
	Logo            string             `json:"logo" bson:"logo"`
	Banner          string             `json:"banner" bson:"banner"`
	Created_At      time.Time          `json:"created_at" bson:"created_at"`
	Updated_At      time.Time          `json:"updated_at" bson:"updated_at"`
	Email           string             `json:"email" bson:"email"`
	Store_Id        string             `json:"store_id" bson:"store_id"`
	Contact_Details *Contact           `json:"contact" bson:"contact"`
	Address_Details *Address           `json:"address" bson:"address"`
}

type StoreRepository interface {
	CreateStore(ctx context.Context, store Store) (primitive.ObjectID, error)
	GetStore(ctx context.Context, storeID string, email ...string) (*Store, error)
	GetStoreByEmail(ctx context.Context, email string) (*Store, error)
	CheckNameExists(ctx context.Context, name string) (bool, error)
	UpdateStore(ctx context.Context, email, storeID string, update bson.D) (*mongo.UpdateResult, error)
	RemoveStore(ctx context.Context, email, storeID string) (*mongo.DeleteResult, error)
	GetStoreByQuery(ctx context.Context, query string) ([]Store, error)
}

type StoreService interface {
	CreateStore(ctx context.Context, email string, req *dto.StoreReq) (*dto.AddStoreRes, error)
	GetStoreByIdAndEmail(ctx context.Context, email, storeID string) (*dto.GetStoreRes, error)
	UpdateStore(ctx context.Context, email, storeID string, req *dto.StoreReq) (*dto.UpdateStoreRes, error)
	DeleteStore(ctx context.Context, email, storeID string) (*dto.DeleteStoreRes, error)
	SearchStore(ctx context.Context, query string) ([]Store, error)
}
