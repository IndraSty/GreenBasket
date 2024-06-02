package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	service domain.AddressService
}

func NewAddressHandler(s domain.AddressService) *AddressHandler {
	return &AddressHandler{
		service: s,
	}
}

// user address handler
func (h *AddressHandler) AddUserAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.Address
		email := ctx.MustGet("email").(string)

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		_, err := h.service.AddUserAddress(ctx, email, req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully add user address"})
	}
}

func (h *AddressHandler) GetUserAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		res, err := h.service.GetUserAddress(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully get user address", "result": res})
	}
}

func (h *AddressHandler) UpdateUserAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.Address
		email := ctx.MustGet("email").(string)

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		_, err := h.service.UpdateUserAddress(ctx, email, req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully update user address"})
	}
}

func (h *AddressHandler) RemoveUserAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		_, err := h.service.RemoveUserAddress(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully remove user address"})
	}
}

// seller address handler
func (h *AddressHandler) AddSellerAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.Address
		email := ctx.MustGet("email").(string)

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		_, err := h.service.AddSellerAddress(ctx, email, req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully add seller address"})
	}
}

func (h *AddressHandler) GetSellerAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		res, err := h.service.GetSellerAddress(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully get seller address", "result": res})
	}
}

func (h *AddressHandler) UpdateSellerAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.Address
		email := ctx.MustGet("email").(string)

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		_, err := h.service.UpdateSellerAddress(ctx, email, req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully update seller address"})
	}
}

func (h *AddressHandler) RemoveSellerAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		_, err := h.service.RemoveSellerAddress(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully remove seller address"})
	}
}

// store address handler
func (h *AddressHandler) AddStoreAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.Address
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		_, err := h.service.AddStoreAddress(ctx, email, storeID, req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully add store address"})
	}
}

func (h *AddressHandler) GetStoreAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		res, err := h.service.GetStoreAddress(ctx, email, storeID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully get store address", "result": res})
	}
}

func (h *AddressHandler) EditStoreAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req domain.Address
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		_, err := h.service.UpdateStoreAddress(ctx, email, storeID, req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully update store address"})
	}
}

func (h *AddressHandler) RemoveStoreAddress() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		_, err := h.service.RemoveStoreAddress(ctx, email, storeID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Successfully remove store address"})
	}
}
