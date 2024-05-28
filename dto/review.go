package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddReviewReq struct {
	Rating float32 `json:"rating" bson:"rating"`
	Review string  `json:"review" bson:"review"`
}

type ResponSellerReq struct {
	Seller_Response string `json:"seller_response" bson:"seller_response"`
}

type AddReviewRes struct {
	InsertId primitive.ObjectID
	Messages string
}

type GetReviewRes struct {
	Review_Id   string    `json:"review_id" bson:"review_id"`
	Product_Id  string    `json:"product_id" bson:"product_id"`
	Email       string    `json:"email" bson:"email"`
	Rating      float32   `json:"rating" bson:"rating"`
	Review      string    `json:"review" bson:"review"`
	Reviewed_At time.Time `json:"reviewed_at" bson:"reviewed_at"`
	Updated_At  time.Time `json:"updated_at" bson:"updated_at"`
}
