package utils

import (
	"krayshoping/model/entity"

	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

var JwtKey = []byte(DotEnv("KEY"))

func GenerateToken(user *entity.User, expirationTime time.Time) (string, error) {
	logger, err := CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	// Set token claims
	claims := jwt.MapClaims{}
	claims["email"] = user.Email
	claims["password"] = user.Password

	claims["user_id"] = user.ID

	claims["exp"] = expirationTime.Unix()

	// Create token with claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		logrus.Errorf("Failed to generate token: %v", err)
		return "", err
	}

	logrus.Info("Claim Token Successfully")

	return tokenString, nil
}

func GenerateTokenMerchant(user *entity.Merchant, expirationTime time.Time) (string, error) {
	logger, err := CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	// Set token claims
	claims := jwt.MapClaims{}
	claims["email"] = user.Email
	claims["password"] = user.Password
	claims["merchant_name"] = user.MerchantName
	claims["merchant_id"] = user.ID

	claims["exp"] = expirationTime.Unix()

	// Create token with claims and secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		logrus.Errorf("Failed to generate token: %v", err)
		return "", err
	}

	logrus.Info("Claim Token Successfully")

	return tokenString, nil
}
