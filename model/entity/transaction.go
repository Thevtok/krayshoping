package entity

type Transaction struct {
	TxID     int        `json:"tx_id"`
	TxType   string     `json:"tx_type"`
	TxDate   string     `json:"tx_date"`
	Deposit  []Deposit  `json:"deposit,omitempty"`
	Transfer []Transfer `json:"transfer,omitempty"`
	Payment  []Payment  `json:"payment,omitempty"`
}

type Deposit struct {
	UserID            string `json:"user_id"`
	DepositAmount     int    `json:"deposit_amount"`
	BankName          string `json:"bank_name"`
	AccountNumber     string `json:"account_number"`
	AccountHolderName string `json:"account_holder_name"`

	Status   string `json:"status"`
	VaNumber string `json:"va_number"`
	Token    string `json:"token"`
}

type Transfer struct {
	SenderID             string `json:"sender_id"`
	RecipientID          string `json:"recipient_id"`
	Transfer_Amount      int    `json:"transfer_amount"`
	SenderPhoneNumber    string `json:"sender_phone_number"`
	RecipientPhoneNumber string `json:"recipient_phone_number"`
	SenderName           string `json:"sender_name"`
	RecipientName        string `json:"recipient_name"`

	Status string `json:"status"`
}

type Payment struct {
	UserID string `json:"user_id"`

	PaymentAmount int `json:"payment_amount"`

	MerchantID      string `json:"merchant_id"`
	MerchantName    string `json:"merchant_name"`
	ProductName     string `json:"product_name"`
	CustomerName    string `json:"customer_name"`
	ProductQuantity int    `json:"product_quantity"`

	Status string `json:"status"`
}
