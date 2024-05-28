package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type MidtransHandler struct {
	service    domain.MidtransService
	paymentSvc domain.PaymentService
}

func NewMidtransHandler(svc domain.MidtransService, paymentSvc domain.PaymentService) *MidtransHandler {
	return &MidtransHandler{
		service:    svc,
		paymentSvc: paymentSvc,
	}
}

func (h *MidtransHandler) PaymentHandlerNotification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var notificationPayload map[string]interface{}
		if err := ctx.BindJSON(&notificationPayload); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		orderId, exists := notificationPayload["order_id"].(string)
		if !exists {
			// do something when key `order_id` not found
			ctx.Status(http.StatusBadRequest)
		}

		success, err := h.service.VerifyPayment(ctx, orderId)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		if success {
			_ = h.paymentSvc.ConfirmedPayment(ctx, orderId)
		}

		ctx.Status(http.StatusOK)
	}
}
