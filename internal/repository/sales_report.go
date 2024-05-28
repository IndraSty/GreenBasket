package repository

import (
	"context"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type salesReportRepository struct {
	Collection *mongo.Collection
}

func NewSalesReportRepository(client *mongo.Client) domain.SalesReportRepository {
	return &salesReportRepository{
		Collection: db.OpenCollection(client, "Sales_Report"),
	}
}

// Insert implements domain.SalesReportRepository.
func (repo *salesReportRepository) Insert(ctx context.Context, input domain.Sales_Report) (primitive.ObjectID, error) {
	result, err := repo.Collection.InsertOne(ctx, input)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// UpdateAverageRating implements domain.SalesReportRepository.
func (repo *salesReportRepository) UpdateAverageRating(ctx context.Context, storeID string, productID string, averageRating float32) (*mongo.UpdateResult, error) {
	filter := bson.M{"store_id": storeID, "products.product_id": productID}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "products.$.average_rating", Value: averageRating}}}}
	return repo.Collection.UpdateOne(ctx, filter, update)
}

// Update implements domain.SalesReportRepository.
func (repo *salesReportRepository) Update(ctx context.Context, storeID string, update primitive.D) (*mongo.UpdateResult, error) {
	filter := bson.M{"store_id": storeID}
	result, err := repo.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetByEmailAndStoreId implements domain.SalesReportRepository.
func (repo *salesReportRepository) GetByEmailAndStoreId(ctx context.Context, email string, storeID string) (*domain.Sales_Report, error) {
	var result domain.Sales_Report
	filter := bson.M{"store_id": storeID, "email": email}
	err := repo.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetByStoreId implements domain.SalesReportRepository.
func (repo *salesReportRepository) GetByStoreId(ctx context.Context, storeID string) (*domain.Sales_Report, error) {
	var result domain.Sales_Report
	filter := bson.M{"store_id": storeID}
	err := repo.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
