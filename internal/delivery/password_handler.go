package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type PasswordHandler struct {
	service domain.PasswordService
}

func NewPasswordHandler(s domain.PasswordService) *PasswordHandler {
	return &PasswordHandler{
		service: s,
	}
}

func (h *PasswordHandler) ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		var req dto.PasswordReq

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		if err := h.service.ChangePassword(ctx, email, req); err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully Change Password"})
	}
}

func (h *PasswordHandler) RecoveryPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		var req dto.PasswordReq

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		if err := h.service.RecoveryPassword(ctx, email, req); err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully Recovery Password"})
	}
}
