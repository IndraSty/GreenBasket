package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SellerOrder struct {
	ID             primitive.ObjectID `bson:"_id"`
	Order_id       string             `json:"order_id" bson:"order_id"`
	Email          string             `json:"email" bson:"email"`
	Ordered_At     time.Time          `json:"ordered_at" bson:"ordered_at"`
	Updated_At     time.Time          `json:"updated_at" bson:"updated_at"`
	Total_Price    float64            `json:"total_price" bson:"total_price"`
	Payment_Status string             `json:"payment_status" bson:"payment_status"`
	Items          []SellerOrderItem  `json:"items" bson:"items"`
}

type SellerOrderItem struct {
	User_Email       string   `json:"user_email" bson:"user_email"`
	Product_Id       string   `json:"product_id" bson:"product_id"`
	Product_Name     string   `json:"product_name" bson:"product_name"`
	Product_Image    []string `json:"product_image" bson:"product_image"`
	Quantity         int      `json:"quantity" bson:"quantity"`
	Price            float64  `json:"price" bson:"price"`
	Status           string   `json:"status" bson:"status"`
	Address_Shipping Address  `json:"address_shipping" bson:"address_shipping"`
}

type SellerOrderRepository interface {
	CreateOrderSeller(ctx context.Context, order SellerOrder) (primitive.ObjectID, error)
	GetAllSellerOrders(ctx context.Context, email string) (*[]SellerOrder, error)
	GetSellerOrderById(ctx context.Context, orderID string) (*[]SellerOrder, error)
	GetSellerOrderByEmailAndId(ctx context.Context, email, orderID string) (*SellerOrder, error)
	UpdateOrderSeller(ctx context.Context, orderID string, req *dto.OrderSellerUpdateReq) (*mongo.UpdateResult, error)
	UpdateOrderSellerByEmail(ctx context.Context, email string, req *dto.OrderSellerUpdateReq) (*mongo.UpdateResult, error)
	UpdateStatusOrderSeller(ctx context.Context, orderID, productID string, req *dto.OrderStatusUpdateReq) (*mongo.UpdateResult, error)
	DeleteItem(ctx context.Context, email, orderID, productID string) (*mongo.UpdateResult, error)
}

type SellerOrderService interface {
	GetAllSellerOrders(ctx context.Context, email string) (*[]SellerOrder, error)
	GetSellerOrderByEmailAndId(ctx context.Context, email, orderID string) (*SellerOrder, error)
	UpdateSellerAndUserOrderStatus(ctx context.Context, email, orderID, productID string, req *dto.OrderStatusUpdateReq) error
	CancelOrder(ctx context.Context, email, orderID, productID string) error
}
