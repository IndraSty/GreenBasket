package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StoreReq struct {
	Name        string `json:"name" valid:"required,min=2,max=100" bson:"name"`
	Description string `json:"description" valid:"required" bson:"description"`
	Logo        string `json:"logo" bson:"logo"`
	Banner      string `json:"banner" bson:"banner"`
}

type AddStoreRes struct {
	InsertId *primitive.ObjectID
}

type GetStoreRes struct {
	Name        string `json:"name" valid:"required,min=2,max=100" bson:"name"`
	Description string `json:"description" valid:"required" bson:"description"`
	Logo        string `json:"logo" bson:"logo"`
	Banner      string `json:"banner" bson:"banner"`
	Email       string `json:"email" bson:"email"`
	Store_Id    string `json:"store_id" bson:"store_id"`
}

type UpdateStoreRes struct {
	UpdateResult *mongo.UpdateResult
}

type DeleteStoreRes struct {
	DeleteResult *mongo.DeleteResult
}
