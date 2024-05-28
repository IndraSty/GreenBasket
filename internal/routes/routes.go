package routes

import (
	"github.com/IndraSty/GreenBasket/internal/delivery"
	"github.com/IndraSty/GreenBasket/internal/middlewares"
	"github.com/IndraSty/GreenBasket/internal/sse"
	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	App                 *gin.Engine
	Middlewares         *middlewares.Middleware
	UserHandler         *delivery.UserHandler
	SellerHandler       *delivery.SellerHandler
	StoreHandler        *delivery.StoreHandler
	ProductHandler      *delivery.ProductHandler
	PaymentHandler      *delivery.PaymentHandler
	OrderHandler        *delivery.OrderHandler
	SellerOrderHandler  *delivery.SellerOrderHandler
	NotificationHandler *delivery.NotificationHandler
	MidtransHandler     *delivery.MidtransHandler
	ContactHandler      *delivery.ContactHandler
	CartHandler         *delivery.CartHandler
	AddressHandler      *delivery.AddressHandler
	ReviewHandler       *delivery.ReviewHandler
	SalesReportHandler  *delivery.SalesReportHandler
	AuthHandler         *delivery.AuthHandler
	PasswordHandler     *delivery.PasswordHandler
	NotificationSSE     *sse.NotificationSSE
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupSellerAuthRoute()
	c.SetupUserAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {
	// user route

	c.App.POST("/api/users/signup", c.UserHandler.RegisterUser())
	c.App.POST("/api/users/login", c.AuthHandler.AuthenticateUser())
	c.App.POST("/api/users/otp", c.AuthHandler.ValidateOTP())
	c.App.POST("/api/users/email-verify-request", c.AuthHandler.RequestVerifyEmail())

	// google login
	c.App.GET("/api/auth/:provider", c.AuthHandler.BeginAuthHandler)
	c.App.GET("/api/auth/callback/:provider", c.AuthHandler.GetAuthCallBackFunc)
	c.App.GET("/api/auth/logout", c.AuthHandler.LogoutHandler)

	// seller route
	c.App.POST("/api/sellers/signup", c.SellerHandler.RegisterSeller())
	c.App.POST("/api/sellers/login", c.SellerHandler.AuthenticateSeller())

	// product for guest
	c.App.GET("/api/products", c.ProductHandler.FetchAllProductForGuest())
	c.App.GET("/api/products/search", c.ProductHandler.SearchProductForGuest())
	c.App.GET("/api/products/:product_id", c.ProductHandler.FetchProductForGuest())
	c.App.GET("/api/products/category", c.ProductHandler.FetchAllProductByCategoryForGuest())
	c.App.GET("/api/products/sort", c.ProductHandler.SortProductForGuest())

	// store for guest
	c.App.GET("/api/stores", c.StoreHandler.SearchStore())

	// midtrans callback
	c.App.POST("/api/midtrans/payment-callback", c.MidtransHandler.PaymentHandlerNotification())
}

func (c *RouteConfig) SetupSellerAuthRoute() {
	sellerRoutes := c.App.Group("/api/sellers")
	{
		sellerRoutes.Use(c.Middlewares.SellerAuthMiddleware())
		sellerRoutes.GET("/current", c.SellerHandler.GetSellerHandler())
		sellerRoutes.PUT("/current", c.SellerHandler.UpdateSellerHandler())

		// seller address
		sellerRoutes.POST("/current/addresses", c.AddressHandler.AddSellerAddress())
		sellerRoutes.GET("/current/addresses", c.AddressHandler.GetSellerAddress())
		sellerRoutes.PUT("/current/addresses", c.AddressHandler.UpdateSellerAddress())
		sellerRoutes.DELETE("/current/addresses", c.AddressHandler.RemoveSellerAddress())

		// seller store
		sellerRoutes.POST("/current/stores", c.StoreHandler.CreateStore())
		sellerRoutes.GET("/current/stores/:store_id", c.StoreHandler.DetailStore())
		sellerRoutes.PUT("/current/stores/:store_id", c.StoreHandler.EditStore())
		sellerRoutes.DELETE("/current/stores/:store_id", c.StoreHandler.DeleteStore())

		// seller store address
		sellerRoutes.POST("/current/stores/:store_id/address", c.AddressHandler.AddStoreAddress())
		sellerRoutes.GET("/current/stores/:store_id/address", c.AddressHandler.GetStoreAddress())
		sellerRoutes.PUT("/current/stores/:store_id/address", c.AddressHandler.EditStoreAddress())
		sellerRoutes.DELETE("/current/stores/:store_id/address", c.AddressHandler.RemoveStoreAddress())

		// seller store contact
		sellerRoutes.POST("/current/stores/:store_id/contact", c.ContactHandler.AddStoreContact())
		sellerRoutes.GET("/current/stores/:store_id/contact", c.ContactHandler.GetStoreContact())
		sellerRoutes.PUT("/current/stores/:store_id/contact", c.ContactHandler.EditStoreContact())
		sellerRoutes.DELETE("/current/stores/:store_id/contact", c.ContactHandler.DeleteStoreContact())

		// seller store product
		sellerRoutes.POST("/current/stores/:store_id/product", c.ProductHandler.AddProduct())
		sellerRoutes.GET("/current/stores/:store_id/product", c.ProductHandler.FetchProductById())
		sellerRoutes.GET("/current/stores/:store_id/products/category", c.ProductHandler.FetchAllProductByCategorySeller())
		sellerRoutes.GET("/current/stores/:store_id/products", c.ProductHandler.FetchAllProductSeller())
		sellerRoutes.GET("/current/stores/:store_id/products/search", c.ProductHandler.SearchProduct())
		sellerRoutes.GET("/current/stores/:store_id/products/sort", c.ProductHandler.SortProduct())
		sellerRoutes.PUT("/current/stores/:store_id/product", c.ProductHandler.UpdateProduct())
		sellerRoutes.DELETE("/current/stores/:store_id/product", c.ProductHandler.DeleteProduct())

		// seller order
		sellerRoutes.GET("/current/order/:order_id", c.SellerOrderHandler.DetailSellerOrder())
		sellerRoutes.GET("/current/orders", c.SellerOrderHandler.GetAllSellerOrders())
		sellerRoutes.PATCH("/current/order/:order_id", c.SellerOrderHandler.UpdateStatusOrder())
		sellerRoutes.DELETE("/current/order/:order_id", c.SellerOrderHandler.CancelOrder())

		// seller review
		sellerRoutes.GET("/current/product/reviews", c.ReviewHandler.GetAllReviewByProductId())
		sellerRoutes.GET("/current/reviews", c.ReviewHandler.GetAllReviewBySellerEmail())
		sellerRoutes.PATCH("/current/reviews/:review_id", c.ReviewHandler.UpdateResponSeller())

		// seller sales report
		sellerRoutes.GET("/current/stores/:store_id/report", c.SalesReportHandler.GetSalesReport())
	}
}

func (c *RouteConfig) SetupUserAuthRoute() {
	c.App.Use(c.Middlewares.UserAuthMiddleware())
	userRoutes := c.App.Group("/api/users")
	{
		userRoutes.Use(c.Middlewares.UserAuthMiddleware())
		userRoutes.GET("/current", c.UserHandler.GetUserHandler())
		userRoutes.PUT("/current", c.UserHandler.UpdateUserHandler())
		userRoutes.PUT("/current/phone", c.UserHandler.AddPhoneNumber())

		// user address
		userRoutes.POST("/current/addresses", c.AddressHandler.AddUserAddress())
		userRoutes.GET("/current/addresses", c.AddressHandler.GetUserAddress())
		userRoutes.PUT("/current/addresses", c.AddressHandler.UpdateUserAddress())
		userRoutes.DELETE("/current/addresses", c.AddressHandler.RemoveUserAddress())

		// user cart
		userRoutes.POST("/current/cart", c.CartHandler.AddToCart())
		userRoutes.GET("/current/cart", c.CartHandler.GetCart())
		userRoutes.PUT("/current/cartitem", c.CartHandler.UpdateItemInCart())
		userRoutes.DELETE("/current/cartitem", c.CartHandler.RemoveItemInCart())
		userRoutes.GET("/current/cartitems", c.CartHandler.GetAllItemCart())

		// user order
		userRoutes.POST("/current/order", c.OrderHandler.CreateOrder())
		userRoutes.GET("/current/order/:order_id", c.OrderHandler.DetailOrder())
		userRoutes.PATCH("/current/order/:order_id", c.OrderHandler.FinishOrder())
		userRoutes.GET("/current/orders", c.OrderHandler.GetAllOrders())
		userRoutes.DELETE("/current/order/:order_id", c.OrderHandler.CancelOrder())

		// user payment
		userRoutes.POST("/current/payment", c.PaymentHandler.InitializePayment())

		// user review
		userRoutes.POST("/current/review/:order_id", c.ReviewHandler.AddReview())
		userRoutes.GET("/current/review/:review_id", c.ReviewHandler.GetUserReviewById())
		userRoutes.GET("/current/reviews", c.ReviewHandler.GetAllReviewByUserEmail())
		userRoutes.PATCH("/current/review/:review_id", c.ReviewHandler.UpdateReview())
		userRoutes.DELETE("/current/review/:review_id", c.ReviewHandler.DeleteReview())

	}

	// sse notification user
	c.App.GET("/sse/notification-stream", c.NotificationSSE.StreamNotification())
}
