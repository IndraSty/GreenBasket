package repository

import (
	"context"

	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type reviewRepository struct {
	Collection *mongo.Collection
}

func NewReviewRepository(client *mongo.Client) domain.ReviewRepository {
	return &reviewRepository{
		Collection: db.OpenCollection(client, "Reviews"),
	}
}

// GetReview implements domain.ReviewRepository.
func (repo *reviewRepository) GetReviewById(ctx context.Context, reviewID string) (*domain.Review, error) {
	var review domain.Review
	filter := bson.M{"review_id": reviewID}
	err := repo.Collection.FindOne(ctx, filter).Decode(&review)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

// GetReviewByProductId implements domain.ReviewRepository.
func (repo *reviewRepository) GetReviewByProductId(ctx context.Context, productID string) (*domain.Review, error) {
	var review domain.Review
	filter := bson.M{"product_id": productID}
	err := repo.Collection.FindOne(ctx, filter).Decode(&review)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

// GetReviewByEmailAndId implements domain.ReviewRepository.
func (repo *reviewRepository) GetUserReviewByEmailAndId(ctx context.Context, email string, reviewID string) (*domain.Review, error) {
	var review domain.Review
	filter := bson.M{"review_id": reviewID, "email": email}
	err := repo.Collection.FindOne(ctx, filter).Decode(&review)
	if err != nil {
		return nil, err
	}

	return &review, nil
}

// InsertReview implements domain.ReviewRepository.
func (repo *reviewRepository) InsertReview(ctx context.Context, input domain.Review) (primitive.ObjectID, error) {
	result, err := repo.Collection.InsertOne(ctx, input)
	if err != nil {
		return primitive.NilObjectID, nil
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

// UpdateReview implements domain.ReviewRepository.
func (repo *reviewRepository) UpdateReview(ctx context.Context, reviewID string, update primitive.D) (*mongo.UpdateResult, error) {
	filter := bson.M{"review_id": reviewID}
	result, err := repo.Collection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteReview implements domain.ReviewRepository.
func (repo *reviewRepository) DeleteReview(ctx context.Context, reviewID string) (*mongo.DeleteResult, error) {
	filter := bson.M{"review_id": reviewID}
	result, err := repo.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetAllReviewByProductId implements domain.ReviewRepository.
func (repo *reviewRepository) GetAllReviewByProductId(ctx context.Context, productID, sellerEmail string) (*[]domain.Review, error) {
	var reviews []domain.Review
	filter := bson.M{"product_id": productID, "seller_email": sellerEmail}
	cur, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var review domain.Review
		err := cur.Decode(&review)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &reviews, nil
}

// GetAllReviewBySellerEmail implements domain.ReviewRepository.
func (repo *reviewRepository) GetAllReviewBySellerEmail(ctx context.Context, sellerEmail string) (*[]domain.Review, error) {
	var reviews []domain.Review
	filter := bson.M{"seller_email": sellerEmail}

	cur, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var review domain.Review
		err := cur.Decode(&review)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &reviews, nil
}

// GetAllReviewByUserEmail implements domain.ReviewRepository.
func (repo *reviewRepository) GetAllReviewByUserEmail(ctx context.Context, email string) (*[]domain.Review, error) {
	var reviews []domain.Review
	filter := bson.M{"email": email}

	cur, err := repo.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var review domain.Review
		err := cur.Decode(&review)
		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return &reviews, nil
}
