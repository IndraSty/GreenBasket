package dto

type SalesReportRes struct {
	Store_Id      string            `json:"store_id" bson:"store_id"`
	Email         string            `json:"email" bson:"email"`
	Total_Sales   int64             `json:"total_sales" bson:"total_sales"`
	Total_Incomes float64           `json:"total_income" bson:"total_income"`
	Products      []ProductSalesRes `json:"products" bson:"products"`
}

type ProductSalesRes struct {
	Product_Id  string `json:"product_id" bson:"product_id"`
	Total_Sales int64  `json:"total_sales" bson:"total_sales"`
	Stock       int32  `json:"stock" bson:"stock"`
}
