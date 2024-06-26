package dto

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SellerRegisterReq struct {
	First_Name string `json:"first_name" valid:"required,minstringlength(2),maxstringlength(100)"`
	Last_Name  string `json:"last_name" valid:"required,minstringlength(2),maxstringlength(100)"`
	Email      string `json:"email" valid:"email,required"`
	Password   string `json:"password" valid:"required,minstringlength(8)"`
	Phone      string `json:"phone" valid:"required,minstringlength(11)"`
}

type SellerRegisterRes struct {
	InsertId primitive.ObjectID `json:"insert_id"`
	Message  []string           `json:"message"`
}

type SellerAuthReq struct {
	Email    string `json:"email" valid:"email,required"`
	Password string `json:"password" valid:"required,minstringlength(8)"`
}

type SellerAuthRes struct {
	Access_Token  string `json:"access_token"`
	Refresh_Token string `json:"refresh_token"`
}

type SellerUserResponse struct {
	First_Name string `json:"first_name"`
	Last_Name  string `json:"last_name"`
	Email      string `json:"email"`
	Image_Url  string `json:"image_url"`
	Phone      string `json:"phone"`
	Role       string `json:"role"`
	User_Id    string `json:"user_id"`
}

type SellerUpdateReq struct {
	First_Name string `json:"first_name" valid:"required,minstringlength(2),maxstringlength(100)"`
	Last_Name  string `json:"last_name" valid:"required,minstringlength(2),maxstringlength(100)"`
	Image_Url  string `json:"image_url"`
}

type SellerUpdateRes struct {
	UpdateResult *mongo.UpdateResult
}
