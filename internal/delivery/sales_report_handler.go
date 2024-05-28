package delivery

import (
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type SalesReportHandler struct {
	service domain.SalesReportService
}

func NewSalesReportHandler(s domain.SalesReportService) *SalesReportHandler {
	return &SalesReportHandler{
		service: s,
	}
}

func (h *SalesReportHandler) GetSalesReport() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		res, err := h.service.GetSalesReport(ctx, email, storeID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Get Sales Report Successfully", "result": res})
	}
}
