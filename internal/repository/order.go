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

type orderRepository struct {
	Collection *mongo.Collection
}

func NewOrderRepository(client *mongo.Client) domain.OrderRepository {
	return &orderRepository{
		Collection: db.OpenCollection(client, "Orders"),
	}
}

// CreateOrder implements domain.OrderRepository.
func (repo *orderRepository) CreateOrder(ctx context.Context, order domain.Orders) (primitive.ObjectID, error) {
	result, err := repo.Collection.InsertOne(ctx, order)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// GetAllOrders implements domain.OrderRepository.
func (repo *orderRepository) GetAllOrders(ctx context.Context, email string) (*[]domain.Orders, error) {
	var orders []domain.Orders
	filter := bson.M{"email": email}
	cur, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var order domain.Orders
		err := cur.Decode(&order)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &orders, nil
}

// GetOrderById implements domain.OrderRepository.
func (repo *orderRepository) GetOrder(ctx context.Context, orderID string, email ...string) (*domain.Orders, error) {
	var order domain.Orders
	filter := bson.M{"order_id": orderID}

	if len(email) > 0 {
		filter["email"] = email[0]
	}

	err := repo.Collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// Update implements domain.OrderRepository.
func (repo *orderRepository) UpdateOrder(ctx context.Context, orderID string, req *dto.UpdatePaymentReq) (*mongo.UpdateResult, error) {
	filter := bson.M{"order_id": orderID}
	var update primitive.D

	if req.Payment_Method != "" {
		update = append(update, bson.E{Key: "payment.payment_method", Value: req.Payment_Method})
	}
	if req.Status != "" {
		update = append(update, bson.E{Key: "payment.status", Value: req.Status})
	}
	if req.TransactionID != "" {
		update = append(update, bson.E{Key: "payment.transaction_id", Value: req.TransactionID})
	}

	res, err := repo.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateStatusOrder implements domain.OrderRepository.
func (repo *orderRepository) UpdateStatusOrder(ctx context.Context, orderID string, productID string, req *dto.OrderStatusUpdateReq) (*mongo.UpdateResult, error) {
	filter := bson.M{"order_id": orderID, "items.product_id": productID}
	var update primitive.D

	if req.Status != "" {
		update = append(update, bson.E{Key: "items.$.order_status", Value: req.Status})
	}

	res, err := repo.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteItem implements domain.OrderRepository.
func (repo *orderRepository) DeleteItem(ctx context.Context, orderID string, productID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"order_id": orderID}
	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{
				"product_id": productID,
			},
		},
	}

	return repo.Collection.UpdateOne(ctx, filter, update)
}
