package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type CartHandler struct {
	service domain.CartService
}

func NewCartHandler(s domain.CartService) *CartHandler {
	return &CartHandler{
		service: s,
	}
}

func (h *CartHandler) GetCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		res, err := h.service.GetUserCart(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Fetching cart user successfully", "data": res})
	}
}

func (h *CartHandler) AddToCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.AddCartReq
		email := ctx.MustGet("email").(string)
		productID := ctx.Query("product_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		err := h.service.AddToCart(ctx, email, productID, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully added product to the cart"})
	}
}

func (h *CartHandler) GetAllItemCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		res, err := h.service.GetAllCartItem(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully Fetch All cart items", "items": res})
	}
}

func (h *CartHandler) UpdateItemInCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.CartItemEditReq
		email := ctx.MustGet("email").(string)
		productID := ctx.Query("product_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		err := h.service.UpdateCartItemById(ctx, email, productID, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully update product in the cart"})
	}
}

func (h *CartHandler) RemoveItemInCart() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		productID := ctx.Query("product_id")

		err := h.service.RemoveCartItemById(ctx, email, productID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully remove product from the cart"})
	}
}
