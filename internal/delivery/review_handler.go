package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	service domain.ReviewService
}

func NewReviewHandler(s domain.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		service: s,
	}
}

func (h *ReviewHandler) AddReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.AddReviewReq
		email := ctx.MustGet("email").(string)
		orderID := ctx.Param("order_id")
		productID := ctx.Query("product_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.service.CreateReview(ctx, email, orderID, productID, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Add Review Successfully", "result": res})
	}
}

func (h *ReviewHandler) GetUserReviewById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		reviewId := ctx.Param("review_id")

		res, err := h.service.GetUserReviewById(ctx, email, reviewId)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Get User Review By Id Successfully", "result": res})
	}
}

func (h *ReviewHandler) DeleteReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		reviewId := ctx.Param("review_id")

		err := h.service.DeleteReview(ctx, email, reviewId)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Delete Review Successfully"})
	}
}

func (h *ReviewHandler) UpdateReview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.AddReviewReq
		email := ctx.MustGet("email").(string)
		reviewId := ctx.Param("review_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		err := h.service.UpdateReview(ctx, email, reviewId, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Update Review Successfully"})
	}
}

func (h *ReviewHandler) GetAllReviewByProductId() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		productId := ctx.Query("product_id")

		res, err := h.service.GetAllReviewByProductId(ctx, email, productId)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Get Review By Product Id Successfully", "result": res})
	}
}

func (h *ReviewHandler) GetAllReviewBySellerEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		res, err := h.service.GetAllReviewBySellerEmail(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Get Review By Seller Email Successfully", "result": res})
	}
}

func (h *ReviewHandler) GetAllReviewByUserEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)

		res, err := h.service.GetAllReviewByUserEmail(ctx, email)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Get Review By User Email Successfully", "result": res})
	}
}

func (h *ReviewHandler) UpdateResponSeller() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.ResponSellerReq
		email := ctx.MustGet("email").(string)
		reviewId := ctx.Param("review_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		err := h.service.UpdateResponSeller(ctx, email, reviewId, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Add Seller Response Successfully"})
	}
}
