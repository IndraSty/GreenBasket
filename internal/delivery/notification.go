package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	service domain.NotificationService
	userSvc domain.UserService
}

func NewNotificationHandler(s domain.NotificationService, us domain.UserService) *NotificationHandler {
	return &NotificationHandler{
		service: s,
		userSvc: us,
	}
}

func (h *NotificationHandler) GetUserNotification() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		user, err := h.userSvc.GetUserByEmail(ctx, email)
		if err != nil {
			msg := "Something went wrong while fetching User data"
			util.HandleError(ctx, err, http.StatusInternalServerError, msg)
			return
		}

		notification, err := h.service.FindByUser(ctx, user.User_Id)
		if err != nil {
			msg := "Something went wrong while fetching User notification"
			util.HandleError(ctx, err, http.StatusInternalServerError, msg)
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully get user notification", "result": notification})
	}
}
