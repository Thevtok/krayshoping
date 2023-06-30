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

type UserController struct {
	userRepo repository.UserRepository
}

func NewUserController(userRepo repository.UserRepository) *UserController {
	controller := UserController{
		userRepo: userRepo,
	}
	return &controller
}

func (c *UserController) Login(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	var user entity.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		logrus.Errorf("invalid json")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid json")
		return
	}

	// Retrieve the user by email and password
	foundUser, err := c.userRepo.Login(user.Email, user.Password)
	if err != nil {
		logrus.Errorf("login failed: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "invalid email or password")
		return
	}
	expirationTime := time.Now().Add(time.Hour * 1)
	// Generate a token
	token, err := utils.GenerateToken(foundUser, expirationTime)
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
func (c *UserController) Logout(ctx *gin.Context) {

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

func (c *UserController) FindUserByPhone(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	phone := ctx.Param("phone_number")

	res, err := c.userRepo.GetByPhone(phone)
	if res == nil {
		logrus.Errorf("User not found : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusNotFound, "User not found")
		return
	}
	logrus.Info("Success to get user")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, res)

}

func (c *UserController) Register(ctx *gin.Context) {
	// Logging
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Failed to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	newUser := entity.User{}
	if err := ctx.BindJSON(&newUser); err != nil {
		logrus.Errorf("Invalid Input : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input")
		return
	}

	err = c.userRepo.Create(&newUser)
	if err != nil {
		logrus.Errorf("Failed to Register User : %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to Register User")
		return
	}

	logrus.Info("Success Register User")
	response.JSONSuccess(ctx.Writer, true, http.StatusCreated, "created user successfully")
}
