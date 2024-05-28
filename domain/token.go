package domain

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

type TokenService interface {
	GenerateAllTokens(email string, firstname string, lastname string, uid string) (signedToken string, signedRefreshToken string, err error)
	UpdateRefreshToken(signedRefreshToken string, userId string, usercol *mongo.Collection)
	ValidateToken(signedToken string) (claims *SignedDetails, msg string)
}
