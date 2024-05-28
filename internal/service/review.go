package service

import (
	"context"
	"errors"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type reviewService struct {
	repo            domain.ReviewRepository
	productRepo     domain.ProductRepository
	orderRepo       domain.OrderRepository
	storeRepo       domain.StoreRepository
	userRepo        domain.UserRepository
	salesReportRepo domain.SalesReportRepository
	notifSvc        domain.NotificationService
	cacheRepo       domain.CacheRepository
}

func NewReviewService(repo domain.ReviewRepository, productRepo domain.ProductRepository,
	orderRepo domain.OrderRepository, storeRepo domain.StoreRepository,
	notifSvc domain.NotificationService, userRepo domain.UserRepository,
	salesReportRepo domain.SalesReportRepository, cacheRepo domain.CacheRepository) domain.ReviewService {
	return &reviewService{
		repo:            repo,
		productRepo:     productRepo,
		orderRepo:       orderRepo,
		storeRepo:       storeRepo,
		notifSvc:        notifSvc,
		userRepo:        userRepo,
		salesReportRepo: salesReportRepo,
		cacheRepo:       cacheRepo,
	}
}

func calculateAverageRatingProduct(reviews []domain.Review, productId string) float32 {
	var totalRating float32
	var count int

	for _, review := range reviews {
		if review.Product_Id == productId {
			totalRating += review.Rating
			count++
		}
	}

	if count == 0 {
		return 0
	}

	averageRating := totalRating / float32(count)
	return averageRating
}

// CreateReview implements domain.ReviewService.
func (s *reviewService) CreateReview(ctx context.Context, email string, orderId string, productID string, req *dto.AddReviewReq) (*dto.AddReviewRes, error) {
	order, err := s.orderRepo.GetOrder(ctx, orderId, email)
	if err != nil {
		return nil, errors.New("Failed to get order by email and id: " + err.Error())
	}

	var insertid primitive.ObjectID
	var messages string

	for _, item := range order.Items {
		if item.Product_Id == productID {
			if item.Order_Status == "FINISHED" {

				store, err := s.storeRepo.GetStore(ctx, item.StoreID)
				if err != nil {
					return nil, errors.New("Failed to get store by id: " + err.Error())
				}

				user, err := s.userRepo.FindUserByEmail(ctx, order.Email)
				if err != nil {
					return nil, errors.New("Failed to get user by email: " + err.Error())
				}

				id := primitive.NewObjectID()
				reviewID := id.Hex()
				input := domain.Review{
					Id:              id,
					Review_Id:       reviewID,
					Product_Id:      productID,
					Email:           order.Email,
					Rating:          req.Rating,
					Review:          req.Review,
					Seller_Email:    store.Email,
					Seller_Response: "",
					Reviewed_At:     time.Now(),
					Updated_At:      time.Now(),
				}

				res, err := s.repo.InsertReview(ctx, input)
				if err != nil {
					return nil, errors.New("Failed to insert review: " + err.Error())
				}

				review, err := s.repo.GetAllReviewByProductId(ctx, productID, store.Email)
				if err != nil {
					return nil, errors.New("Failed to get all reviews by product id: " + err.Error())
				}

				averageRating := calculateAverageRatingProduct(*review, productID)
				_, err = s.salesReportRepo.UpdateAverageRating(ctx, store.Store_Id, productID, averageRating)
				if err != nil {
					return nil, errors.New("Failed to update sales report: " + err.Error())
				}

				go s.notificationProductReviewed(store.Email, user.First_Name, productID)
				insertid = res
				messages = "Success to Insert Review for product " + item.Product_Id
			}
		}
	}

	return &dto.AddReviewRes{
		InsertId: insertid,
		Messages: messages,
	}, nil
}

// GetUserReviewById implements domain.ReviewService.
func (s *reviewService) GetUserReviewById(ctx context.Context, email string, reviewID string) (*dto.GetReviewRes, error) {
	review, err := s.repo.GetUserReviewByEmailAndId(ctx, email, reviewID)
	if err != nil {
		return nil, errors.New("failed to get review by email and id: " + err.Error())
	}

	if review == nil {
		return nil, errors.New("review not found with that id")
	}

	return &dto.GetReviewRes{
		Review_Id:   review.Review_Id,
		Product_Id:  review.Product_Id,
		Email:       review.Email,
		Rating:      review.Rating,
		Review:      review.Review,
		Reviewed_At: review.Reviewed_At,
		Updated_At:  review.Updated_At,
	}, nil
}

// DeleteReview implements domain.ReviewService.
func (s *reviewService) DeleteReview(ctx context.Context, email string, reviewID string) error {
	review, err := s.repo.GetUserReviewByEmailAndId(ctx, email, reviewID)
	if err != nil {
		return errors.New("failed to get review by email and id: " + err.Error())
	}

	if review == nil {
		return errors.New("review not found with that id")
	}

	_, err = s.repo.DeleteReview(ctx, reviewID)
	if err != nil {
		return errors.New("failed to delete the review: " + err.Error())
	}

	return nil
}

// UpdateReview implements domain.ReviewService.
func (s *reviewService) UpdateReview(ctx context.Context, email string, reviewID string, req *dto.AddReviewReq) error {
	var update primitive.D
	review, err := s.repo.GetUserReviewByEmailAndId(ctx, email, reviewID)
	if err != nil {
		return errors.New("failed to get review by email and id: " + err.Error())
	}

	if review == nil {
		return errors.New("review not found with that id")
	}

	if req.Rating != 0.0 {
		update = append(update, bson.E{Key: "rating", Value: req.Rating})
	}
	if req.Review != "" {
		update = append(update, bson.E{Key: "review", Value: req.Review})
	}

	updateAt := time.Now()
	update = append(update, bson.E{Key: "updated_at", Value: updateAt})

	_, err = s.repo.UpdateReview(ctx, reviewID, update)
	if err != nil {
		return errors.New("failed to update the review: " + err.Error())
	}

	return nil
}

// GetAllReviewByProductId implements domain.ReviewService.
func (s *reviewService) GetAllReviewByProductId(ctx context.Context, productID, sellerEmail string) (*[]dto.GetReviewRes, error) {
	product, err := s.productRepo.GetProductById(ctx, productID)
	if err != nil {
		return nil, errors.New("Failed to get product by id: " + err.Error())
	}

	if product == nil {
		return nil, errors.New("product not found with this id " + productID)
	}

	reviews, err := s.repo.GetAllReviewByProductId(ctx, productID, sellerEmail)
	if err != nil {
		return nil, errors.New("Failed to get all reviews by this product id: " + err.Error())
	}

	reviewRes := make([]dto.GetReviewRes, len(*reviews))
	for i, review := range *reviews {
		reviewRes[i] = dto.GetReviewRes{
			Review_Id:   review.Review_Id,
			Product_Id:  review.Product_Id,
			Email:       review.Email,
			Rating:      review.Rating,
			Review:      review.Review,
			Reviewed_At: review.Reviewed_At,
			Updated_At:  review.Updated_At,
		}
	}

	return &reviewRes, nil
}

// GetAllReviewBySellerEmail implements domain.ReviewService.
func (s *reviewService) GetAllReviewBySellerEmail(ctx context.Context, email string) (*[]dto.GetReviewRes, error) {
	reviews, err := s.repo.GetAllReviewBySellerEmail(ctx, email)
	if err != nil {
		return nil, errors.New("Failed to get all reviews by this seller email: " + err.Error())
	}

	reviewRes := make([]dto.GetReviewRes, len(*reviews))
	for i, review := range *reviews {
		reviewRes[i] = dto.GetReviewRes{
			Review_Id:   review.Review_Id,
			Product_Id:  review.Product_Id,
			Email:       review.Email,
			Rating:      review.Rating,
			Review:      review.Review,
			Reviewed_At: review.Reviewed_At,
			Updated_At:  review.Updated_At,
		}
	}

	return &reviewRes, nil
}

// GetAllReviewByUserEmail implements domain.ReviewService.
func (s *reviewService) GetAllReviewByUserEmail(ctx context.Context, email string) (*[]dto.GetReviewRes, error) {
	reviews, err := s.repo.GetAllReviewByUserEmail(ctx, email)
	if err != nil {
		return nil, errors.New("Failed to get all reviews by this user email: " + err.Error())
	}

	reviewRes := make([]dto.GetReviewRes, len(*reviews))
	for i, review := range *reviews {
		reviewRes[i] = dto.GetReviewRes{
			Review_Id:   review.Review_Id,
			Product_Id:  review.Product_Id,
			Email:       review.Email,
			Rating:      review.Rating,
			Review:      review.Review,
			Reviewed_At: review.Reviewed_At,
			Updated_At:  review.Updated_At,
		}
	}

	return &reviewRes, nil
}

// UpdateResponSeller implements domain.ReviewService.
func (s *reviewService) UpdateResponSeller(ctx context.Context, email string, reviewID string, req *dto.ResponSellerReq) error {
	review, err := s.repo.GetReviewById(ctx, reviewID)
	if err != nil {
		return errors.New("Failed to get review by id: " + err.Error())
	}

	if review == nil {
		return errors.New("review not found with this id: " + reviewID)
	}

	if review.Seller_Response != "" {
		return errors.New("this review already has a seller response")
	}

	if review.Seller_Email != email {
		return errors.New("you do not have the right to provide a response to this review")
	}
	var update primitive.D

	if req.Seller_Response != "" {
		update = append(update, bson.E{Key: "seller_response", Value: req.Seller_Response})
	}
	updateAt := time.Now()
	update = append(update, bson.E{Key: "updated_at", Value: updateAt})

	res, err := s.repo.UpdateReview(ctx, reviewID, update)
	if err != nil {
		return errors.New("failed to update the seller response in review: " + err.Error())
	}

	if res.ModifiedCount == 0 {
		return errors.New("failed to update the seller response in review")
	}

	store, err := s.storeRepo.GetStoreByEmail(ctx, email)
	if err != nil {
		return errors.New("failed to get store by id: " + err.Error())
	}

	err = s.notificationSellerResponse(review.Email, store.Name, review.Product_Id)
	if err != nil {
		return errors.New("failed to insert user notification :" + err.Error())
	}
	return nil
}

func (s *reviewService) notificationProductReviewed(sellerEmail, username, productID string) error {
	data := map[string]string{
		"product_id": productID,
		"username":   username,
	}
	go s.notifSvc.Insert(context.Background(), sellerEmail, "SELLER_PRODUCT_REVIEWED", data)

	return nil
}

func (s *reviewService) notificationSellerResponse(userEmail, storeName, productID string) error {
	data := map[string]string{
		"store_name": storeName,
		"product_id": productID,
	}
	err := s.notifSvc.Insert(context.Background(), userEmail, "SELLER_RESPONSE_REVIEWED", data)
	if err != nil {
		return errors.New("failed to insert user notification :" + err.Error())
	}

	return nil
}
