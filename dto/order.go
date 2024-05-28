package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type InsertOrderRes struct {
	InsertId primitive.ObjectID
}

type OrderSellerUpdateReq struct {
	Status         string `json:"status" bson:"status"`
	Payment_Status string `json:"payment_status" bson:"payment_status"`
}

type OrderStatusUpdateReq struct {
	Status string `json:"status" bson:"status"`
}
