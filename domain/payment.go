package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Payment struct {
	ID             primitive.ObjectID `bson:"_id"`
	OrderID        string             `json:"order_id" bson:"order_id"`
	UserID         string             `json:"user_id" bson:"user_id"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdateAt       time.Time          `json:"updated_at" bson:"updated_at"`
	Payment_Method string             `json:"payment_method" bson:"payment_method"`
	Amount         float64            `json:"amount" bson:"amount"`
	Status         string             `json:"status" bson:"status"`
	TransactionID  string             `json:"transaction_id" bson:"transaction_id"`
	Snap_Url       string             `json:"snap_url"`
}

type PaymentRepository interface {
	FindByOrderId(ctx context.Context, orderID string) (*Payment, error)
	Insert(ctx context.Context, p *Payment) error
	Update(ctx context.Context, orderID string, req *dto.UpdatePaymentReq) (*mongo.UpdateResult, error)
}

type PaymentService interface {
	ConfirmedPayment(ctx context.Context, orderID string) error
	InitializePayment(ctx context.Context, req *dto.PaymentReq) (*dto.PaymentRes, error)
	UpdatePayment(ctx context.Context, orderID string, req *dto.UpdatePaymentReq) error
}
