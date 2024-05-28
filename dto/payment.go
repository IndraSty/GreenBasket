package dto

type PaymentRes struct {
	Snap_Url string `json:"snap_url"`
}

type PaymentReq struct {
	OrderID string  `json:"-"`
	UserID  string  `json:"-"`
	Amount  float64 `json:"amount"`
}

type UpdatePaymentReq struct {
	Payment_Method string `json:"payment_method" bson:"payment_method"`
	Status         string `json:"status" bson:"status"`
	TransactionID  string `json:"transaction_id" bson:"transaction_id"`
}
