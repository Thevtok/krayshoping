package controller

import (
	"krayshoping/model/entity"
	"krayshoping/model/response"
	"krayshoping/repository"
	"krayshoping/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ProductController struct {
	productRepo repository.ProductRepository
}

func NewProductController(productRepo repository.ProductRepository) *ProductController {
	controller := ProductController{
		productRepo: productRepo,
	}
	return &controller
}

func (c *ProductController) AddProduct(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)
	mID := ctx.Param("merchant_id")

	var newProduct entity.MerchantProduct
	err = ctx.BindJSON(&newProduct)
	if err != nil {
		logrus.Errorf("Invalid request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid request body")
		return
	}
	newProduct.ProductName = strings.ToLower(newProduct.ProductName)

	if newProduct.ProductName == "" || newProduct.ProductPrice == 0 || newProduct.ProductQuantity == 0 {
		logrus.Errorf("Invalid Input: Required fields are empty")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input: Required fields are empty")
		return
	}
	newProduct.MerchantID = mID

	err = c.productRepo.AddProduct(&newProduct)
	if err != nil {
		logrus.Errorf("failed to add product %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "failed to add product")
		return
	}

	// Nomor akun sudah ada dalam database
	logrus.Info("add product succesfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, "add product succesfully")

}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Mendapatkan product_id dari URL parameter
	productIDStr := ctx.Param("product_id")
	productID, _ := strconv.Atoi(productIDStr)

	var updatedProduct entity.MerchantProduct
	err = ctx.BindJSON(&updatedProduct)
	if err != nil {
		logrus.Errorf("Invalid request body: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid request body")
		return
	}
	updatedProduct.ProductName = strings.ToLower(updatedProduct.ProductName)

	if updatedProduct.ProductName == "" || updatedProduct.ProductPrice == 0 || updatedProduct.ProductQuantity == 0 {
		logrus.Errorf("Invalid Input: Required fields are empty")
		response.JSONErrorResponse(ctx.Writer, false, http.StatusBadRequest, "Invalid Input: Required fields are empty")
		return
	}

	// Set product_id pada updatedProduct
	updatedProduct.Product_ID = productID

	err = c.productRepo.UpdateProduct(&updatedProduct)
	if err != nil {
		logrus.Errorf("Failed to update product: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to update product")
		return
	}

	logrus.Info("Update product successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, "Update product successfully")
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	// Mendapatkan product_id dari URL parameter
	productIDStr := ctx.Param("product_id")
	productID, _ := strconv.Atoi(productIDStr)

	err = c.productRepo.DeleteProduct(productID)
	if err != nil {
		logrus.Errorf("Failed to delete product: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	logrus.Info("Delete product successfully")
	response.JSONSuccess(ctx.Writer, true, http.StatusOK, "Delete product successfully")
}

func (c *ProductController) GetAllProducts(ctx *gin.Context) {
	logger, err := utils.CreateLogFile()
	if err != nil {
		log.Fatalf("Fatal to create log file: %v", err)
	}

	logrus.SetOutput(logger)

	products, err := c.productRepo.GetAllProducts()
	if err != nil {
		logrus.Errorf("Failed to get products: %v", err)
		response.JSONErrorResponse(ctx.Writer, false, http.StatusInternalServerError, "Failed to get products")
		return
	}

	response.JSONSuccess(ctx.Writer, true, http.StatusOK, products)
}
