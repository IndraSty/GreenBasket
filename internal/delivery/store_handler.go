package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	service domain.StoreService
}

func NewStoreHandler(s domain.StoreService) *StoreHandler {
	return &StoreHandler{
		service: s,
	}
}

func (h *StoreHandler) CreateStore() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.StoreReq
		email := ctx.MustGet("email").(string)
		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.service.CreateStore(ctx, email, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"message": "Create store data successfully", "result": res})
	}
}

func (h *StoreHandler) EditStore() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.StoreReq
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")
		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		_, err := h.service.UpdateStore(ctx, email, storeID, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Update store data successfully"})

	}
}

func (h *StoreHandler) DetailStore() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		res, err := h.service.GetStoreByIdAndEmail(ctx, email, storeID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Fetching store data successfully", "data": res})
	}
}

func (h *StoreHandler) SearchStore() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query := ctx.Query("key")

		res, err := h.service.SearchStore(ctx, query)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		if len(res) == 0 {
			util.HandleError(ctx, nil, http.StatusNotFound, "no store was found")
			return
		}
		ctx.JSON(http.StatusCreated, gin.H{"message": "Search store by Name successfully", "data": res})
	}
}

func (h *StoreHandler) DeleteStore() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		_, err := h.service.DeleteStore(ctx, email, storeID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Delete store successfully"})
	}
}
