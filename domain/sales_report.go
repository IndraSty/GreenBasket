package domain

import (
	"context"

	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Sales_Report struct {
	ID            primitive.ObjectID `bson:"_id"`
	Store_Id      string             `json:"store_id" bson:"store_id"`
	Email         string             `json:"email" bson:"email"`
	Total_Sales   int64              `json:"total_sales" bson:"total_sales"`
	Total_Incomes float64            `json:"total_income" bson:"total_income"`
	Products      []Product_Sales    `json:"products" bson:"products"`
}

type Product_Sales struct {
	Product_Id     string  `json:"product_id" bson:"product_id"`
	Total_Sales    int64   `json:"total_sales" bson:"total_sales"`
	Stock          int32   `json:"stock" bson:"stock"`
	Average_Rating float32 `json:"average_rating" bson:"average_rating"`
}

type SalesReportRepository interface {
	Insert(ctx context.Context, input Sales_Report) (primitive.ObjectID, error)
	Update(ctx context.Context, storeID string, update bson.D) (*mongo.UpdateResult, error)
	UpdateAverageRating(ctx context.Context, storeID, rpdocutID string, averageRating float32) (*mongo.UpdateResult, error)
	GetByEmailAndStoreId(ctx context.Context, email, storeID string) (*Sales_Report, error)
	GetByStoreId(ctx context.Context, storeID string) (*Sales_Report, error)
}

type SalesReportService interface {
	UpdateSalesReport(ctx context.Context, storeID, email string) error
	GetSalesReport(ctx context.Context, email, storeID string) (*dto.SalesReportRes, error)
}
