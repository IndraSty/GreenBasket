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

type userRepository struct {
	Collection *mongo.Collection
}

func NewUserRepository(client *mongo.Client) domain.UserRepository {
	return &userRepository{
		Collection: db.OpenCollection(client, "Users"),
	}
}

func (ur *userRepository) CreateUser(ctx context.Context, user domain.User) (primitive.ObjectID, error) {
	insertResult, err := ur.Collection.InsertOne(ctx, user)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return insertResult.InsertedID.(primitive.ObjectID), nil
}

func (ur *userRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	count, err := ur.Collection.CountDocuments(ctx, bson.M{"email": email})
	return count > 0, err
}

func (ur *userRepository) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	count, err := ur.Collection.CountDocuments(ctx, bson.M{"phone": phone})
	return count > 0, err
}

func (ur *userRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := ur.Collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserById implements domain.UserRepository.
func (ur *userRepository) FindUserById(ctx context.Context, userId string) (*domain.User, error) {
	var user domain.User
	err := ur.Collection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) UpdateUser(ctx context.Context, email string, update bson.D) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	opt := options.Update().SetUpsert(true)
	return ur.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}}, opt)
}
