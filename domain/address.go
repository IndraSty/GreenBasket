package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/mongo"
)

type Address struct {
	House   string `json:"house_name" valid:"required" bson:"house_name"`
	Street  string `json:"street_name" valid:"required" bson:"street_name"`
	City    string `json:"city_name" valid:"required" bson:"city_name"`
	Pincode string `json:"pin_code" valid:"required" bson:"pin_code"`
}

type AddressRepository interface {
	AddUserAddress(ctx context.Context, email string, address Address, updateAt time.Time) (*mongo.UpdateResult, error)
	GetUserAddress(ctx context.Context, email string) (*Address, error)
	UpdateAddressUser(ctx context.Context, email string, address Address, updateAt time.Time) (*mongo.UpdateResult, error)
	RemoveAddressUser(ctx context.Context, email string, updateAt time.Time) (*mongo.UpdateResult, error)

	AddSellerAddress(ctx context.Context, email string, address Address, updateAt time.Time) (*mongo.UpdateResult, error)
	GetSellerAddress(ctx context.Context, email string) (*Address, error)
	UpdateSellerAddress(ctx context.Context, email string, address Address, updateAt time.Time) (*mongo.UpdateResult, error)
	RemoveSellerAddress(ctx context.Context, email string, updateAt time.Time) (*mongo.UpdateResult, error)

	AddStoreAddress(ctx context.Context, email string, storeID string, address Address, updateAt time.Time) (*mongo.UpdateResult, error)
	GetStoreAddress(ctx context.Context, email string, storeID string) (*Address, error)
	UpdateStoreAddress(ctx context.Context, email string, storeID string, address Address, updateAt time.Time) (*mongo.UpdateResult, error)
	RemoveStoreAddress(ctx context.Context, email string, storeID string, updateAt time.Time) (*mongo.UpdateResult, error)
}

type AddressService interface {
	AddUserAddress(ctx context.Context, email string, req Address) (*dto.AddressRes, error)
	GetUserAddress(ctx context.Context, email string) (*Address, error)
	UpdateUserAddress(ctx context.Context, email string, req Address) (*dto.AddressRes, error)
	RemoveUserAddress(ctx context.Context, email string) (*dto.AddressDelRes, error)

	AddSellerAddress(ctx context.Context, email string, req Address) (*dto.AddressRes, error)
	GetSellerAddress(ctx context.Context, email string) (*Address, error)
	UpdateSellerAddress(ctx context.Context, email string, req Address) (*dto.AddressRes, error)
	RemoveSellerAddress(ctx context.Context, email string) (*dto.AddressDelRes, error)

	AddStoreAddress(ctx context.Context, email, storeID string, req Address) (*dto.AddressRes, error)
	GetStoreAddress(ctx context.Context, email, storeID string) (*Address, error)
	UpdateStoreAddress(ctx context.Context, email, storeID string, req Address) (*dto.AddressRes, error)
	RemoveStoreAddress(ctx context.Context, email, storeID string) (*dto.AddressDelRes, error)
}
