package repository

import (
	"database/sql"
	"fmt"
	"krayshoping/model/entity"
)

type BankRepository interface {
	GetByUserID(id string) ([]*entity.Bank, error)
	GetByAccountNumber(number string) (*entity.Bank, error)
	Create(id string, newBankAcc *entity.Bank) (any, error)

	Delete(number string) error
}

type bankRepository struct {
	db *sql.DB
}

func (r *bankRepository) GetByUserID(id string) ([]*entity.Bank, error) {
	var bankAccs []*entity.Bank
	query := "SELECT user_id, account_id, bank_name, account_number, account_holder_name FROM banks WHERE user_id = $1"
	rows, err := r.db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var bankAcc entity.Bank
		err = rows.Scan(&bankAcc.UserID, &bankAcc.AccountID, &bankAcc.BankName, &bankAcc.AccountNumber, &bankAcc.AccountHolderName)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		bankAccs = append(bankAccs, &bankAcc)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate rows: %v", err)
	}

	return bankAccs, nil
}

func (r *bankRepository) GetByAccountNumber(number string) (*entity.Bank, error) {
	var bankAcc entity.Bank
	query := "SELECT account_id, bank_name, account_number, account_holder_name, user_id FROM banks WHERE account_number = $1"
	row := r.db.QueryRow(query, number)
	err := row.Scan(&bankAcc.AccountID, &bankAcc.BankName, &bankAcc.AccountNumber, &bankAcc.AccountHolderName, &bankAcc.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account number not found")
		}
		return nil, fmt.Errorf("failed to get data: %v", err)
	}

	return &bankAcc, nil
}

func (r *bankRepository) Create(id string, newBankAcc *entity.Bank) (any, error) {
	query := "INSERT INTO banks(user_id, bank_name, account_number, account_holder_name) VALUES ($1, $2, $3, $4)"
	_, err := r.db.Exec(query, id, newBankAcc.BankName, newBankAcc.AccountNumber, newBankAcc.AccountHolderName)
	if err != nil {
		return nil, fmt.Errorf("failed to create data: %v", err)
	}

	return newBankAcc, nil
}

func (r *bankRepository) Delete(number string) error {
	_, err := r.GetByAccountNumber(number)
	if err != nil {

		return fmt.Errorf("failed to get data: %v", err)
	}

	query := "DELETE FROM banks WHERE account_number = $1"
	_, err = r.db.Exec(query, number)
	if err != nil {

		return fmt.Errorf("failed to delete data: %v", err)
	}
	return nil
}

func NewBankAccRepository(db *sql.DB) BankRepository {
	repo := new(bankRepository)
	repo.db = db
	return repo
}
