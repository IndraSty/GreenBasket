package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Cart struct {
	ID         primitive.ObjectID `bson:"_id"`
	Email      string             `json:"email" bson:"email"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	TotalPrice float64            `json:"total_price" bson:"total_price"`
	Items      []CartItem         `json:"items" bson:"items"`
}

type CartItem struct {
	Product_Id    string    `json:"product_id" bson:"product_id"`
	Product_Name  string    `json:"product_name" bson:"product_name"`
	Product_Image []string  `json:"product_image" bson:"product_image"`
	StoreID       string    `json:"store_id" bson:"store_id"`
	Quantity      int       `json:"quantity" bson:"quantity"`
	AddedAt       time.Time `json:"added_at" bson:"added_at"`
	Selected      bool      `json:"selected" bson:"selected"`
	Price         float64   `json:"price" bson:"price"`
}

type CartRepository interface {
	CreateCart(ctx context.Context, cart *Cart) error
	GetUserCart(ctx context.Context, email string) (*Cart, error)
	CheckUserCart(ctx context.Context, email string) (bool, error)
	AddToCart(ctx context.Context, email string, item *CartItem) (*mongo.UpdateResult, error)
	GetAllCartItem(ctx context.Context, email string) (*[]CartItem, error)
	UpdateCartItemById(ctx context.Context, email, productID string, value *dto.CartItemEditRepo, updateAt time.Time) (*mongo.UpdateResult, error)
	UpdateTotalPrice(ctx context.Context, email string, value float64) error
	RemoveCartItemById(ctx context.Context, email, productID string, updateAt time.Time) (*mongo.UpdateResult, error)
}

type CartService interface {
	CreateCart(ctx context.Context, email string) error
	GetUserCart(ctx context.Context, email string) (*Cart, error)
	AddToCart(ctx context.Context, email, productID string, req *dto.AddCartReq) error
	GetAllCartItem(ctx context.Context, email string) (*[]dto.GetCartItemRes, error)
	UpdateCartItemById(ctx context.Context, email, productID string, input *dto.CartItemEditReq) error
	RemoveCartItemById(ctx context.Context, email, productID string) error
}
