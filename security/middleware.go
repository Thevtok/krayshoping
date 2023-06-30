package security

import (
	"krayshoping/model/response"
	"krayshoping/utils"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func AuthMiddleware() gin.HandlerFunc {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")

		if tokenString == "" {
			logrus.Errorf("unauthorized %v", err)

			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			logrus.Info("Claim Token Succesfully")
			return utils.JwtKey, nil

		})

		if err != nil || !token.Valid {
			logrus.Errorf("unauthorized %v", err)
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		claims := token.Claims.(*jwt.MapClaims)
		email, ok := (*claims)["email"].(string)
		if !ok {
			logrus.Errorf("invalid claim email")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid email claim")
			c.Abort()
			return
		}
		password, ok := (*claims)["password"].(string)
		if !ok {
			logrus.Errorf("invalid claim password ")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid password claim")
			c.Abort()
			return
		}
		userid, ok := (*claims)["user_id"].(string)
		if !ok {
			logrus.Errorf("invalid claim userID")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid userID claim")
			c.Abort()
			return
		}
		requestedID := c.Param("user_id")
		if userid != requestedID {
			response.JSONErrorResponse(c.Writer, false, http.StatusForbidden, "you do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("password", password)
		c.Set("user_id", userid)

		c.Next()
		logrus.Info("Success parsing midleware")
	}
}

func AuthMiddlewareMerchant() gin.HandlerFunc {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")

		if tokenString == "" {
			logrus.Errorf("unauthorized %v", err)

			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			logrus.Info("Claim Token Succesfully")
			return utils.JwtKey, nil

		})

		if err != nil || !token.Valid {
			logrus.Errorf("unauthorized %v", err)
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}

		claims := token.Claims.(*jwt.MapClaims)
		email, ok := (*claims)["email"].(string)
		if !ok {
			logrus.Errorf("invalid claim email")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid email claim")
			c.Abort()
			return
		}
		password, ok := (*claims)["password"].(string)
		if !ok {
			logrus.Errorf("invalid claim password ")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid password claim")
			c.Abort()
			return
		}
		merchantID, ok := (*claims)["merchant_id"].(string)
		if !ok {
			logrus.Errorf("invalid claim merchantID")
			response.JSONErrorResponse(c.Writer, false, http.StatusUnauthorized, "invalid merchantID claim")
			c.Abort()
			return
		}
		requestedID := c.Param("merchant_id")
		if merchantID != requestedID {
			response.JSONErrorResponse(c.Writer, false, http.StatusForbidden, "you do not have permission to access this resource")
			c.Abort()
			return
		}

		c.Set("email", email)
		c.Set("password", password)
		c.Set("merchant_id", merchantID)

		c.Next()
		logrus.Info("Success parsing midleware")
	}
}
