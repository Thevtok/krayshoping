package entity

type Merchant struct {
	ID string `json:"merchant_id"`

	MerchantName string `json:"merchant_name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Phone_Number string `json:"phone_number"`
	Address      string `json:"address"`
	Balance      int    `json:"balance"`
}
