package repository

import (
	"context"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type sellerRepository struct {
	Collection *mongo.Collection
}

func NewSellerRepository(client *mongo.Client) domain.SellerRepository {
	return &sellerRepository{
		Collection: db.OpenCollection(client, "Sellers"),
	}
}

func (sr *sellerRepository) CreateSeller(ctx context.Context, seller domain.Seller) (primitive.ObjectID, error) {
	insertResult, err := sr.Collection.InsertOne(ctx, seller)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return insertResult.InsertedID.(primitive.ObjectID), nil
}

func (sr *sellerRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	count, err := sr.Collection.CountDocuments(ctx, bson.M{"email": email})
	return count > 0, err
}

func (sr *sellerRepository) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	count, err := sr.Collection.CountDocuments(ctx, bson.M{"phone": phone})
	return count > 0, err
}

func (sr *sellerRepository) FindSellerByEmail(ctx context.Context, email string) (*domain.Seller, error) {
	var seller domain.Seller
	err := sr.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&seller)
	if err != nil {
		return nil, err
	}

	return &seller, nil
}

func (sr *sellerRepository) FindSellerByStoreId(ctx context.Context, storeID string) (*domain.Seller, error) {
	var seller domain.Seller
	err := sr.Collection.FindOne(ctx, bson.M{"store_id": storeID}).Decode(&seller)
	if err != nil {
		return nil, err
	}

	return &seller, nil
}

func (sr *sellerRepository) UpdateSeller(ctx context.Context, email string, update bson.D) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	opt := options.Update().SetUpsert(true)
	return sr.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}}, opt)
}

// AddStoreId implements domain.SellerRepository.
func (sr *sellerRepository) AddStoreId(ctx context.Context, email string, storeID string) error {
	filter := bson.M{"email": email}
	_, err := sr.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "store_id", Value: storeID}}}})
	if err != nil {
		return err
	}

	return nil
}
