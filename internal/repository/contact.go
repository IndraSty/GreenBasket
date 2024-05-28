package repository

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type contactRepository struct {
	Collection *mongo.Collection
}

func NewContactRepository(client *mongo.Client) domain.ContactRepository {
	return &contactRepository{
		Collection: db.OpenCollection(client, "Stores"),
	}
}

func (rep *contactRepository) AddStoreContact(ctx context.Context, email string, storeID string, contact domain.Contact, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email, "store_id": storeID}
	return rep.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "contact", Value: contact}, {Key: "updated_at", Value: updateAt}}}})
}

func (rep *contactRepository) GetStoreContact(ctx context.Context, email string, storeID string) (*domain.Contact, error) {
	var store domain.Store
	var contact domain.Contact
	filter := bson.M{"email": email, "store_id": storeID}
	err := rep.Collection.FindOne(ctx, filter).Decode(&store)
	if err != nil {
		return nil, err
	}

	contact = *store.Contact_Details

	return &contact, nil
}

func (rep *contactRepository) UpdateStoreContact(ctx context.Context, email string, storeID string, contact domain.Contact, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email, "store_id": storeID}
	result, err := rep.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "contact", Value: contact}, {Key: "updated_at", Value: updateAt}}}})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (rep *contactRepository) RemoveStoreContact(ctx context.Context, email string, storeID string, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email, "store_id": storeID}
	result, err := rep.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "contact", Value: nil}, {Key: "updated_at", Value: updateAt}}}})
	if err != nil {
		return nil, err
	}

	return result, nil
}
