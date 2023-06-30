package repository

import (
	"database/sql"
	"fmt"
	"krayshoping/model/entity"
	"time"
)

var now = time.Now().Local()
var date = now.Format("2006-01-02")

type TransactionRepository interface {
	CreateDepositBank(tx *entity.Deposit) error
	CreatePayment(tx *entity.Payment) error

	CreateTransfer(tx *entity.Transfer) error

	GetTransactions(userID string) ([]*entity.Transaction, error)
	GetPaymentByTxID(txID int) (*entity.Payment, error)

	UpdateDepositStatus(vaNumber, token string) error
	UpdatePaymentStatus(txID int, status string) error
}

type transactionRepository struct {
	db *sql.DB
}

func NewTxRepository(db *sql.DB) TransactionRepository {
	repo := new(transactionRepository)
	repo.db = db
	return repo
}
func (r *transactionRepository) GetTransactions(userID string) ([]*entity.Transaction, error) {
	query := `
	SELECT 
		t.tx_id, t.tx_type, t.tx_date,
		d.bank_name, d.account_number, d.account_holder_name, d.deposit_amount, d.status,d.user_id,
		tr.sender_name, tr.sender_phone_number, tr.recipient_name, tr.recipient_phone_number, tr.transfer_amount, tr.status,tr.sender_id,tr.recipient_id,
		p.customer_name, p.merchant_name, p.product_name, p.product_quantity, p.payment_amount, p.status,p.user_id,p.merchant_id
	FROM transactions t
	LEFT JOIN tx_deposit d ON t.tx_id = d.tx_id
	LEFT JOIN tx_transfer tr ON t.tx_id = tr.tx_id
	LEFT JOIN tx_payment p ON t.tx_id = p.tx_id
	WHERE (t.sender_id = $1 OR t.recipient_id = $1)
	ORDER BY t.tx_id DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	transactions := make([]*entity.Transaction, 0)

	for rows.Next() {
		var (
			txID                   int
			transactionType        string
			transactionDate        string
			depositBankName        sql.NullString
			depositAccountNumber   sql.NullString
			depositAccountHolder   sql.NullString
			depositAmount          sql.NullInt64
			depositStatus          sql.NullString
			dUserID                sql.NullString
			transferSenderName     sql.NullString
			transferSenderPhone    sql.NullString
			transferRecipientName  sql.NullString
			transferRecipientPhone sql.NullString
			transferAmount         sql.NullFloat64
			transferStatus         sql.NullString
			tSenderID              sql.NullString
			tRecipientID           sql.NullString
			paymentCustomerName    sql.NullString
			paymentMerchantName    sql.NullString
			paymentProductName     sql.NullString
			paymentProductQuantity sql.NullInt64
			paymentAmount          sql.NullFloat64
			paymentStatus          sql.NullString
			pUserID                sql.NullString
			pMerchantID            sql.NullString
		)

		err := rows.Scan(
			&txID, &transactionType, &transactionDate,
			&depositBankName, &depositAccountNumber, &depositAccountHolder, &depositAmount, &depositStatus, &dUserID,
			&transferSenderName, &transferSenderPhone, &transferRecipientName, &transferRecipientPhone, &transferAmount, &transferStatus, &tSenderID, &tRecipientID,
			&paymentCustomerName, &paymentMerchantName, &paymentProductName, &paymentProductQuantity, &paymentAmount, &paymentStatus, &pUserID, &pMerchantID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction row: %v", err)
		}

		transaction := &entity.Transaction{
			TxID:     txID,
			TxType:   transactionType,
			TxDate:   transactionDate,
			Deposit:  make([]entity.Deposit, 0),
			Transfer: make([]entity.Transfer, 0),
			Payment:  make([]entity.Payment, 0),
		}

		if depositBankName.Valid && depositAccountNumber.Valid && depositAccountHolder.Valid && depositAmount.Valid && depositStatus.Valid && dUserID.Valid {
			deposit := &entity.Deposit{
				BankName:          depositBankName.String,
				AccountNumber:     depositAccountNumber.String,
				AccountHolderName: depositAccountHolder.String,
				DepositAmount:     int(depositAmount.Int64),
				Status:            depositStatus.String,
				UserID:            dUserID.String,
			}
			transaction.Deposit = append(transaction.Deposit, *deposit)
		}
		if transferSenderName.Valid && transferSenderPhone.Valid && transferRecipientName.Valid && transferRecipientPhone.Valid && transferAmount.Valid && transferStatus.Valid && tSenderID.Valid && tRecipientID.Valid {
			transfer := &entity.Transfer{
				SenderName:           transferSenderName.String,
				SenderPhoneNumber:    transferSenderPhone.String,
				RecipientName:        transferRecipientName.String,
				RecipientPhoneNumber: transferRecipientPhone.String,
				Transfer_Amount:      int(transferAmount.Float64),
				Status:               transferStatus.String,
				SenderID:             tSenderID.String,
				RecipientID:          tRecipientID.String,
			}
			transaction.Transfer = append(transaction.Transfer, *transfer)

		}
		if paymentCustomerName.Valid && paymentMerchantName.Valid && paymentProductName.Valid && paymentProductQuantity.Valid && paymentAmount.Valid && paymentStatus.Valid && pUserID.Valid && pMerchantID.Valid {
			payment := &entity.Payment{
				CustomerName:    paymentCustomerName.String,
				MerchantName:    paymentMerchantName.String,
				ProductName:     paymentProductName.String,
				ProductQuantity: int(paymentProductQuantity.Int64),
				PaymentAmount:   int(paymentAmount.Float64),
				Status:          paymentStatus.String,
				UserID:          pUserID.String,
				MerchantID:      pMerchantID.String,
			}
			transaction.Payment = append(transaction.Payment, *payment)
		}

		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate through result set: %v", err)
	}
	return transactions, nil
}
func (r *transactionRepository) CreateDepositBank(tx *entity.Deposit) error {
	query := "INSERT INTO transactions (tx_type, tx_date, sender_id) VALUES ($1, $2, $3)"
	_, err := r.db.Exec(query, "Deposit", date, tx.UserID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}

	var txID int
	err = r.db.QueryRow("SELECT lastval()").Scan(&txID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction ID: %v", err)
	}

	query = "INSERT INTO tx_deposit (tx_id, deposit_amount, bank_name, account_number, account_holder_name,status,va_number,token,user_id) VALUES ($1, $2, $3, $4, $5,$6,$7,$8,$9)"
	_, err = r.db.Exec(query, txID, tx.DepositAmount, tx.BankName, tx.AccountNumber, tx.AccountHolderName, "Pending", tx.VaNumber, tx.Token, tx.UserID)
	if err != nil {
		return fmt.Errorf("failed to insert deposit: %v", err)
	}

	return nil
}

func (r *transactionRepository) CreateTransfer(tx *entity.Transfer) error {
	query := "INSERT INTO transactions (tx_type, tx_date, sender_id,recipient_id) VALUES ($1, $2, $3,$4)"
	_, err := r.db.Exec(query, "Transfer", date, tx.SenderID, tx.RecipientID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}

	var txID int
	err = r.db.QueryRow("SELECT lastval()").Scan(&txID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction ID: %v", err)
	}

	query = "INSERT INTO tx_transfer (tx_id, sender_name, recipient_name, transfer_amount, sender_phone_number, recipient_phone_number,sender_id,recipient_id,status) VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9)"
	_, err = r.db.Exec(query, txID, tx.SenderName, tx.RecipientName, tx.Transfer_Amount, tx.SenderPhoneNumber, tx.RecipientPhoneNumber, tx.SenderID, tx.RecipientID, "Success")
	if err != nil {
		return fmt.Errorf("failed to insert transfer: %v", err)
	}

	return nil
}

func (r *transactionRepository) CreatePayment(tx *entity.Payment) error {
	query := "INSERT INTO transactions (tx_type, tx_date, sender_id,merchant_id) VALUES ($1, $2, $3,$4)"
	_, err := r.db.Exec(query, "Payment", date, tx.UserID, tx.MerchantID)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %v", err)
	}

	var txID int
	err = r.db.QueryRow("SELECT lastval()").Scan(&txID)
	if err != nil {
		return fmt.Errorf("failed to retrieve transaction ID: %v", err)
	}

	query = "INSERT INTO tx_payment (tx_id, payment_amount, customer_name, merchant_name, product_name,status,product_quantity,user_id,merchant_id) VALUES ($1, $2, $3, $4, $5,$6,$7,$8,$9)"
	_, err = r.db.Exec(query, txID, tx.PaymentAmount, tx.CustomerName, tx.MerchantName, tx.ProductName, "Pending", tx.ProductQuantity, tx.UserID, tx.MerchantID)
	if err != nil {
		return fmt.Errorf("failed to insert deposit: %v", err)
	}

	return nil
}

func (r *transactionRepository) UpdateDepositStatus(vaNumber, token string) error {
	query := "UPDATE tx_deposit SET status = $1, va_number = $2 WHERE token = $3"
	_, err := r.db.Exec(query, "Success", vaNumber, token)
	if err != nil {
		return fmt.Errorf("failed to update deposit status: %v", err)
	}

	return nil
}

func (r *transactionRepository) UpdatePaymentStatus(txID int, status string) error {
	var dbStatus string
	if status == "success" {
		dbStatus = "success"
	} else if status == "failed" {
		dbStatus = "failed"
	} else {
		return fmt.Errorf("invalid status value: %s", status)
	}

	query := "UPDATE tx_payment SET status = $1 WHERE tx_id = $2"
	_, err := r.db.Exec(query, dbStatus, txID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}
	return nil
}

func (r *transactionRepository) GetPaymentByTxID(txID int) (*entity.Payment, error) {
	var payment entity.Payment
	err := r.db.QueryRow("SELECT tx_id,merchant_id,user_id,payment_amount, product_quantity,merchant_name,product_name,customer_name FROM tx_payment WHERE tx_id = $1", txID).Scan(&txID, &payment.MerchantID, &payment.UserID, &payment.PaymentAmount, &payment.ProductQuantity, &payment.MerchantName, &payment.ProductName, &payment.CustomerName)
	if err != nil {

		return nil, fmt.Errorf("failed to get payment: %v", err)
	}
	return &payment, nil
}
