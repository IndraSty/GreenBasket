package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service domain.OrderService
}

func NewOrderHandler(s domain.OrderService) *OrderHandler {
	return &OrderHandler{
		service: s,
	}
}

func (h *OrderHandler) CreateOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		res, err := h.service.CreateOrder(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully create an order", "result": res})
	}
}

func (h *OrderHandler) GetAllOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		res, err := h.service.GetAllOrders(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully Fetch all Orders", "result": res})
	}
}

func (h *OrderHandler) DetailOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		orderID := ctx.Param("order_id")

		res, err := h.service.GetOrderByEmailAndId(ctx, email, orderID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully Fetch the Order", "result": res})
	}
}

func (h *OrderHandler) FinishOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.OrderStatusUpdateReq
		email := ctx.MustGet("email").(string)
		orderID := ctx.Param("order_id")
		productId := ctx.Query("product_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		err := h.service.FinishOrder(ctx, email, orderID, productId, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully Finish status the Order"})
	}
}

func (h *OrderHandler) CancelOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		orderID := ctx.Param("order_id")
		productID := ctx.Query("product_id")

		err := h.service.CancelOrder(ctx, email, orderID, productID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully Cancel the Order"})
	}
}
