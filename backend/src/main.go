package main

import (
	"log"
	"net/http"
	"server/src/config"
	"server/src/delivery/handlers"
	"server/src/delivery/middlewares"
	"server/src/delivery/routes"
	"server/src/repository"
	"server/src/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	config.InitDB()
	config.ConnectRdb()
	limiter := middlewares.NewIPLimiter(5, 10)

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	err := config.Rdb.Ping(config.Ctx).Err()

	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}

	r := gin.Default()

	r.Use(middlewares.RateLimitMiddleware(limiter))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Route not found!",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	// Inisialisasi repo, service, & handler
	authRepo := repository.NewAuthRepo(config.DB, config.Rdb)
	authService := service.NewAuthService(authRepo)
	authHandlers := handlers.NewAuthHandler(authService)

	clientRepo := repository.NewClientRepo(config.DB, config.Rdb)
	clientService := service.NewClientService(clientRepo)
	clientHandlers := handlers.NewClientHandler(clientService)

	productRepo := repository.NewProductRepo(config.DB, config.Rdb)
	productService := service.NewProductService(productRepo)
	productHandlers := handlers.NewProductHandler(productService)

	studentRepo := repository.NewStudentRepo(config.DB, config.Rdb)
	studentService := service.NewStudentService(studentRepo)
	studentHandlers := handlers.NewStudentHandler(studentService)

	transactionRepo := repository.NewTransactionRepo(config.DB, config.Rdb)
	transactionService := service.NewTransactionService(transactionRepo, studentRepo, clientRepo)
	transactionHandlers := handlers.NewTransactionHandler(transactionService)

	// API
	api := r.Group("/api")
	{
		routes.AuthRoutes(api, authHandlers)
		routes.ClientRoutes(api, clientHandlers)
		routes.Product(api, productHandlers)
		routes.Student(api, studentHandlers)
		routes.TransactionRoutes(api, transactionHandlers)
	}

	r.Run("localhost:3000")
}
