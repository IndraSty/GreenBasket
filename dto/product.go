package dto

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GetProductRes struct {
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	Stok           int       `json:"stock"`
	Average_Rating float32   `json:"average_rating"`
	Total_Sales    int64     `json:"total_sales"`
	Product_id     string    `json:"product_id"`
	Category       string    `json:"category"`
	Created_at     time.Time `json:"created_at"`
	Store_Name     string    `json:"store_name"`
	City           string    `json:"city"`
	Images         []string  `json:"images"`
}

type PagedProducts struct {
	Products  []GetProductRes `json:"products"`
	Page      int             `json:"page"`
	TotalItem int             `json:"total_item"`
	LastPage  int             `json:"last_page"`
}

type ProductReq struct {
	Name        string   `json:"name" valid:"required,min=2,max=200" bson:"name"`
	Description string   `json:"description" valid:"required" bson:"description"`
	Price       float64  `json:"price" valid:"required" bson:"price"`
	Stok        int      `json:"stock" valid:"required" bson:"stock"`
	Category    string   `json:"category" valid:"required" bson:"category"`
	Images      []string `json:"images" valid:"required" bson:"images"`
}

type AddProductRes struct {
	InsertId *primitive.ObjectID
}

type EditProductRes struct {
	UpdateResult *mongo.UpdateResult
}

type DeleteProductRes struct {
	DeleteResult *mongo.DeleteResult
}
