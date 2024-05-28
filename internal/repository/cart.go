package repository

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type cartRepository struct {
	Collection *mongo.Collection
}

func NewCartRepository(client *mongo.Client) domain.CartRepository {
	return &cartRepository{
		Collection: db.OpenCollection(client, "Cart"),
	}
}

// AddToCart implements domain.CartRepository.
func (repo *cartRepository) AddToCart(ctx context.Context, email string, item *domain.CartItem) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email}
	result, err := repo.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$push", Value: bson.D{{Key: "items", Value: item}}}})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CheckUserCart implements domain.CartRepository.
func (repo *cartRepository) CheckUserCart(ctx context.Context, email string) (bool, error) {
	filter := bson.M{"email": email}
	count, err := repo.Collection.CountDocuments(ctx, filter)
	return count > 0, err
}

// CreateCart implements domain.CartRepository.
func (repo *cartRepository) CreateCart(ctx context.Context, cart *domain.Cart) error {
	_, err := repo.Collection.InsertOne(ctx, cart)
	if err != nil {
		return err
	}

	return nil
}

// GetAllCartItem implements domain.CartService.
func (repo *cartRepository) GetAllCartItem(ctx context.Context, email string) (*[]domain.CartItem, error) {
	var cart domain.Cart
	filter := bson.M{"email": email}
	err := repo.Collection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		return nil, err
	}

	return &cart.Items, nil
}

// GetUserCart implements domain.CartRepository.
func (repo *cartRepository) GetUserCart(ctx context.Context, email string) (*domain.Cart, error) {
	var cart domain.Cart
	filter := bson.M{"email": email}
	err := repo.Collection.FindOne(ctx, filter).Decode(&cart)
	if err != nil {
		return nil, err
	}

	return &cart, nil
}

// RemoveCartItemById implements domain.CartRepository.
func (repo *cartRepository) RemoveCartItemById(ctx context.Context, email string, productID string, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email, "items.product_id": productID}
	update := bson.M{
		"$pull": bson.M{"items": bson.M{"product_id": productID}},
		"$set":  bson.M{"updated_at": time.Now()},
	}
	result, err := repo.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateCartItemById implements domain.CartRepository.
func (repo *cartRepository) UpdateCartItemById(ctx context.Context, email string, productID string, value *dto.CartItemEditRepo, updateAt time.Time) (*mongo.UpdateResult, error) {
	filter := bson.M{"email": email, "items.product_id": productID}
	update := bson.M{
		"$set": bson.M{
			"items.$.quantity": value.Quantity,
			"items.$.selected": value.Selected,
			"updated_at":       time.Now(),
		},
		"$inc": bson.M{"total_price": value.Total_Price},
	}

	result, err := repo.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, err
}

// UpdateTotalPrice implements domain.CartRepository.
func (repo *cartRepository) UpdateTotalPrice(ctx context.Context, email string, value float64) error {
	filter := bson.M{"email": email}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "total_price", Value: value}}}}

	_, err := repo.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
