package main

import (
	"github.com/gin-gonic/gin"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/config"
	authinfra "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/infrastructure/auth"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/infrastructure/mongo"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/controllers"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/middlewares"
	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/interface/routes"

	authusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/auth"
	cartitemusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/cartitem"

	bundleusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/bundle"
	orderusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/order"
	productusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/product"
	reviewusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/review"
	trustusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/trust"
	userusecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/user"
	warehouse_usecase "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/usecase/warehouse"
)

func main() {
	// Load .env variables
	config.LoadEnv()

	// Set Gin to release mode
	gin.SetMode(gin.DebugMode)

	// Load grouped app config
	appConfig := config.LoadAppConfig()

	// Connect to MongoDB
	db := config.ConnectMongo(appConfig.DBURI, appConfig.DBName)

	// Init shared services
	jwtSvc := authinfra.NewJWTService(appConfig.JWTSecret)
	passSvc := authinfra.NewPasswordService()

	// Init Repositories
	userRepo := mongo.NewMongoUserRepository(db)
	productRepo := mongo.NewMongoProductRepository(db)
	bundleRepo := mongo.NewBundleRepository(db)
	orderRepo := mongo.NewMongoOrderRepository(db) // Add order repository
	cartItemRepo := mongo.NewCartItemRepository(db)
	reviewRepo := mongo.NewReviewRepository(db)            // Add review repository
	warehouseRepo := mongo.NewMongoWarehouseRepository(db) // Add warehouse repository
	paymentRepo := mongo.NewMongoPaymentRepository(db)     // Add payment repository

	// Init Usecases
	userUC := userusecase.NewUserUsecase(userRepo)
	authUC := authusecase.NewAuthUsecase(userRepo, passSvc, jwtSvc)
	productUC := productusecase.NewProductUsecase(productRepo, bundleRepo)
	bundleUC := bundleusecase.NewBundleUsecase(bundleRepo)
	trustUC := trustusecase.NewTrustUsecase(productRepo, bundleRepo, userRepo)
	orderUC := orderusecase.NewOrderUsecase(
		bundleRepo,
		orderRepo,
		warehouseRepo,
		paymentRepo,
		userRepo,
		productRepo,
	)
	cartItemUC := cartitemusecase.NewCartItemUsecase(cartItemRepo, productRepo, paymentRepo, orderUC, orderRepo)

	reviewUC := reviewusecase.NewReviewUsecase(reviewRepo, orderRepo) // Add review usecase
	warehouseSvc := warehouse_usecase.NewWarehouseUseCase(warehouseRepo, bundleRepo)

	// Init Controllers
	authCtrl := controllers.NewAuthController(authUC)
	adminCtrl := controllers.NewAdminController(userUC, orderUC)
	productCtrl := controllers.NewProductController(productUC, trustUC, bundleUC, warehouseRepo)
	bundleCtrl := controllers.NewBundleController(bundleUC, userUC)
	consumerCtrl := controllers.NewConsumerController(orderRepo)
	supplierCtrl := controllers.NewSupplierController(orderUC) // Add consumer controller
	cartItemCtrl := controllers.NewCartItemController(cartItemUC, productUC)
	reviewCtrl := controllers.NewReviewController(reviewUC, trustUC, productUC) // Add trust and product usecases
	warehouseCtrl := controllers.NewWarehouseController(warehouseSvc)
	orderCtrl := controllers.NewOrderController(orderUC) // Add order controller

	// Init Gin Engine and Routes
	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())

	// Add a health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK"})
	})

	routes.RegisterAuthRoutes(r, authCtrl)
	routes.RegisterProductRoutes(r, productCtrl, jwtSvc, reviewCtrl, trustUC, productUC)
	routes.RegisterAdminRoutes(r, adminCtrl, jwtSvc)
	routes.RegisterBundleRoutes(r, bundleCtrl, jwtSvc)
	routes.RegisterCartItemRoutes(r, cartItemCtrl, jwtSvc)
	routes.RegisterOrderRoutes(r, orderCtrl, consumerCtrl, jwtSvc)
	routes.RegisterSupplierRoutes(r, supplierCtrl, jwtSvc)
	routes.RegisterWarehouseRoutes(r, warehouseCtrl, jwtSvc)
	routes.RegisterResellerRoutes(r, supplierCtrl, jwtSvc)
	routes.SetupUserRoutes(r, userUC, jwtSvc) // Add user routes

	// Run server
	r.Run(":8080")
}
