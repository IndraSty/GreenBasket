package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Seller struct {
	ID              primitive.ObjectID `bson:"_id"`
	First_Name      string             `json:"first_name" valid:"required,min=2,max=100"`
	Last_Name       string             `json:"last_name" valid:"required,min=2,max=100"`
	Email           string             `json:"email" valid:"email,required"`
	Password        string             `json:"password" valid:"required,min=8"`
	Image_Url       string             `json:"image_url"`
	Phone           string             `json:"phone" valid:"required"`
	Refresh_Token   string             `json:"refresh_token"`
	Role            string             `json:"role"`
	Created_At      time.Time          `json:"created_at"`
	Updated_At      time.Time          `json:"updated_at"`
	EmailVerified   bool               `json:"email_verified"`
	PhoneVerified   bool               `json:"phone_verified"`
	Seller_Id       string             `json:"seller_id"`
	Store_Id        string             `json:"store_id" bson:"store_id"`
	Address_Details *Address           `json:"address" bson:"address"`
}

type SellerRepository interface {
	CreateSeller(ctx context.Context, seller Seller) (primitive.ObjectID, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)
	FindSellerByEmail(ctx context.Context, email string) (*Seller, error)
	FindSellerByStoreId(ctx context.Context, storeID string) (*Seller, error)
	UpdateSeller(ctx context.Context, email string, update bson.D) (*mongo.UpdateResult, error)
	AddStoreId(ctx context.Context, email string, storeID string) error
}

type SellerService interface {
	RegisterSeller(ctx context.Context, req *dto.SellerRegisterReq) (*dto.SellerRegisterRes, error)
	ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error
	AuthenticateSeller(ctx context.Context, req *dto.SellerAuthReq) (*dto.SellerAuthRes, error)
	GetSellerByEmail(ctx context.Context, email string) (*Seller, error)
	UpdateSeller(ctx context.Context, email string, req *dto.SellerUpdateReq) (*dto.SellerUpdateRes, error)
}
