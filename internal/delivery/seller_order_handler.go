package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type SellerOrderHandler struct {
	service domain.SellerOrderService
}

func NewSellerOrderHandler(s domain.SellerOrderService) *SellerOrderHandler {
	return &SellerOrderHandler{
		service: s,
	}
}

func (h *SellerOrderHandler) GetAllSellerOrders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		res, err := h.service.GetAllSellerOrders(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully Fetch all Orders", "result": res})
	}
}

func (h *SellerOrderHandler) DetailSellerOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		orderID := ctx.Param("order_id")

		res, err := h.service.GetSellerOrderByEmailAndId(ctx, email, orderID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully Fetch the Order", "result": res})
	}
}

func (h *SellerOrderHandler) UpdateStatusOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.OrderStatusUpdateReq
		email := ctx.MustGet("email").(string)
		orderID := ctx.Param("order_id")
		productID := ctx.Query("product_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		err := h.service.UpdateSellerAndUserOrderStatus(ctx, email, orderID, productID, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully Update the Order Status"})
	}
}

func (h *SellerOrderHandler) CancelOrder() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		orderID := ctx.Param("order_id")
		productID := ctx.Query("product_id")

		err := h.service.CancelOrder(ctx, email, orderID, productID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully Cancel the Order"})
	}
}
