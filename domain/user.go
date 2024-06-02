package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id"`
	First_Name      string             `json:"first_name" valid:"required,min=2,max=100" bson:"first_name"`
	Last_Name       string             `json:"last_name" valid:"required,min=2,max=100" bson:"last_name"`
	Email           string             `json:"email" valid:"email,required" bson:"email"`
	Password        string             `json:"password" valid:"required,minstringlength(8)" bson:"password"`
	Image_Url       string             `json:"image_url" bson:"image_url"`
	Phone           string             `json:"phone" bson:"phone"`
	Refresh_Token   string             `json:"refresh_token" bson:"refresh_token"`
	Role            string             `json:"role" bson:"role"`
	Created_At      time.Time          `json:"created_at" bson:"created_at"`
	Updated_At      time.Time          `json:"updated_at" bson:"updated_at"`
	User_Id         string             `json:"user_id" bson:"user_id"`
	EmailVerified   bool               `json:"email_verified" bson:"email_verified"`
	Oauth_Id        string             `json:"oauth_id"`
	Address_Details *Address           `json:"address" bson:"address"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user User) (primitive.ObjectID, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)
	FindUserByEmail(ctx context.Context, email string) (*User, error)
	FindUserById(ctx context.Context, userId string) (*User, error)
	UpdateUser(ctx context.Context, email string, update bson.D) (*mongo.UpdateResult, error)
}

type UserService interface {
	RegisterUser(ctx context.Context, req *dto.UserRegisterReq) (*dto.UserRegisterRes, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, email string, req *dto.UserUpdateReq) (*dto.UserUpdateRes, error)
	AddPhoneNumber(ctx context.Context, email string, req *dto.AddPhone) error
}
