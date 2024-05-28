package repository

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type addressRepository struct {
	userCol   *mongo.Collection
	storeCol  *mongo.Collection
	sellerCol *mongo.Collection
}

func NewAddressRepository(client *mongo.Client) domain.AddressRepository {
	return &addressRepository{
		userCol:   db.OpenCollection(client, "Users"),
		storeCol:  db.OpenCollection(client, "Stores"),
		sellerCol: db.OpenCollection(client, "Sellers"),
	}
}

func (rep *addressRepository) AddUserAddress(ctx context.Context, email string, address domain.Address, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	return rep.userCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: address}, {Key: "updated_at", Value: updateAt}}}})
}

func (rep *addressRepository) GetUserAddress(ctx context.Context, email string) (*domain.Address, error) {
	var user domain.User
	var address domain.Address
	filter := bson.M{"email": email}
	err := rep.userCol.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	address = *user.Address_Details

	return &address, nil
}

func (rep *addressRepository) UpdateAddressUser(ctx context.Context, email string, address domain.Address, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	return rep.userCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: address}, {Key: "updated_at", Value: updateAt}}}})
}

func (rep *addressRepository) RemoveAddressUser(ctx context.Context, email string, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	return rep.userCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: nil}, {Key: "updated_at", Value: updateAt}}}})
}

func (rep *addressRepository) AddSellerAddress(ctx context.Context, email string, address domain.Address, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	return rep.sellerCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: address}, {Key: "updated_at", Value: updateAt}}}})
}

func (rep *addressRepository) GetSellerAddress(ctx context.Context, email string) (*domain.Address, error) {
	var seller domain.Seller
	var address domain.Address
	filter := bson.M{"email": email}
	err := rep.sellerCol.FindOne(ctx, filter).Decode(&seller)
	if err != nil {
		return nil, err
	}

	address = *seller.Address_Details

	return &address, nil
}

func (rep *addressRepository) UpdateSellerAddress(ctx context.Context, email string, address domain.Address, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	return rep.sellerCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: address}, {Key: "updated_at", Value: updateAt}}}})
}

func (rep *addressRepository) RemoveSellerAddress(ctx context.Context, email string, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	return rep.sellerCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: nil}, {Key: "updated_at", Value: updateAt}}}})
}

func (rep *addressRepository) AddStoreAddress(ctx context.Context, email string, storeID string, address domain.Address, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email, "store_id": storeID}
	return rep.storeCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: address}, {Key: "updated_at", Value: updateAt}}}})
}

func (rep *addressRepository) GetStoreAddress(ctx context.Context, email string, storeID string) (*domain.Address, error) {
	var store domain.Store
	var address domain.Address
	filter := bson.M{"email": email, "store_id": storeID}
	err := rep.storeCol.FindOne(ctx, filter).Decode(&store)
	if err != nil {
		return nil, err
	}

	address = *store.Address_Details

	return &address, nil
}

func (rep *addressRepository) UpdateStoreAddress(ctx context.Context, email string, storeID string, address domain.Address, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email, "store_id": storeID}
	return rep.storeCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: address}, {Key: "updated_at", Value: updateAt}}}})
}

func (rep *addressRepository) RemoveStoreAddress(ctx context.Context, email string, storeID string, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email, "store_id": storeID}
	return rep.storeCol.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: bson.D{{Key: "address", Value: nil}, {Key: "updated_at", Value: updateAt}}}})
}
