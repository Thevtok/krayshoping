package repository

import (
	"database/sql"
	"fmt"
	"krayshoping/model/entity"
	"krayshoping/utils"

	"github.com/google/uuid"
)

var newUUID = uuid.New()

type UserRepository interface {
	Login(email string, password string) (*entity.User, error)

	GetByID(id string) (*entity.User, error)
	Create(user *entity.User) error

	UpdateBalance(userID string, newBalance int) error

	GetByPhone(phoneNumber string) (*entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) UpdateBalance(userID string, newBalance int) error {
	_, err := r.GetByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	query := "UPDATE users SET balance = $1 WHERE user_id = $2"
	_, err = r.db.Exec(query, newBalance, userID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %v", err)
	}

	return nil
}

func (r *userRepository) GetByPhone(phoneNumber string) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow("SELECT user_id,name,email, phone_number, address  FROM users WHERE phone_number = $1", phoneNumber).Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Address)
	if err != nil {

		return nil, fmt.Errorf("failed to get user by phone number: %v", err)
	}
	return &user, nil
}

func (r *userRepository) GetByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow(`
	SELECT name, user_id, email, phone_number, address, balance
	FROM users
	WHERE user_id = $1
`, id).Scan(&user.Name, &user.ID, &user.Email, &user.PhoneNumber, &user.Address, &user.Balance)

	if err != nil {

		return nil, fmt.Errorf("failed to get user by ID: %v", err)
	}

	return &user, nil
}

func (r *userRepository) Create(user *entity.User) error {
	hashedPassword, err := utils.HasingPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	user.Password = hashedPassword

	_, err = r.db.Exec("INSERT INTO users (user_id, name,email, password, phone_number, address, balance, role) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", newUUID, user.Name, user.Email, user.Password, user.PhoneNumber, user.Address, 0, "user")
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	return nil
}

func (r *userRepository) Login(email string, password string) (*entity.User, error) {
	var u entity.User
	query := "SELECT user_id, password FROM users WHERE email = $1"
	row := r.db.QueryRow(query, email)
	var hashedPassword string
	err := row.Scan(&u.ID, &hashedPassword)
	if err != nil {

		return nil, fmt.Errorf("failed to get user")
	}
	// Verify the password
	err = utils.CheckPasswordHash(password, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials \n password = %s\n hashed = %s", password, hashedPassword)
	}

	user := &entity.User{
		Password: hashedPassword,

		ID: u.ID,

		Email: u.Email,
	}

	// Save the device token for the user

	return user, nil
}
