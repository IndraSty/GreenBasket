package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Orders struct {
	ID               primitive.ObjectID `bson:"_id"`
	Order_id         string             `json:"order_id" bson:"order_id"`
	Email            string             `json:"email" bson:"email"`
	Order_Date       time.Time          `json:"order_date" bson:"order_date"`
	Updated_At       time.Time          `json:"updated_at" bson:"updated_at"`
	Total_Price      float64            `json:"total_price" bson:"total_price"`
	Address_Shipping Address            `json:"address_shipping" bson:"address_shipping"`
	Payment          *PaymentOrder      `json:"payment" bson:"payment"`
	Items            []OrderItem        `json:"items" bson:"items"`
}

type PaymentOrder struct {
	Status         string `json:"status" bson:"status"`
	TransactionID  string `json:"transaction_id" bson:"transaction_id"`
	Payment_Method string `json:"payment_method" bson:"payment_method"`
}

type OrderItem struct {
	Product_Id    string   `json:"product_id" bson:"product_id"`
	Product_Name  string   `json:"product_name" bson:"product_name"`
	Product_Image []string `json:"product_image" bson:"product_image"`
	StoreID       string   `json:"store_id" bson:"store_id"`
	Order_Status  string   `json:"order_status" bson:"order_status"`
	Quantity      int      `json:"quantity" bson:"quantity"`
	Price         float64  `json:"price" bson:"price"`
}

type OrderRepository interface {
	CreateOrder(ctx context.Context, order Orders) (primitive.ObjectID, error)
	GetAllOrders(ctx context.Context, email string) (*[]Orders, error)
	GetOrder(ctx context.Context, orderID string, email ...string) (*Orders, error)
	UpdateOrder(ctx context.Context, orderID string, req *dto.UpdatePaymentReq) (*mongo.UpdateResult, error)
	UpdateStatusOrder(ctx context.Context, orderID, productID string, req *dto.OrderStatusUpdateReq) (*mongo.UpdateResult, error)
	DeleteItem(ctx context.Context, orderID, productID string) (*mongo.UpdateResult, error)
}

type OrderService interface {
	CreateOrder(ctx context.Context, email string) (*dto.InsertOrderRes, error)
	GetAllOrders(ctx context.Context, email string) (*[]Orders, error)
	GetOrderByEmailAndId(ctx context.Context, email, orderID string) (*Orders, error)
	FinishOrder(ctx context.Context, email, orderID, productID string, req *dto.OrderStatusUpdateReq) error
	CancelOrder(ctx context.Context, email, orderID, productID string) error
}
