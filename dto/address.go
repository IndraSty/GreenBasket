package dto

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type AddressReq struct {
	House   string `json:"house_name" valid:"required" bson:"house_name"`
	Street  string `json:"street_name" valid:"required" bson:"street_name"`
	City    string `json:"city_name" valid:"required" bson:"city_name"`
	Pincode string `json:"pin_code" valid:"required" bson:"pin_code"`
}

type AddressRes struct {
	UpdateResult *mongo.UpdateResult
}

type AddressDelRes struct {
	UpdateResult *mongo.UpdateResult
}

type GetAddresRes struct {
	House   string `json:"house_name" valid:"required" bson:"house_name"`
	Street  string `json:"street_name" valid:"required" bson:"street_name"`
	City    string `json:"city_name" valid:"required" bson:"city_name"`
	Pincode string `json:"pin_code" valid:"required" bson:"pin_code"`
}
