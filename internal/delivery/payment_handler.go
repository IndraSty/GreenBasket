package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service domain.PaymentService
}

func NewPaymentHandler(s domain.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		service: s,
	}
}

func (h *PaymentHandler) InitializePayment() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.PaymentReq

		userID := ctx.MustGet("uid").(string)
		req.UserID = userID

		orderId := ctx.Query("order_id")
		req.OrderID = orderId

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.service.InitializePayment(ctx, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Initialize payment successfully", "result": res})
	}
}
