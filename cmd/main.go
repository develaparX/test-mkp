package main

import (
	"fmt"
	"log"
	"net/http"

	"sinibeli/internal/app/company"
	"sinibeli/internal/app/customer"
	"sinibeli/internal/app/product"
	"sinibeli/internal/app/transaction"
	"sinibeli/internal/config"
	"sinibeli/internal/infrastructure/cache"
	"sinibeli/internal/infrastructure/database"
	logger "sinibeli/internal/pkg/logging"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger.Init()

	db, err := database.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redisCache := cache.NewRedisCache(cfg.Cache)
	defer redisCache.Close()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := router.Group("/api/v1")

	txRepo := transaction.NewTransactionRepo(db.DB)
	customerRepo := customer.NewCustomerRepo(db.DB)
	productRepo := product.NewProductRepo(db.DB)
	companyRepo := company.NewCompanyRepo(db.DB)

	txService := transaction.NewTransactionService(txRepo, customerRepo, productRepo)
	transactionHandler := transaction.NewTransactionHandler(txService)
	trx := v1.Group("/transactions")
	{
		trx.POST("", transactionHandler.Create)
		trx.GET("", transactionHandler.GetAll)
		trx.GET("/:id", transactionHandler.GetByID)
		trx.GET("/summary", transactionHandler.GetTransactionSummaryFiltered)
		trx.GET("/reports", transactionHandler.GetCustomerActivity)
	}

	companyService := company.NewCompanyService(companyRepo)
	companyHandler := company.NewCompanyHandler(companyService)
	comp := v1.Group("/companies")
	{
		comp.POST("", companyHandler.Create)
		comp.GET("", companyHandler.GetAll)
		comp.GET("/:id", companyHandler.GetByID)
		comp.PUT("/:id", companyHandler.Update)
		comp.DELETE("/:id", companyHandler.Delete)
	}

	customerService := customer.NewCustomerService(customerRepo)
	customerHandler := customer.NewCustomerHandler(customerService)

	cust := v1.Group("/customers")
	{
		cust.POST("", customerHandler.Create)
		cust.GET("", customerHandler.GetAll)
		cust.GET("/:id", customerHandler.GetByID)
		cust.PUT("/:id", customerHandler.Update)
		cust.DELETE("/:id", customerHandler.Delete)
	}

	productService := product.NewProductService(productRepo)
	productHandler := product.NewProductHandler(productService)

	prod := v1.Group("/products")
	{
		prod.POST("", productHandler.Create)
		prod.GET("", productHandler.GetAll)
		prod.GET("/:id", productHandler.GetByID)
		prod.PUT("/:id", productHandler.Update)
		prod.DELETE("/:id", productHandler.Delete)
	}

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
