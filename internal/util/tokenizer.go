package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/config"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tokenService struct {
	cnf *config.Config
}

func NewTokenService(cnf *config.Config) domain.TokenService {
	return &tokenService{cnf: cnf}
}

func (ts *tokenService) GenerateAllTokens(email string, firstname string, lastname string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &domain.SignedDetails{
		Email:      email,
		First_Name: firstname,
		Last_Name:  lastname,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Minute * time.Duration(5)).Unix(),
		},
	}

	refreshClaims := &domain.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(ts.cnf.Token.Secret_Key))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshClaims).SignedString([]byte(ts.cnf.Token.Secret_Key))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func (ts *tokenService) UpdateRefreshToken(signedRefreshToken string, userId string, usercol *mongo.Collection) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	update_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "update_at", Value: update_at})

	upsert := true
	filter := bson.M{"user_id": userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := usercol.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateObj},
	}, &opt)

	if err != nil {
		log.Panic(err)
		return
	}

}

func (ts *tokenService) ValidateToken(signedToken string) (claims *domain.SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&domain.SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(ts.cnf.Token.Secret_Key), nil
		},
	)

	// check token is invalid
	claims, ok := token.Claims.(*domain.SignedDetails)
	if !ok {
		msg = fmt.Sprintf("token is invalid: %v", err)
		return
	}

	// check token is expired
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = fmt.Sprintf("token is expired: %v", err)
		return
	}

	return claims, msg
}
