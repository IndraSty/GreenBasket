package delivery

import (
	"log"
	"net/http"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	dto "github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type SellerHandler struct {
	service domain.SellerService
}

func NewSellerHandler(s domain.SellerService) *SellerHandler {
	return &SellerHandler{
		service: s,
	}
}

func (sh *SellerHandler) RegisterSeller() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.SellerRegisterReq
		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		res, err := sh.service.RegisterSeller(ctx, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"result": res})
	}
}

func (sh *SellerHandler) AuthenticateSeller() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.SellerAuthReq
		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		res, err := sh.service.AuthenticateSeller(ctx, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:     "refresh_token",
			Value:    res.Refresh_Token,
			Expires:  time.Now().Add(168 * time.Hour),
			HttpOnly: true,
		})

		ctx.JSON(http.StatusOK, gin.H{"access_token": res.Access_Token})
	}
}

func (uh *SellerHandler) GetSellerHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		seller, err := uh.service.GetSellerByEmail(ctx, email)
		if err != nil {
			msg := "Something went wrong while fetching Seller data"
			util.HandleError(ctx, err, http.StatusInternalServerError, msg)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Fetch Seller Successfully", "data": seller})
	}
}

func (uh *SellerHandler) UpdateSellerHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		var sellerInput dto.SellerUpdateReq
		if err := ctx.BindJSON(&sellerInput); err != nil {
			log.Println("Error Input req seller:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := uh.service.UpdateSeller(ctx, email, &sellerInput)
		if err != nil {
			msg := "Error updating seller" + err.Error()
			util.HandleError(ctx, err, http.StatusInternalServerError, msg)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Update seller Successfully", "result": result})
	}
}
