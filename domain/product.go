package domain

import (
	"context"
	"time"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Products struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `json:"name" valid:"required,min=2,max=200" bson:"name"`
	Description string             `json:"description" valid:"required" bson:"description"`
	Price       float64            `json:"price" valid:"required" bson:"price"`
	Stock       int                `json:"stock" valid:"required" bson:"stock"`
	Product_id  string             `json:"product_id" bson:"product_id"`
	Category    string             `json:"category" valid:"required" bson:"category"`
	Created_at  time.Time          `json:"created_at" bson:"created_at"`
	Updated_at  time.Time          `json:"updated_at" bson:"updated_at"`
	Store_id    string             `json:"store_id" bson:"store_id"`
	Images      []string           `json:"images" valid:"required" bson:"images"`
}

type SalesData struct {
	Average_rating float32 `bson:"average_rating"`
	Total_sales    int64   `bson:"total_sales"`
}

type ProductWithSalesData struct {
	ID          string     `bson:"_id"`
	Name        string     `bson:"name"`
	Description string     `bson:"description"`
	Price       float64    `bson:"price"`
	Stock       int        `bson:"stock"`
	Product_id  string     `bson:"product_id"`
	Category    string     `bson:"category"`
	Created_at  time.Time  `bson:"created_at"`
	Updated_at  time.Time  `bson:"updated_at"`
	Store_id    string     `bson:"store_id"`
	Images      []string   `bson:"images"`
	SalesData   *SalesData `bson:"sales_data"`
}

type PagedProducts struct {
	Products  []ProductWithSalesData `json:"products"`
	Page      int                    `json:"page"`
	TotalItem int                    `json:"total_item"`
	LastPage  int                    `json:"last_page"`
}

type ProductRepository interface {
	CreateProduct(ctx context.Context, product Products) (primitive.ObjectID, error)
	CheckNameExists(ctx context.Context, name string) (bool, error)
	UpdateProduct(ctx context.Context, storeID, productID string, update bson.D) (*mongo.UpdateResult, error)
	UpdateStockProduct(ctx context.Context, storeID, productID string, stock int, updateAt time.Time) (*mongo.UpdateResult, error)
	DeleteProductById(ctx context.Context, storeID, productID string) (*mongo.DeleteResult, error)
	GetAllByCategory(ctx context.Context, category string, page int, storeID ...string) (*PagedProducts, error)
	GetAllProductByQuery(ctx context.Context, query string, page int, storeID ...string) (*PagedProducts, error)
	GetProductById(ctx context.Context, productID string, storeID ...string) (*ProductWithSalesData, error)
	GetAllProduct(ctx context.Context, page int, storeID ...string) (*PagedProducts, error)
	GetAllProductWithNoPage(ctx context.Context, storeID string) (*[]ProductWithSalesData, error)
	GetAllProductByQueryForCust(ctx context.Context, page int, query ...string) (*PagedProducts, error)
	GetAllProductSorted(ctx context.Context, sortParams map[string]string, page int, storeID ...string) (*PagedProducts, error)
}

type ProductService interface {
	// seller
	CreateProduct(ctx context.Context, storeID, email string, req *dto.ProductReq) (*dto.AddProductRes, error)
	GetProductById(ctx context.Context, storeID, email, productID string) (*dto.GetProductRes, error)
	GetAllProduct(ctx context.Context, storeID, email string, page int) (*dto.PagedProducts, error)
	SearchProduct(ctx context.Context, email, storeID, query string, page int) (*dto.PagedProducts, error)
	GetAllByCategory(ctx context.Context, email, storeID, category string, page int) (*dto.PagedProducts, error)
	UpdateProduct(ctx context.Context, storeID, email, productID string, req *dto.ProductReq) (*dto.EditProductRes, error)
	DeleteProductById(ctx context.Context, storeID, email, productID string) (*dto.DeleteProductRes, error)
	GetAllProductSorted(ctx context.Context, sortParams map[string]string, page int, email, storeID string) (*dto.PagedProducts, error)

	// user / guest
	GetAllProductForGuest(ctx context.Context, page int) (*dto.PagedProducts, error)
	GetProductByIdForGuest(ctx context.Context, productID string) (*dto.GetProductRes, error)
	GetAllByCategoryForGuest(ctx context.Context, category string, page int) (*dto.PagedProducts, error)
	SearchProductForGuest(ctx context.Context, page int, query ...string) (*dto.PagedProducts, error)
	GetAllProductSortedForCust(ctx context.Context, sortParams map[string]string, page int) (*dto.PagedProducts, error)
}
