package repository

import (
	"context"
	"errors"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type storeRepository struct {
	Collection *mongo.Collection
}

func NewStoreRepository(client *mongo.Client) domain.StoreRepository {
	return &storeRepository{
		Collection: db.OpenCollection(client, "Stores"),
	}
}

func (repo *storeRepository) CreateStore(ctx context.Context, store domain.Store) (primitive.ObjectID, error) {
	insertResult, err := repo.Collection.InsertOne(ctx, store)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertResult.InsertedID.(primitive.ObjectID), nil
}

// GetStore implements domain.StoreRepository.
func (repo *storeRepository) GetStore(ctx context.Context, storeID string, email ...string) (*domain.Store, error) {
	var store domain.Store
	filter := bson.M{"store_id": storeID}

	if len(email) > 0 {
		filter["email"] = email[0]
	}

	err := repo.Collection.FindOne(ctx, filter).Decode(&store)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

// GetStoreByEmail implements domain.StoreRepository.
func (repo *storeRepository) GetStoreByEmail(ctx context.Context, email string) (*domain.Store, error) {
	var store domain.Store
	filter := bson.M{"email": email}
	err := repo.Collection.FindOne(ctx, filter).Decode(&store)
	if err != nil {
		return nil, err
	}

	return &store, nil
}

// GetStoreByQuery implements domain.StoreRepository.
func (repo *storeRepository) GetStoreByQuery(ctx context.Context, query string) ([]domain.Store, error) {
	var stores []domain.Store
	filter := bson.M{}
	if query != "" {
		filter["name"] = bson.M{
			"$regex": primitive.Regex{
				Pattern: query,
				Options: "i",
			},
		}
	}

	cur, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var store domain.Store
		err := cur.Decode(&store)
		if err != nil {
			return nil, err
		}
		stores = append(stores, store)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return stores, nil
}

// CheckNameExists implements domain.StoreRepository.
func (repo *storeRepository) CheckNameExists(ctx context.Context, name string) (bool, error) {
	count, err := repo.Collection.CountDocuments(ctx, bson.M{"name": name})
	return count > 0, err
}

func (repo *storeRepository) UpdateStore(ctx context.Context, email, storeID string, update bson.D) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email, "store_id": storeID}
	return repo.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
}

func (repo *storeRepository) RemoveStore(ctx context.Context, email, storeID string) (*mongo.DeleteResult, error) {
	filter := bson.M{"email": email, "store_id": storeID}
	result, err := repo.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("no store was deleted")
	}

	return result, nil
}
