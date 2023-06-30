package controller

import (
	"encoding/json"
	"io"
	"krayshoping/controller/midtranss"
	"krayshoping/model/entity"
	"krayshoping/model/response"
	"krayshoping/repository"
	"krayshoping/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransactionController struct {
	txRepository       repository.TransactionRepository
	userRepository     repository.UserRepository
	bankRepository     repository.BankRepository
	productRepository  repository.ProductRepository
	merchantRepository repository.MerchantRepository
}

func NewTransactionController(trp repository.TransactionRepository, urp repository.UserRepository, brp repository.BankRepository, prp repository.ProductRepository, mrp repository.MerchantRepository) *TransactionController {
	controller := TransactionController{
		txRepository:       trp,
		userRepository:     urp,
		bankRepository:     brp,
		productRepository:  prp,
		merchantRepository: mrp,
	}
	return &controller
}

var depositToken string
var userIDdepo string
var depoAmount int

func (c *TransactionController) HandlePaymentNotification(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	if ctx.Request.Method != http.MethodPost {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}
	defer ctx.Request.Body.Close()

	logrus.Println("Received payment notification:")
	logrus.Println(string(body))

	// Mengambil nomor virtual account (VA)
	var notification response.PaymentNotification
	err = json.Unmarshal(body, &notification)
	if err != nil {
		logrus.Errorf("Failed to decode notification payload: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Failed to decode notification payload"})
		return
	}

	go func() {
		user, err := c.userRepository.GetByID(userIDdepo)
		if err != nil {
			logrus.Errorf("Failed to get user: %v", err)
			return
		}
		newBalance := user.Balance + depoAmount

		if notification.TransactionStatus == "settlement" && len(notification.VANumbers) > 0 {
			vaNumber := notification.VANumbers[0].VANumber

			err = c.userRepository.UpdateBalance(userIDdepo, newBalance)
			if err != nil {
				logrus.Errorf("Failed to update balance user: %v", err)
				return
			}

			logrus.Infof("useridDEPO: %s", userIDdepo)
			logrus.Infof("depoAmount: %d", depoAmount)

			err = c.txRepository.UpdateDepositStatus(vaNumber, depositToken)
			if err != nil {
				logrus.Errorf("Failed to update deposit status: %v", err)
				return
			}

			logrus.Infof("Successfully updated balance and deposit status")
		} else {
			logrus.Infof("Notification ignored")
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{"message": "Notification received"})
}

func (c *TransactionController) GetTxByID(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userId := ctx.Param("user_id")

	// Get sender by ID
	_, err = c.userRepository.GetByID(userId)
	if err != nil {
		logrus.Errorf("Failed to get User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get User")
		return
	}

	txs, err := c.txRepository.GetTransactions(userId)
	if err != nil {
		logrus.Errorf("Failed to get Transaction %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get Transaction")
		return
	}
	if len(txs) == 0 {
		logrus.Errorf("Transaction not found")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Transaction not found")
		return
	}

	logrus.Info("Transaction Log loaded Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, txs)
}

func (c *TransactionController) CreateDepositBank(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse user_id parameter
	userID := ctx.Param("user_id")

	// Parse bank_account_id parameter
	accNumber := ctx.Param("account_number")

	user, err := c.userRepository.GetByID(userID)
	if err != nil {
		logrus.Errorf("user not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "user not found")
		return
	}
	// Retrieve bank account by bank_account_id
	bankAcc, err := c.bankRepository.GetByAccountNumber(accNumber)
	if err != nil {
		logrus.Errorf("Bank_account_id not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank_account_id not found")
		return
	}

	// Parse request body
	var reqBody entity.Deposit
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		logrus.Errorf("Incorrect request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Incorrect request body")
		return
	}
	reqBody.UserID = userID
	reqBody.BankName = bankAcc.BankName
	reqBody.AccountHolderName = bankAcc.AccountHolderName
	reqBody.AccountNumber = bankAcc.AccountNumber

	if reqBody.DepositAmount < 10000 {
		logrus.Errorf("Minimum deposit 10.000: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Minimum deposit 10.000")
		return
	}

	token, err := midtranss.CreateMidtransTransactionFromDeposit(&reqBody, user)
	if err != nil {
		logrus.Errorf("Failed to create Midtrans transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Midtrans transaction")
		return
	}

	userIDdepo = userID
	reqBody.Token = token
	depositToken = token
	depoAmount = reqBody.DepositAmount

	// Create the deposit transaction
	if err := c.txRepository.CreateDepositBank(&reqBody); err != nil {
		logrus.Errorf("Failed to create Deposit Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Deposit Transaction")
		return
	}

	// Kirim respons sukses
	logrus.Info("Deposit Transaction created Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, gin.H{
		"body_midtrans": token,
	})

}

func (c *TransactionController) CreateTransferTransaction(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse transfer data from request body
	var newTransfer entity.Transfer
	if err := ctx.BindJSON(&newTransfer); err != nil {
		logrus.Errorf("Failed to parse transfer data: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to parse transfer data: invalid JSON format")
		return
	}

	userID := ctx.Param("user_id")

	// Get sender by ID
	sender, err := c.userRepository.GetByID(userID)
	if err != nil {
		logrus.Errorf("Failed to get Sender User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get Sender User")
		return
	}

	recipient, err := c.userRepository.GetByPhone(newTransfer.RecipientPhoneNumber)
	if err != nil {
		logrus.Errorf("Failed to get Recipient User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get Recipient User")
		return
	}

	if sender.Balance < newTransfer.Transfer_Amount {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Insufficient balance")
		return
	}
	if sender.PhoneNumber == recipient.PhoneNumber {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusForbidden, "Input the recipient correctly")
		return
	}

	if newTransfer.Transfer_Amount < 10000 {
		logrus.Errorf("Minimum transfer amount is 10,000")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Minimum transfer amount is 10,000")
		return
	}
	newBalanceS := sender.Balance - newTransfer.Transfer_Amount - 2500
	err = c.userRepository.UpdateBalance(sender.ID, newBalanceS)
	if err != nil {
		logrus.Errorf("Failed to update balance sender %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to update balance sender")
		return
	}
	newBalanceR := recipient.Balance + newTransfer.Transfer_Amount
	err = c.userRepository.UpdateBalance(recipient.ID, newBalanceR)
	if err != nil {
		logrus.Errorf("Failed to update balance recipient %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to update balance recipient")
		return

	}

	newTransfer.SenderName = sender.Name
	newTransfer.RecipientName = recipient.Name
	newTransfer.SenderPhoneNumber = sender.PhoneNumber
	newTransfer.RecipientPhoneNumber = recipient.PhoneNumber
	newTransfer.SenderID = sender.ID
	newTransfer.RecipientID = recipient.ID

	// Create transfer transaction in use case layer
	err = c.txRepository.CreateTransfer(&newTransfer)
	logrus.Info("Processing transfer transaction...")
	logrus.Infof("Sender: %s, Recipient: %s, Amount: %d", sender.Name, recipient.Name, newTransfer.Transfer_Amount)

	if err != nil {
		logrus.Errorf("Failed to create Transfer Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Transfer Transaction")
		return
	}

	logrus.Info("Transfer Transaction created Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "Transfer Successfully")

}

func (c *TransactionController) CreatePaymentTransaction(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Parse transfer data from request body
	var newPayment entity.Payment
	if err := ctx.BindJSON(&newPayment); err != nil {
		logrus.Errorf("Failed to parse transfer data: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to parse transfer data: invalid JSON format")
		return
	}

	userID := ctx.Param("user_id")
	productIDStr := ctx.Param("product_id")
	productID, _ := strconv.Atoi(productIDStr)
	merchantID := ctx.Param("merchant_id")

	// Get sender by ID
	sender, err := c.userRepository.GetByID(userID)
	if err != nil {
		logrus.Errorf("Failed to get Sender User: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get Sender User")
		return
	}
	product, err := c.productRepository.GetByProductID(productID)
	if err != nil {
		logrus.Errorf("Failed to get product: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Failed to get product")
		return
	}

	if sender.Balance < product.ProductPrice*newPayment.ProductQuantity {
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Insufficient balance")
		return
	}
	totalPaymentAmount := product.ProductPrice * newPayment.ProductQuantity
	newPayment.CustomerName = sender.Name
	newPayment.MerchantID = merchantID
	newPayment.MerchantName = product.MerchantName
	newPayment.ProductName = product.ProductName
	newPayment.PaymentAmount = totalPaymentAmount
	newPayment.UserID = userID

	// Create transfer transaction in use case layer
	err = c.txRepository.CreatePayment(&newPayment)

	if err != nil {
		logrus.Errorf("Failed to create Payment Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Payment Transaction")
		return
	}
	sender.Balance = sender.Balance - totalPaymentAmount
	err = c.userRepository.UpdateBalance(userID, sender.Balance)
	if err != nil {
		logrus.Errorf("Failed to update user balance: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to update user balance")
		return
	}

	logrus.Info("Payment Transaction created Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "Payment Successfully")

}

func (c *TransactionController) UpdatePaymentStatus(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	txIDStr := ctx.Param("tx_id")
	txID, _ := strconv.Atoi(txIDStr)
	merchantID := ctx.Param("merchant_id")
	userID := ctx.Param("user_id")

	var newPayment entity.Payment
	if err := ctx.BindJSON(&newPayment); err != nil {
		logrus.Errorf("Failed to parse status data: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Failed to parse status data: invalid JSON format")
		return
	}
	payment, err := c.txRepository.GetPaymentByTxID(txID)
	if err != nil {
		logrus.Errorf("Failed to update get Payment Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get Payment Transaction")
		return
	}

	err = c.txRepository.UpdatePaymentStatus(txID, newPayment.Status)
	if err != nil {
		logrus.Errorf("Failed to update status Payment Transaction: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to update status Payment Transaction")
		return
	}

	if newPayment.Status == "success" {

		merchant, err := c.merchantRepository.GetByiDMerchant(merchantID)
		if err != nil {
			logrus.Errorf("Failed to get merchant: %v", err)
			response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get merchant")
			return
		}

		merchant.Balance += payment.PaymentAmount
		err = c.merchantRepository.UpdateBalanceMerchant(merchantID, merchant.Balance)
		if err != nil {
			logrus.Errorf("Failed to update merchant balance: %v", err)
			response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to update merchant balance")
			return
		}

		logrus.Info("Payment Transaction update status Successfully")
		response.JSONSuccess(ctx.Writer, true, http.StatusOK, "Update status Payment Successfully")
	} else if newPayment.Status == "failed" {

		user, err := c.userRepository.GetByID(userID)
		if err != nil {
			logrus.Errorf("Failed to get user: %v", err)
			response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get user")
			return
		}
		user.Balance += payment.PaymentAmount
		err = c.userRepository.UpdateBalance(userID, user.Balance)
		if err != nil {
			logrus.Errorf("Failed to update user balance: %v", err)
			response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to update user balance")
			return
		}
		logrus.Info("Payment Transaction update status Success")
		response.JSONSuccess(ctx.Writer, true, http.StatusOK, "Update status Payment Success")
	} else {
		logrus.Info("Invalid status value")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid status value")
	}
}
