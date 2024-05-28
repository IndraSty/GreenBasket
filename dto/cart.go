package dto

import "time"

type GetCartItemRes struct {
	Product_Id    string    `json:"product_id" bson:"product_id"`
	Product_Name  string    `json:"product_name" bson:"product_name"`
	Product_Image []string  `json:"product_image" bson:"product_image"`
	Store_Name    string    `json:"store_name" bson:"store_name"`
	Quantity      int       `json:"quantity" bson:"quantity"`
	AddedAt       time.Time `json:"added_at" bson:"added_at"`
	Selected      bool      `json:"selected" bson:"selected"`
	Price         float64   `json:"price" bson:"price"`
}

type AddCartReq struct {
	Quantity int `json:"quantity"`
}

type CartItemEditReq struct {
	Quantity int  `json:"quantity"`
	Selected bool `json:"selected"`
}

type CartItemEditRepo struct {
	Quantity    int     `json:"quantity"`
	Selected    bool    `json:"selected"`
	Total_Price float64 `json:"total_price"`
}
