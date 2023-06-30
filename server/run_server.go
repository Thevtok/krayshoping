package server

import (
	"krayshoping/config"
	"krayshoping/controller"
	"krayshoping/repository"
	"krayshoping/security"
	"krayshoping/utils"
	"log"

	"github.com/gin-gonic/gin"
)

func RunServer() {

	db := config.LoadDatabase()
	defer db.Close()

	r := gin.Default()
	midlleware := security.AuthMiddleware()
	midllewareMerchant := security.AuthMiddlewareMerchant()

	userRepo := repository.NewUserRepository(db)
	userController := controller.NewUserController(userRepo)
	bankRepo := repository.NewBankAccRepository(db)
	bankController := controller.NewBankController(bankRepo)
	merchantRepo := repository.NewMerchantRepository(db)
	merchantController := controller.NewMerchantController(merchantRepo)
	productRepo := repository.NewProductRepository(db)
	productController := controller.NewProductController(productRepo)
	txRepo := repository.NewTxRepository(db)
	txController := controller.NewTransactionController(txRepo, userRepo, bankRepo, productRepo, merchantRepo)

	userRouter := r.Group("/user")

	userRouter.GET("/:phone_number", userController.FindUserByPhone)
	userRouter.POST("/login", userController.Login)
	userRouter.POST("/logout", userController.Logout)
	userRouter.POST("/register", userController.Register)

	bankRouter := r.Group("/bank")
	bankRouter.Use(midlleware)

	bankRouter.GET("/:user_id", bankController.FindBankAccByUserID)
	bankRouter.GET("/:user_id/:account_number", bankController.FindBankAccByAccountNumber)
	bankRouter.POST("/add/:user_id", bankController.CreateBankAccount)
	bankRouter.DELETE("delete/:user_id/:account_number", bankController.Unreg)

	merchantRouter := r.Group("/merchant")

	merchantRouter.POST("/login", merchantController.Login)
	merchantRouter.POST("/logout", merchantController.Logout)
	merchantRouter.POST("/register", merchantController.Register)

	productRouter := r.Group("/product")
	productRouter.Use(midllewareMerchant)

	productRouter.POST("/add/:merchant_id", productController.AddProduct)
	productRouter.PUT("/update/:merchant_id/:product_id", productController.UpdateProduct)
	productRouter.DELETE("/delete/:merchant_id/:product_id", productController.DeleteProduct)
	r.GET("/product", productController.GetAllProducts)

	txRouter := r.Group("/transaction")
	txRouter.Use(midlleware)

	txRouter.POST("deposit/:user_id/:account_number", txController.CreateDepositBank)
	txRouter.POST("transfer/:user_id", txController.CreateTransferTransaction)
	txRouter.GET("/:user_id", txController.GetTxByID)
	r.POST("notification", txController.HandlePaymentNotification)
	txRouter.POST("payment/:user_id/:merchant_id/:product_id", txController.CreatePaymentTransaction)
	txRouter.PUT("payment/:user_id/:merchant_id/:tx_id", txController.UpdatePaymentStatus)

	if err := r.Run(utils.DotEnv("SERVER_PORT")); err != nil {
		log.Fatal(err)
	}
}
