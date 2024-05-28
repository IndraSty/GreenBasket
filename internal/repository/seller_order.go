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

type sellerOrderRepository struct {
	Collection *mongo.Collection
}

func NewSellerOrderRepository(client *mongo.Client) domain.SellerOrderRepository {
	return &sellerOrderRepository{
		Collection: db.OpenCollection(client, "Seller_Orders"),
	}
}

// CreateOrderSeller implements domain.SellerOrderRepository.
func (repo *sellerOrderRepository) CreateOrderSeller(ctx context.Context, order domain.SellerOrder) (primitive.ObjectID, error) {
	result, err := repo.Collection.InsertOne(ctx, order)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// GetAllSellerOrders implements domain.SellerOrderRepository.
func (repo *sellerOrderRepository) GetAllSellerOrders(ctx context.Context, email string) (*[]domain.SellerOrder, error) {
	var orders []domain.SellerOrder
	filter := bson.M{"email": email}
	cur, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var order domain.SellerOrder
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

// GetSellerOrderByEmailAndId implements domain.SellerOrderRepository.
func (repo *sellerOrderRepository) GetSellerOrderByEmailAndId(ctx context.Context, email string, orderID string) (*domain.SellerOrder, error) {
	var order domain.SellerOrder
	filter := bson.M{"email": email, "order_id": orderID}
	err := repo.Collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// GetSellerOrderById implements domain.SellerOrderRepository.
func (repo *sellerOrderRepository) GetSellerOrderById(ctx context.Context, orderID string) (*[]domain.SellerOrder, error) {
	var orders []domain.SellerOrder
	filter := bson.M{"order_id": orderID}
	cur, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var order domain.SellerOrder
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

// UpdateOrderSeller implements domain.SellerOrderRepository.
func (repo *sellerOrderRepository) UpdateOrderSeller(ctx context.Context, orderID string, req *dto.OrderSellerUpdateReq) (*mongo.UpdateResult, error) {
	filter := bson.M{"order_id": orderID}
	var update primitive.D

	if req.Payment_Status != "" {
		update = append(update, bson.E{Key: "payment_status", Value: req.Payment_Status})
	}
	if req.Status != "" {
		update = append(update, bson.E{Key: "items.$[].status", Value: req.Status})
	}

	res, err := repo.Collection.UpdateMany(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateOrderSellerByEmail implements domain.SellerOrderRepository.
func (repo *sellerOrderRepository) UpdateOrderSellerByEmail(ctx context.Context, email string, req *dto.OrderSellerUpdateReq) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	var update primitive.D

	if req.Payment_Status != "" {
		update = append(update, bson.E{Key: "payment_status", Value: req.Payment_Status})
	}

	res, err := repo.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateStatusOrderSeller implements domain.SellerOrderRepository.
func (repo *sellerOrderRepository) UpdateStatusOrderSeller(ctx context.Context, orderID, productID string, req *dto.OrderStatusUpdateReq) (*mongo.UpdateResult, error) {
	filter := bson.M{"order_id": orderID, "items.product_id": productID}
	var update primitive.D

	if req.Status != "" {
		update = append(update, bson.E{Key: "items.$.status", Value: req.Status})
	}

	res, err := repo.Collection.UpdateMany(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteItem implements domain.SellerOrderRepository.
func (repo *sellerOrderRepository) DeleteItem(ctx context.Context, email, orderID string, productID string) (*mongo.UpdateResult, error) {
	filter := bson.M{"order_id": orderID, "email": email}
	update := bson.M{
		"$pull": bson.M{
			"items": bson.M{
				"product_id": productID,
			},
		},
	}

	return repo.Collection.UpdateOne(ctx, filter, update)
}
