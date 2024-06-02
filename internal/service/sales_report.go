package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type salesReportService struct {
	repo            domain.SalesReportRepository
	sellerOrderRepo domain.SellerOrderRepository
	storeRepo       domain.StoreRepository
	productRepo     domain.ProductRepository
	reviewRepo      domain.ReviewRepository
	cacheRepo       domain.CacheRepository
}

func NewSalesRepository(repo domain.SalesReportRepository, sellerOrderRepo domain.SellerOrderRepository,
	storeRepo domain.StoreRepository, productRepo domain.ProductRepository,
	reviewRepo domain.ReviewRepository, cacheRepo domain.CacheRepository) domain.SalesReportService {
	return &salesReportService{
		repo:            repo,
		sellerOrderRepo: sellerOrderRepo,
		storeRepo:       storeRepo,
		productRepo:     productRepo,
		reviewRepo:      reviewRepo,
		cacheRepo:       cacheRepo,
	}
}

func calculateSalesAndIncome(orders []domain.SellerOrder) (totalSales int32, totalIncome float64) {
	for _, order := range orders {
		if order.Payment_Status == "SUCCESS" {
			for _, item := range order.Items {
				if item.Status == "FINISHED" {
					totalSales += int32(item.Quantity)
					totalIncome += float64(item.Quantity) * item.Price
				}
			}
		} else {
			log.Println("payment status not 'SUCCESS'")
		}
	}
	return
}

func calculateProductSales(orders []domain.SellerOrder) map[string]int64 {
	productSalesMap := make(map[string]int64)

	for _, order := range orders {
		if order.Payment_Status == "SUCCESS" {
			for _, item := range order.Items {
				if item.Status == "FINISHED" {
					productSalesMap[item.Product_Id] += int64(item.Quantity)
				}
			}
		} else {
			log.Println("payment status not 'SUCCESS'")
		}
	}

	return productSalesMap
}

func CalculateAverageRatingProduct(reviews []domain.Review, productId string) float32 {
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

func (s *salesReportService) setRedisSR(data dto.SalesReportRes, email string) error {
	spData, err := json.Marshal(data)
	if err != nil {
		return errors.New("failed to marshal sales report data: " + err.Error())
	}

	err = s.cacheRepo.Set("sales-report_seller:"+email, spData, time.Hour*24)
	if err != nil {
		return errors.New("failed to set sales report data in cache: " + err.Error())
	}

	return nil
}

func (s *salesReportService) delRedisSR(email string) error {
	err := s.cacheRepo.Del("sales-report_seller:" + email)
	if err != nil {
		return errors.New("failed to delete sales report data in cache: " + err.Error())
	}

	return nil
}

func (s *salesReportService) updateRedisSR(ctx context.Context, email, storeID string) error {
	cartItems, err := s.getSalesReportWithNoAct(ctx, email, storeID)
	if err != nil {
		return errors.New("failed to get sales report data: " + err.Error())
	}

	return s.setRedisSR(*cartItems, email)
}

func (s *salesReportService) getSalesReportWithNoAct(ctx context.Context, email string, storeID string) (*dto.SalesReportRes, error) {
	var productSales []dto.ProductSalesRes
	_, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil {
		return nil, errors.New("failed to get store by email and id: " + err.Error())
	}

	result, err := s.repo.GetByEmailAndStoreId(ctx, email, storeID)
	if err != nil {
		return nil, errors.New("failed to get sales report: " + err.Error())
	}

	for _, item := range result.Products {
		result := dto.ProductSalesRes{
			Product_Id:  item.Product_Id,
			Total_Sales: item.Total_Sales,
			Stock:       item.Stock,
		}

		productSales = append(productSales, result)

	}

	return &dto.SalesReportRes{
		Store_Id:      result.Store_Id,
		Email:         result.Email,
		Total_Sales:   result.Total_Sales,
		Total_Incomes: result.Total_Incomes,
		Products:      productSales,
	}, nil
}

// UpdateSalesReport implements domain.SalesReportService.
func (s *salesReportService) UpdateSalesReport(ctx context.Context, storeID string, email string) error {
	err := s.delRedisSR(email)
	if err != nil {
		return err
	}

	_, err = s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil {
		return errors.New("failed to get store by email and id: " + err.Error())
	}

	orders, err := s.sellerOrderRepo.GetAllSellerOrders(ctx, email)
	if err != nil {
		return errors.New("failed to get all seller order with email: " + err.Error())
	}

	totalSales, totalIncome := calculateSalesAndIncome(*orders)
	productSalesMap := calculateProductSales(*orders)

	var productSales []domain.Product_Sales
	for productID, totalProdSales := range productSalesMap {
		product, err := s.productRepo.GetProductById(ctx, productID)
		if err != nil || product == nil {
			return errors.New("failed to get product by id: " + err.Error())
		}

		reviews, err := s.reviewRepo.GetAllReviewByProductId(ctx, productID, email)
		if err != nil {
			return errors.New("failed to get all reviews: " + err.Error())
		}

		averageRating := CalculateAverageRatingProduct(*reviews, productID)

		productSales = append(productSales, domain.Product_Sales{
			Product_Id:     productID,
			Total_Sales:    totalProdSales,
			Stock:          int32(product.Stock),
			Average_Rating: averageRating,
		})
	}

	var update primitive.D
	if totalSales != 0 {
		update = append(update, bson.E{Key: "total_sales", Value: totalSales})
	}
	if totalIncome != 0.0 {
		update = append(update, bson.E{Key: "total_income", Value: totalIncome})
	}
	if len(productSales) > 0 {
		update = append(update, bson.E{Key: "products", Value: productSales})
	}

	if len(update) == 0 {
		return errors.New("no updates to be made")
	}

	result, err := s.repo.Update(ctx, storeID, update)
	if err != nil {
		return errors.New("failed to update sales report: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return errors.New("no data was updated")
	}

	defer func() {
		if err := s.updateRedisSR(ctx, email, storeID); err != nil {
			log.Println("failed to update sales report data in cache: ", err)
		}
	}()

	return nil
}

// GetSalesReport implements domain.SalesReportService.
func (s *salesReportService) GetSalesReport(ctx context.Context, email string, storeID string) (*dto.SalesReportRes, error) {
	val, err := s.cacheRepo.Get("sales-report_seller:" + email)
	if err == nil {
		var data dto.SalesReportRes
		err = json.Unmarshal(val, &data)
		if err != nil {
			return nil, errors.New("failed to unmarshal sales report data: " + err.Error())
		}
		return &data, nil
	}

	var productSales []dto.ProductSalesRes
	_, err = s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil {
		return nil, errors.New("failed to get store by email and id: " + err.Error())
	}

	result, err := s.repo.GetByEmailAndStoreId(ctx, email, storeID)
	if err != nil {
		return nil, errors.New("failed to get sales report: " + err.Error())
	}

	if result == nil {
		return nil, errors.New("this seller has no store")
	}

	for _, item := range result.Products {
		result := dto.ProductSalesRes{
			Product_Id:  item.Product_Id,
			Total_Sales: item.Total_Sales,
			Stock:       item.Stock,
		}

		productSales = append(productSales, result)

	}

	data := dto.SalesReportRes{
		Store_Id:      result.Store_Id,
		Email:         result.Email,
		Total_Sales:   result.Total_Sales,
		Total_Incomes: result.Total_Incomes,
		Products:      productSales,
	}

	err = s.setRedisSR(data, email)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
