package bootstrap

import (
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/config"
	"github.com/IndraSty/GreenBasket/internal/delivery"
	"github.com/IndraSty/GreenBasket/internal/middlewares"
	"github.com/IndraSty/GreenBasket/internal/repository"
	"github.com/IndraSty/GreenBasket/internal/routes"
	"github.com/IndraSty/GreenBasket/internal/service"
	"github.com/IndraSty/GreenBasket/internal/sse"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApplicationConfig struct {
	Config *config.Config
	Client *mongo.Client
	App    *gin.Engine
}

func Application(cnf *ApplicationConfig) {
	authSetup := config.NewAuthSetup(cnf.Config)
	hub := &dto.Hub{
		NotificationChannel: map[string]chan dto.NotificationRes{},
	}

	// setup repository
	cacheRepository := repository.NewRedisClient(cnf.Config)
	userRepository := repository.NewUserRepository(cnf.Client)
	sellerRepository := repository.NewSellerRepository(cnf.Client)
	storeRepository := repository.NewStoreRepository(cnf.Client)
	addressRepository := repository.NewAddressRepository(cnf.Client)
	contactRepository := repository.NewContactRepository(cnf.Client)
	productRepository := repository.NewProductRepository(cnf.Client)
	cartRepository := repository.NewCartRepository(cnf.Client)
	orderRepository := repository.NewOrderRepository(cnf.Client)
	notificationRepository := repository.NewNotificationRepository(cnf.Client)
	paymentRepository := repository.NewPaymentRepository(cnf.Client)
	sellerOrderRepository := repository.NewSellerOrderRepository(cnf.Client)
	templateRepository := repository.NewTemplateRepository(cnf.Client)
	reviewRepository := repository.NewReviewRepository(cnf.Client)
	salesReportRepository := repository.NewSalesReportRepository(cnf.Client)

	// setup service
	tokenService := util.NewTokenService(cnf.Config)
	addressService := service.NewAddressService(addressRepository, sellerRepository, userRepository, storeRepository)
	cartService := service.NewCartService(cartRepository, productRepository, storeRepository, cacheRepository)
	contactService := service.NewContactService(contactRepository, storeRepository)
	emailService := service.NewEmailService(cnf.Config)
	notificationService := service.NewNotificationService(notificationRepository, templateRepository, hub)
	salesReportService := service.NewSalesRepository(salesReportRepository, sellerOrderRepository, storeRepository, productRepository, reviewRepository, cacheRepository)
	orderService := service.NewOrderService(orderRepository, userRepository, cartRepository, sellerRepository,
		storeRepository, notificationService, sellerOrderRepository, salesReportService, cacheRepository)
	sellerService := service.NewSellerService(sellerRepository, tokenService, cacheRepository, emailService)
	midtransService := service.NewMidtransService(cnf.Config, paymentRepository, orderRepository, sellerOrderRepository)
	paymentService := service.NewPaymentService(notificationService, paymentRepository, userRepository, midtransService)
	productService := service.NewProductService(productRepository, storeRepository, salesReportRepository, cacheRepository)
	sellerOrderService := service.NewSellerOrderService(sellerOrderRepository, sellerRepository, orderRepository, productRepository, notificationService, cacheRepository)
	userService := service.NewUserService(userRepository, emailService, cacheRepository, cartService)
	reviewService := service.NewReviewService(reviewRepository, productRepository, orderRepository, storeRepository, notificationService, userRepository, salesReportRepository, cacheRepository)
	storeService := service.NewStoreService(storeRepository, sellerRepository, salesReportRepository, cacheRepository)
	authService := service.NewAuthService(userRepository, cacheRepository, tokenService, emailService)

	// setup handler
	authHandler := delivery.NewAuthHandler(userRepository, *authSetup, cnf.Config, authService)
	addressHandler := delivery.NewAddressHandler(addressService)
	cartHandler := delivery.NewCartHandler(cartService)
	contactHandler := delivery.NewContactHandler(contactService)
	midtransHandler := delivery.NewMidtransHandler(midtransService, paymentService)
	notificationHandler := delivery.NewNotificationHandler(notificationService, userService)
	orderHandler := delivery.NewOrderHandler(orderService)
	paymentHandler := delivery.NewPaymentHandler(paymentService)
	productHandler := delivery.NewProductHandler(productService)
	sellerHandler := delivery.NewSellerHandler(sellerService)
	storeHandler := delivery.NewStoreHandler(storeService)
	sellerOrderHandler := delivery.NewSellerOrderHandler(sellerOrderService)
	userHandler := delivery.NewUserHandler(userService)
	reviewHandler := delivery.NewReviewHandler(reviewService)
	salesReportHandler := delivery.NewSalesReportHandler(salesReportService)
	notificationSSE := sse.NewNotificationSSE(hub, userRepository)

	// setup middleware
	middleware := middlewares.NewMiddleware(tokenService)

	// setup routes
	routeConfig := routes.RouteConfig{
		App:                 cnf.App,
		Middlewares:         middleware,
		UserHandler:         userHandler,
		SellerHandler:       sellerHandler,
		StoreHandler:        storeHandler,
		ProductHandler:      productHandler,
		PaymentHandler:      paymentHandler,
		OrderHandler:        orderHandler,
		NotificationHandler: notificationHandler,
		MidtransHandler:     midtransHandler,
		ContactHandler:      contactHandler,
		CartHandler:         cartHandler,
		AddressHandler:      addressHandler,
		NotificationSSE:     notificationSSE,
		SellerOrderHandler:  sellerOrderHandler,
		SalesReportHandler:  salesReportHandler,
		ReviewHandler:       reviewHandler,
		AuthHandler:         authHandler,
	}

	routeConfig.Setup()

	// setup sse
	sse.NewNotificationSSE(hub, userRepository)
}
