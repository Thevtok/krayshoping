package controller

import (
	"krayshoping/model/entity"
	"krayshoping/model/response"
	"krayshoping/repository"
	"krayshoping/utils"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type MerchantController struct {
	merchantRepo repository.MerchantRepository
}

func NewMerchantController(merchantRepo repository.MerchantRepository) *MerchantController {
	controller := MerchantController{
		merchantRepo: merchantRepo,
	}
	return &controller
}

func (c *MerchantController) Login(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	var user entity.Merchant

	if err := ctx.ShouldBindJSON(&user); err != nil {
		logrus.Errorf("invalid json")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid json")
		return
	}

	// Retrieve the user by email and password
	foundUser, err := c.merchantRepo.LoginMerchant(user.Email, user.Password)
	if err != nil {
		logrus.Errorf("login failed: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid email or password")
		return
	}
	expirationTime := time.Now().Add(time.Hour * 1)
	// Generate a token
	token, err := utils.GenerateTokenMerchant(foundUser, expirationTime)
	if err != nil {
		logrus.Errorf("failed to generate token")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "failed to generate token")
		return
	}

	// Set the token as a secure HTTP-only cookie

	ctx.SetCookie("token", token, 3600, "/", "", true, true)

	// Return a success response
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, gin.H{"message": "Login successful"})

	logrus.Info("Successfully logged in")
}
func (c *MerchantController) Logout(ctx *gin.Context) {
	// Ambil token dari cookie
	_, err := ctx.Cookie("token")
	if err != nil {
		logrus.Errorf("Failed to get token from cookie: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get token")
		return
	}

	ctx.SetCookie("token", "", -1, "/", "", true, true)

	// Proses tambahan yang mungkin diperlukan, seperti menghapus token dari database, membersihkan sesi, dll.

	response.JSONSuccess(ctx.Writer, true, http.StatusOK, gin.H{"message": "Logout successful"})
}

func (c *MerchantController) Register(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	newUser := entity.Merchant{}
	if err := ctx.BindJSON(&newUser); err != nil {
		logrus.Errorf("Invalid Input : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input")
		return
	}

	err = c.merchantRepo.CreateMerchant(&newUser)
	if err != nil {
		logrus.Errorf("Failed to Register Merchant : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to Register Merchant")
		return
	}

	logrus.Info("Success Register Merchant")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "created Merchant successfully")
}
