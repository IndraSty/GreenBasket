package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	service domain.ContactService
}

func NewContactHandler(s domain.ContactService) *ContactHandler {
	return &ContactHandler{
		service: s,
	}
}

func (h *ContactHandler) AddStoreContact() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.Contact
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		_, err := h.service.AddStoreContact(ctx, email, storeID, req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Store Contact successfully added"})
	}
}

func (h *ContactHandler) GetStoreContact() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		res, err := h.service.GetStoreContact(ctx, email, storeID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"store_contact": res})
	}
}

func (h *ContactHandler) EditStoreContact() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.Contact
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		_, err := h.service.UpdateStoreContact(ctx, email, storeID, req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Store Contact successfully updated"})
	}
}

func (h *ContactHandler) DeleteStoreContact() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		_, err := h.service.RemoveStoreContact(ctx, email, storeID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Store Contact successfully deleted"})
	}
}
