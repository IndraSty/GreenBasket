package service

import (
	"context"
	"errors"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type productService struct {
	repo            domain.ProductRepository
	storeRepo       domain.StoreRepository
	salesReportRepo domain.SalesReportRepository
	cacheRepo       domain.CacheRepository
}

func NewProductService(repo domain.ProductRepository, storeRepo domain.StoreRepository,
	salesReportRepo domain.SalesReportRepository,
	cacheRepo domain.CacheRepository) domain.ProductService {
	return &productService{
		repo:            repo,
		storeRepo:       storeRepo,
		salesReportRepo: salesReportRepo,
		cacheRepo:       cacheRepo,
	}
}

var validCategories = []string{"Vegetables", "Fruits", "Protein", "Ready to Eat", "Staples",
	"Snacks", "Mother & Baby", "Spices", "Milk & Dairy", "Breakfast"}

// CreateProduct implements domain.ProductService.
func (s *productService) CreateProduct(ctx context.Context, storeID, email string, req *dto.ProductReq) (*dto.AddProductRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil {
		return nil, errors.New("failed to find store" + err.Error())
	}

	if store == nil {
		return nil, errors.New("store not found")
	}

	isValidCategory := false
	for _, category := range validCategories {
		if category == req.Category {
			isValidCategory = true
			break
		}
	}

	if !isValidCategory {
		return nil, errors.New("invalid category")
	}

	products, err := s.repo.GetAllProductWithNoPage(ctx, storeID)
	if err != nil {
		return nil, errors.New("failed to get all product in this store" + err.Error())
	}

	for _, item := range *products {
		if item.Name == req.Name {
			return nil, errors.New("name product already added")
		}
	}

	id := primitive.NewObjectID()
	productID := id.Hex()
	var num = util.ToFixed(req.Price, 2)
	req.Price = num

	product := domain.Products{
		ID:          id,
		Product_id:  productID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stok,
		Category:    req.Category,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
		Store_id:    storeID,
		Images:      req.Images,
	}

	result, err := s.repo.CreateProduct(ctx, product)
	if err != nil {
		return nil, errors.New("failed to created store" + err.Error())
	}

	return &dto.AddProductRes{
		InsertId: &result,
	}, nil
}

// DeleteProductById implements domain.ProductService.
func (s *productService) DeleteProductById(ctx context.Context, storeID, email, productID string) (*dto.DeleteProductRes, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	product, err := s.repo.GetProductById(ctx, productID, storeID)
	if err != nil || product == nil {
		return nil, errors.New("product not found" + err.Error())
	}

	result, err := s.repo.DeleteProductById(ctx, storeID, productID)
	if err != nil {
		return nil, errors.New("failed to delete product: " + err.Error())
	}

	return &dto.DeleteProductRes{
		DeleteResult: result,
	}, nil
}

// GetAllProductSeller implements domain.ProductService.
func (s *productService) GetAllProduct(ctx context.Context, storeID, email string, page int) (*dto.PagedProducts, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	products, err := s.repo.GetAllProduct(ctx, page, storeID)
	if err != nil {
		return nil, errors.New("failed to get all products: " + err.Error())
	}

	productRes := make([]dto.GetProductRes, len(products.Products))
	for i, product := range products.Products {
		averageRating := float32(0)
		totalSales := int64(0)
		if product.SalesData != nil {
			averageRating = product.SalesData.Average_rating
			totalSales = product.SalesData.Total_sales
		}

		productRes[i] = dto.GetProductRes{
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			Stok:           product.Stock,
			Product_id:     product.Product_id,
			Category:       product.Category,
			Created_at:     product.Created_at,
			Images:         product.Images,
			Store_Name:     store.Name,
			City:           store.Address_Details.City,
			Average_Rating: averageRating,
			Total_Sales:    totalSales,
		}
	}

	return &dto.PagedProducts{
		Products:  productRes,
		Page:      page,
		TotalItem: products.TotalItem,
		LastPage:  products.LastPage,
	}, nil
}

// GetAllProductSellerByCategory implements domain.ProductService.
func (s *productService) GetAllByCategory(ctx context.Context, email, storeID string, category string, page int) (*dto.PagedProducts, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	products, err := s.repo.GetAllByCategory(ctx, category, page, storeID)
	if err != nil {
		return nil, errors.New("failed to get all products: " + err.Error())
	}

	productRes := make([]dto.GetProductRes, len(products.Products))
	for i, product := range products.Products {
		store, err := s.storeRepo.GetStore(ctx, product.Store_id)
		if err != nil {
			return nil, errors.New("failed to get store by id: " + err.Error())
		}

		averageRating := float32(0)
		totalSales := int64(0)
		if product.SalesData != nil {
			averageRating = product.SalesData.Average_rating
			totalSales = product.SalesData.Total_sales
		}

		productRes[i] = dto.GetProductRes{
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			Stok:           product.Stock,
			Product_id:     product.Product_id,
			Category:       product.Category,
			Created_at:     product.Created_at,
			Images:         product.Images,
			Store_Name:     store.Name,
			City:           store.Address_Details.City,
			Average_Rating: averageRating,
			Total_Sales:    totalSales,
		}
	}

	return &dto.PagedProducts{
		Products:  productRes,
		Page:      page,
		TotalItem: products.TotalItem,
		LastPage:  products.LastPage,
	}, nil
}

// GetProductById implements domain.ProductService.
func (s *productService) GetProductById(ctx context.Context, storeID, email, productID string) (*dto.GetProductRes, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	product, err := s.repo.GetProductById(ctx, productID, storeID)
	if err != nil {
		return nil, errors.New("failed to get the product: " + err.Error())
	}

	if product == nil {
		return nil, errors.New("product not found")
	}

	averageRating := float32(0)
	totalSales := int64(0)
	if product.SalesData != nil {
		averageRating = product.SalesData.Average_rating
		totalSales = product.SalesData.Total_sales
	}

	return &dto.GetProductRes{
		Name:           product.Name,
		Description:    product.Description,
		Price:          product.Price,
		Stok:           product.Stock,
		Product_id:     product.Product_id,
		Category:       product.Category,
		Created_at:     product.Created_at,
		Store_Name:     store.Name,
		Average_Rating: averageRating,
		Total_Sales:    totalSales,
		City:           store.Address_Details.City,
		Images:         product.Images,
	}, nil
}

// SearchSellerProduct implements domain.ProductService.
func (s *productService) SearchProduct(ctx context.Context, email, storeID, query string, page int) (*dto.PagedProducts, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	products, err := s.repo.GetAllProductByQuery(ctx, query, page, storeID)
	if err != nil {
		return nil, errors.New("failed to get all products by the query: " + err.Error())
	}

	productRes := make([]dto.GetProductRes, len(products.Products))
	for i, product := range products.Products {
		averageRating := float32(0)
		totalSales := int64(0)
		if product.SalesData != nil {
			averageRating = product.SalesData.Average_rating
			totalSales = product.SalesData.Total_sales
		}

		productRes[i] = dto.GetProductRes{
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			Stok:           product.Stock,
			Product_id:     product.Product_id,
			Category:       product.Category,
			Created_at:     product.Created_at,
			Images:         product.Images,
			Store_Name:     store.Name,
			City:           store.Address_Details.City,
			Average_Rating: averageRating,
			Total_Sales:    totalSales,
		}
	}

	return &dto.PagedProducts{
		Products:  productRes,
		Page:      page,
		TotalItem: products.TotalItem,
		LastPage:  products.LastPage,
	}, nil
}

// UpdateProduct implements domain.ProductService.
func (s *productService) UpdateProduct(ctx context.Context, storeID, email, productID string, req *dto.ProductReq) (*dto.EditProductRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	product, err := s.repo.GetProductById(ctx, productID, storeID)
	if err != nil || product == nil {
		return nil, errors.New("product not found" + err.Error())
	}

	updateAT := time.Now()
	var update primitive.D
	if product.Name != "" {
		update = append(update, bson.E{Key: "name", Value: product.Name})
	}
	if product.Description != "" {
		update = append(update, bson.E{Key: "description", Value: product.Description})
	}
	if len(product.Images) != 0 {
		update = append(update, bson.E{Key: "images", Value: product.Images})
	}
	if product.Price != 0 {
		update = append(update, bson.E{Key: "price", Value: product.Price})
	}
	if product.Category != "" {
		update = append(update, bson.E{Key: "category", Value: product.Category})
	}
	if product.Stock != 0 {
		update = append(update, bson.E{Key: "stock", Value: product.Stock})
	}

	update = append(update, bson.E{Key: "updated_at", Value: updateAT})

	result, err := s.repo.UpdateProduct(ctx, storeID, productID, update)
	if err != nil {
		return nil, errors.New("Failed to update the product: " + err.Error())
	}

	return &dto.EditProductRes{
		UpdateResult: result,
	}, nil
}

// GetAllProductSorted implements domain.ProductService.
func (s *productService) GetAllProductSorted(ctx context.Context, sortParams map[string]string, page int, email, storeID string) (*dto.PagedProducts, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	products, err := s.repo.GetAllProductSorted(ctx, sortParams, page, storeID)
	if err != nil {
		return nil, errors.New("failed to get all sorted products: " + err.Error())
	}

	productRes := make([]dto.GetProductRes, len(products.Products))
	for i, product := range products.Products {
		averageRating := float32(0)
		totalSales := int64(0)
		if product.SalesData != nil {
			averageRating = product.SalesData.Average_rating
			totalSales = product.SalesData.Total_sales
		}

		productRes[i] = dto.GetProductRes{
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			Stok:           product.Stock,
			Product_id:     product.Product_id,
			Category:       product.Category,
			Created_at:     product.Created_at,
			Images:         product.Images,
			Store_Name:     store.Name,
			City:           store.Address_Details.City,
			Average_Rating: averageRating,
			Total_Sales:    totalSales,
		}
	}

	return &dto.PagedProducts{
		Products:  productRes,
		Page:      page,
		TotalItem: products.TotalItem,
		LastPage:  products.LastPage,
	}, nil
}

// user / guest

// GetAllProductForGuest implements domain.ProductService.
func (s *productService) GetAllProductForGuest(ctx context.Context, page int) (*dto.PagedProducts, error) {
	products, err := s.repo.GetAllProduct(ctx, page)
	if err != nil {
		return nil, errors.New("failed to get all products: " + err.Error())
	}

	productRes := make([]dto.GetProductRes, len(products.Products))
	for i, product := range products.Products {
		store, err := s.storeRepo.GetStore(ctx, product.Store_id)
		if err != nil {
			return nil, errors.New("failed to get store by id: " + err.Error())
		}

		var average_rating float32
		var total_sales int64
		salesReport, err := s.salesReportRepo.GetByStoreId(ctx, product.Store_id)
		if err != nil {
			return nil, errors.New("failed to get sales report by store id: " + err.Error())
		}

		for _, item := range salesReport.Products {
			if item.Product_Id == product.Product_id {
				average_rating = item.Average_Rating
				total_sales = item.Total_Sales
				break
			}
		}

		productRes[i] = dto.GetProductRes{
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			Stok:           product.Stock,
			Product_id:     product.Product_id,
			Category:       product.Category,
			Created_at:     product.Created_at,
			Images:         product.Images,
			Store_Name:     store.Name,
			City:           store.Address_Details.City,
			Average_Rating: average_rating,
			Total_Sales:    total_sales,
		}
	}

	return &dto.PagedProducts{
		Products:  productRes,
		Page:      page,
		TotalItem: products.TotalItem,
		LastPage:  products.LastPage,
	}, nil
}

// GetAllByCategory implements domain.ProductService.
func (s *productService) GetAllByCategoryForGuest(ctx context.Context, category string, page int) (*dto.PagedProducts, error) {
	isValidCategory := false
	for _, categoryItem := range validCategories {
		if category == categoryItem {
			isValidCategory = true
			break
		}
	}

	if !isValidCategory {
		return nil, errors.New("invalid category")
	}

	products, err := s.repo.GetAllByCategory(ctx, category, page)
	if err != nil {
		return nil, errors.New("failed to get all products: " + err.Error())
	}

	productRes := make([]dto.GetProductRes, len(products.Products))
	for i, product := range products.Products {
		store, err := s.storeRepo.GetStore(ctx, product.Store_id)
		if err != nil {
			return nil, errors.New("failed to get store by id: " + err.Error())
		}

		averageRating := float32(0)
		totalSales := int64(0)
		if product.SalesData != nil {
			averageRating = product.SalesData.Average_rating
			totalSales = product.SalesData.Total_sales
		}

		productRes[i] = dto.GetProductRes{
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			Stok:           product.Stock,
			Product_id:     product.Product_id,
			Category:       product.Category,
			Created_at:     product.Created_at,
			Images:         product.Images,
			Store_Name:     store.Name,
			City:           store.Address_Details.City,
			Average_Rating: averageRating,
			Total_Sales:    totalSales,
		}
	}

	return &dto.PagedProducts{
		Products:  productRes,
		Page:      page,
		TotalItem: products.TotalItem,
		LastPage:  products.LastPage,
	}, nil
}

// GetProductByIdForGuest implements domain.ProductService.
func (s *productService) GetProductByIdForGuest(ctx context.Context, productID string) (*dto.GetProductRes, error) {
	product, err := s.repo.GetProductById(ctx, productID)
	if err != nil {
		return nil, errors.New("failed to get product by id: " + err.Error())
	}

	store, err := s.storeRepo.GetStore(ctx, product.Store_id)
	if err != nil {
		return nil, errors.New("failed to get store by id: " + err.Error())
	}

	var average_rating float32
	var total_sales int64
	salesReport, err := s.salesReportRepo.GetByStoreId(ctx, product.Store_id)
	if err != nil {
		return nil, errors.New("failed to get sales report by store id: " + err.Error())
	}

	for _, item := range salesReport.Products {
		if item.Product_Id == product.Product_id {
			average_rating = item.Average_Rating
			total_sales = item.Total_Sales
			break
		}
	}

	productRes := dto.GetProductRes{
		Name:           product.Name,
		Description:    product.Description,
		Price:          product.Price,
		Stok:           product.Stock,
		Product_id:     product.Product_id,
		Category:       product.Category,
		Created_at:     product.Created_at,
		Images:         product.Images,
		Store_Name:     store.Name,
		City:           store.Address_Details.City,
		Average_Rating: average_rating,
		Total_Sales:    total_sales,
	}

	return &productRes, nil
}

// SearchProductForGuest implements domain.ProductService.
func (s *productService) SearchProductForGuest(ctx context.Context, page int, query ...string) (*dto.PagedProducts, error) {
	products, err := s.repo.GetAllProductByQueryForCust(ctx, page, query[0], query[1])
	if err != nil {
		return nil, errors.New("failed to get all products by the query: " + err.Error())
	}

	productRes := make([]dto.GetProductRes, len(products.Products))
	for i, product := range products.Products {
		store, err := s.storeRepo.GetStore(ctx, product.Store_id)
		if err != nil {
			return nil, errors.New("failed to get store by id: " + err.Error())
		}

		averageRating := float32(0)
		totalSales := int64(0)
		if product.SalesData != nil {
			averageRating = product.SalesData.Average_rating
			totalSales = product.SalesData.Total_sales
		}

		productRes[i] = dto.GetProductRes{
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			Stok:           product.Stock,
			Product_id:     product.Product_id,
			Category:       product.Category,
			Created_at:     product.Created_at,
			Images:         product.Images,
			Store_Name:     store.Name,
			City:           store.Address_Details.City,
			Average_Rating: averageRating,
			Total_Sales:    totalSales,
		}
	}

	return &dto.PagedProducts{
		Products:  productRes,
		Page:      page,
		TotalItem: products.TotalItem,
		LastPage:  products.LastPage,
	}, nil
}

// GetAllProductSorted implements domain.ProductService.
func (s *productService) GetAllProductSortedForCust(ctx context.Context, sortParams map[string]string, page int) (*dto.PagedProducts, error) {
	products, err := s.repo.GetAllProductSorted(ctx, sortParams, page)
	if err != nil {
		return nil, errors.New("failed to get all sorted products: " + err.Error())
	}

	productRes := make([]dto.GetProductRes, len(products.Products))
	for i, product := range products.Products {
		store, err := s.storeRepo.GetStore(ctx, product.Store_id)
		if err != nil {
			return nil, errors.New("failed to get store by id: " + err.Error())
		}

		averageRating := float32(0)
		totalSales := int64(0)
		if product.SalesData != nil {
			averageRating = product.SalesData.Average_rating
			totalSales = product.SalesData.Total_sales
		}

		productRes[i] = dto.GetProductRes{
			Name:           product.Name,
			Description:    product.Description,
			Price:          product.Price,
			Stok:           product.Stock,
			Product_id:     product.Product_id,
			Category:       product.Category,
			Created_at:     product.Created_at,
			Images:         product.Images,
			Store_Name:     store.Name,
			City:           store.Address_Details.City,
			Average_Rating: averageRating,
			Total_Sales:    totalSales,
		}
	}

	return &dto.PagedProducts{
		Products:  productRes,
		Page:      page,
		TotalItem: products.TotalItem,
		LastPage:  products.LastPage,
	}, nil
}
