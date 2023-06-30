package repository

import (
	"database/sql"
	"fmt"
	"krayshoping/model/entity"
	"krayshoping/utils"
)

type MerchantRepository interface {
	LoginMerchant(email string, password string) (*entity.Merchant, error)

	GetByiDMerchant(id string) (*entity.Merchant, error)
	CreateMerchant(Merchant *entity.Merchant) error

	UpdateBalanceMerchant(userID string, newBalance int) error
}

type merchantRepository struct {
	db *sql.DB
}

func NewMerchantRepository(db *sql.DB) MerchantRepository {
	return &merchantRepository{db: db}
}

func (r *merchantRepository) UpdateBalanceMerchant(merchantID string, newBalance int) error {
	_, err := r.GetByiDMerchant(merchantID)
	if err != nil {
		return fmt.Errorf("failed to get merchant data: %v", err)
	}

	query := "UPDATE merchants SET balance = $1 WHERE merchant_id = $2"
	_, err = r.db.Exec(query, newBalance, merchantID)
	if err != nil {

		return fmt.Errorf("failed to update merchant balance: %v", err)
	}
	return nil
}

func (r *merchantRepository) GetByiDMerchant(id string) (*entity.Merchant, error) {
	var user entity.Merchant

	err := r.db.QueryRow("SELECT merchant_name, merchant_id, email, phone_number, address,balance FROM merchants WHERE merchant_id = $1", id).Scan(&user.MerchantName, &user.ID, &user.Email, &user.Phone_Number, &user.Address, &user.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant ID not found: %v", err)
		}
		return nil, fmt.Errorf("failed to get merchant data: %v", err)
	}
	return &user, nil
}

func (r *merchantRepository) CreateMerchant(user *entity.Merchant) error {

	hashedPassword, err := utils.HasingPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = hashedPassword

	_, err = r.db.Exec("INSERT INTO merchants (merchant_id,merchant_name, email, password, phone_number, address, balance) VALUES ($1, $2, $3, $4, $5, $6, $7)", newUUID, user.MerchantName, user.Email, user.Password, user.Phone_Number, user.Address, 0)

	if err != nil {
		return fmt.Errorf("failed to create merchant: %v", err)
	}
	return nil
}

func (r *merchantRepository) LoginMerchant(email string, password string) (*entity.Merchant, error) {
	var m entity.Merchant
	query := "SELECT merchant_id, merchant_name, password FROM merchants WHERE email = $1"
	row := r.db.QueryRow(query, email)
	var hashedPassword string
	err := row.Scan(&m.ID, &m.MerchantName, &hashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("merchant not found")
		}

		return nil, fmt.Errorf("failed to get merchant")
	}

	// Verify that the retrieved password is a valid hash
	err = utils.CheckPasswordHash(password, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials \n password = %s\n hashed = %s", password, hashedPassword)
	}

	user := &entity.Merchant{
		Password:     hashedPassword,
		MerchantName: m.MerchantName,
		ID:           m.ID,
	}

	// Save the device token for the user

	return user, nil
}
