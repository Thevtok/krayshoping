package entity

type MerchantProduct struct {
	Product_ID      int    `json:"product_id"`
	MerchantID      string `json:"merchant_id"`
	MerchantName    string `json:"merchant_name"`
	ProductName     string `json:"product_name"`
	ProductPrice    int    `json:"product_price"`
	ProductQuantity int    `json:"product_quantity"`
}
