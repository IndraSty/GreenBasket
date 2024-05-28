package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Get() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error when load envi %s", err.Error())
	}

	return &Config{
		Token{
			Secret_Key: os.Getenv("JWT_SECRET_KEY"),
		},
		Email{
			Host:     os.Getenv("MAIL_HOST"),
			Name:     os.Getenv("MAIL_NAME"),
			Password: os.Getenv("APP_PASSWORD"),
		},
		Redis{
			Addr: os.Getenv("REDIS_ADDR"),
			Pass: os.Getenv("REDIS_PASS"),
		},
		Midtrans{
			Key:    os.Getenv("MIDTRANS_KEY"),
			IsProd: os.Getenv("MIDTRANS_ENV") == "production",
		},
		MongoDB{
			URI: os.Getenv("MONGO_URI"),
		},
		Server{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Auth{
			Secret_Key:          os.Getenv("AUTH_SECRET_KEY"),
			MaxAge:              os.Getenv("AUTH_MAX_AGE"),
			IsProd:              os.Getenv("AUTH_IS_PROD"),
			GoogleCallBackUrl:   os.Getenv("GOOGLE_AUTH_CALLBACK_URL"),
			FacebookCallBackUrl: os.Getenv("FACEBOOK_AUTH_CALLBACK_URL"),
		},
		Google{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			ScopeEmail:   os.Getenv("GOOGLE_SCOPE_EMAIL"),
			ScopeProfile: os.Getenv("GOOGLE_SCOPE_PROFILE"),
			State:        os.Getenv("GOOGLE_STATE"),
			TokenUrl:     os.Getenv("GOOGLE_TOKEN_URL"),
		},
		Facebook{
			ClientID:     os.Getenv("FACEBOOK_CLIENT_ID"),
			ClientSecret: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		},
	}
}
