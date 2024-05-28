package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Review struct {
	Id              primitive.ObjectID `bson:"_id"`
	Review_Id       string             `json:"review_id" bson:"review_id"`
	Product_Id      string             `json:"product_id" bson:"product_id"`
	Seller_Email    string             `json:"seller_email" bson:"seller_email"`
	Email           string             `json:"email" bson:"email"`
	Rating          float32            `json:"rating" bson:"rating"`
	Review          string             `json:"review" bson:"review"`
	Seller_Response string             `json:"seller_response" bson:"seller_response"`
	Reviewed_At     time.Time          `json:"reviewed_at" bson:"reviewed_at"`
	Updated_At      time.Time          `json:"updated_at" bson:"updated_at"`
}

type ReviewRepository interface {
	InsertReview(ctx context.Context, input Review) (primitive.ObjectID, error)
	UpdateReview(ctx context.Context, reviewID string, update bson.D) (*mongo.UpdateResult, error)
	DeleteReview(ctx context.Context, reviewID string) (*mongo.DeleteResult, error)
	GetUserReviewByEmailAndId(ctx context.Context, email, reviewID string) (*Review, error)
	GetReviewById(ctx context.Context, reviewID string) (*Review, error)
	GetReviewByProductId(ctx context.Context, productID string) (*Review, error)
	GetAllReviewByProductId(ctx context.Context, productID, sellerEmail string) (*[]Review, error)
	GetAllReviewByUserEmail(ctx context.Context, email string) (*[]Review, error)
	GetAllReviewBySellerEmail(ctx context.Context, sellerEmail string) (*[]Review, error)
}

type ReviewService interface {
	CreateReview(ctx context.Context, email, orderID, productID string, req *dto.AddReviewReq) (*dto.AddReviewRes, error)
	UpdateReview(ctx context.Context, email, reviewID string, req *dto.AddReviewReq) error
	DeleteReview(ctx context.Context, email, reviewID string) error
	GetUserReviewById(ctx context.Context, email, reviewID string) (*dto.GetReviewRes, error)
	GetAllReviewByUserEmail(ctx context.Context, email string) (*[]dto.GetReviewRes, error)
	GetAllReviewBySellerEmail(ctx context.Context, email string) (*[]dto.GetReviewRes, error)
	GetAllReviewByProductId(ctx context.Context, productID, sellerEmail string) (*[]dto.GetReviewRes, error)
	UpdateResponSeller(ctx context.Context, email, reviewID string, req *dto.ResponSellerReq) error
}
