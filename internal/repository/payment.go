package repository

import (
	"context"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type paymentRepository struct {
	Collection *mongo.Collection
}

func NewPaymentRepository(client *mongo.Client) domain.PaymentRepository {
	return &paymentRepository{
		Collection: db.OpenCollection(client, "Payments"),
	}
}

// FindById implements domain.PaymentRepository.
func (r *paymentRepository) FindByOrderId(ctx context.Context, orderID string) (*domain.Payment, error) {
	var payment domain.Payment
	filter := bson.M{"order_id": orderID}
	err := r.Collection.FindOne(ctx, filter).Decode(&payment)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

// Insert implements domain.PaymentRepository.
func (r *paymentRepository) Insert(ctx context.Context, p *domain.Payment) error {
	_, err := r.Collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

// Update implements domain.PaymentRepository.
func (r *paymentRepository) Update(ctx context.Context, orderID string, req *dto.UpdatePaymentReq) (*mongo.UpdateResult, error) {
	filter := bson.M{"order_id": orderID}
	var update primitive.D

	if req.Payment_Method != "" {
		update = append(update, bson.E{Key: "payment_method", Value: req.Payment_Method})
	}
	if req.Status != "" {
		update = append(update, bson.E{Key: "status", Value: req.Status})
	}
	if req.TransactionID != "" {
		update = append(update, bson.E{Key: "transaction_id", Value: req.TransactionID})
	}

	res, err := r.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}

	return res, nil
}
