package dto

import "go.mongodb.org/mongo-driver/mongo"

type ContactReq struct {
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
}

type ContactRes struct {
	UpdateResult *mongo.UpdateResult
}

type ContactDelRes struct {
	UpdateResult *mongo.UpdateResult
}
