package entity

type User struct {
	ID   string `json:"user_id"`
	Name string `json:"name"`

	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Balance     int    `json:"balance"`
}
