package controller

import (
	"krayshoping/model/entity"
	"krayshoping/model/response"
	"krayshoping/repository"
	"krayshoping/utils"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BankController struct {
	bankRepo repository.BankRepository
}

func NewBankController(bankRepo repository.BankRepository) *BankController {
	controller := BankController{
		bankRepo: bankRepo,
	}
	return &controller
}

func (c *BankController) FindBankAccByUserID(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userID := ctx.Param("user_id")

	existingUser, err := c.bankRepo.GetByUserID(userID)
	if err != nil {
		logrus.Errorf("Bank Account not found: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank Account not found")
		return
	}

	logrus.Info("Bank Account loaded Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, existingUser)
}

func (c *BankController) FindBankAccByAccountNumber(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	accountNumber := ctx.Param("account_number")

	existingUser, err := c.bankRepo.GetByAccountNumber(accountNumber)
	if err != nil {
		logrus.Errorf("Bank Account not found: %v", err)

		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank not found")
		return
	}
	logrus.Info("Bank Account loaded Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, existingUser)
}

func (c *BankController) CreateBankAccount(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	userID := ctx.Param("user_id")

	var newBankAcc entity.Bank
	err = ctx.BindJSON(&newBankAcc)
	if err != nil {
		logrus.Errorf("Invalid request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid request body")
		return
	}
	newBankAcc.BankName = strings.ToLower(newBankAcc.BankName)

	if newBankAcc.BankName == "" || newBankAcc.AccountNumber == "" || newBankAcc.AccountHolderName == "" {
		logrus.Errorf("Invalid Input: Required fields are empty")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input: Required fields are empty")
		return
	}

	_, err = c.bankRepo.GetByAccountNumber(newBankAcc.AccountNumber)
	if err != nil {
		if err.Error() == "account number not found" {
			// Nomor akun belum ada dalam database, maka dapat dibuat akun baru
			result, err := c.bankRepo.Create(userID, &newBankAcc)
			if err != nil {
				logrus.Errorf("Failed to create Bank Account: %v", err)
				response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Bank Account")
				return
			}
			logrus.Info("Bank Account created Successfully")
			response.JSONSuccess(ctx.Writer, true, http.StatusCreated, result)
			return
		}

		// Terjadi error lain saat pengecekan nomor akun
		logrus.Errorf("Failed to check existing bank account: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to create Bank Account")
		return
	}

	// Nomor akun sudah ada dalam database
	logrus.Errorf("Bank account with the given account number already exists")
	response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Bank account with the given account number already exists")

}
func (c *BankController) Unreg(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	accountNumber := ctx.Param("account_number")

	// Pengecekan apakah nomor akun ada dalam database
	_, err = c.bankRepo.GetByAccountNumber(accountNumber)
	if err != nil {
		logrus.Errorf("Failed to retrieve Bank Account: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "Bank Account not found")
		return
	}

	// Hapus akun bank
	err = c.bankRepo.Delete(accountNumber)
	if err != nil {
		logrus.Errorf("Failed to delete Bank Account: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to delete Bank Account")
		return
	}

	logrus.Info("Bank Account deleted Successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, "Bank account deleted successfully")
}
