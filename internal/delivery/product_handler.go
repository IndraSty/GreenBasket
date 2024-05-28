package delivery

import (
	"net/http"
	"strconv"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(s domain.ProductService) *ProductHandler {
	return &ProductHandler{
		service: s,
	}
}

// seller
func (h *ProductHandler) AddProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.ProductReq
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.service.CreateProduct(ctx, storeID, email, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "Add Product successfully", "result": res})
	}
}

func (h *ProductHandler) UpdateProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.ProductReq
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")

		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.service.CreateProduct(ctx, email, storeID, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Update Product successfully", "result": res})
	}
}

func (h *ProductHandler) FetchProductById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")
		productID := ctx.Query("id")

		res, err := h.service.GetProductById(ctx, storeID, email, productID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Fetch product Successfully", "product": res})
	}
}

func (h *ProductHandler) FetchAllProductSeller() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")
		pageStr := ctx.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)

		res, err := h.service.GetAllProduct(ctx, storeID, email, page)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Fetch All product Successfully", "result": res})
	}
}

func (h *ProductHandler) FetchAllProductByCategorySeller() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")
		category := ctx.Query("key")
		pageStr := ctx.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)

		res, err := h.service.GetAllByCategory(ctx, email, storeID, category, page)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Fetch All product seller by category Successfully", "result": res})
	}
}

func (h *ProductHandler) SearchProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")
		query := ctx.Query("key")
		pageStr := ctx.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)

		res, err := h.service.SearchProduct(ctx, email, storeID, query, page)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Search product Successfully", "result": res})
	}
}

func (h *ProductHandler) SortProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")
		pageStr := ctx.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)

		sortParams := make(map[string]string)
		for _, param := range []string{"price", "stock", "created_at", "total_sales", "average_rating"} {
			if direction := ctx.Query(param); direction != "" {
				sortParams[param] = direction
			}
		}

		res, err := h.service.GetAllProductSorted(ctx, sortParams, page, email, storeID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Sort product Successfully", "result": res})
	}
}

func (h *ProductHandler) DeleteProduct() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		storeID := ctx.Param("store_id")
		productID := ctx.Query("id")

		res, err := h.service.DeleteProductById(ctx, storeID, email, productID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Product successfully deleted", "result": res})
	}
}

// user / guest

func (h *ProductHandler) FetchAllProductForGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageStr := ctx.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)
		res, err := h.service.GetAllProductForGuest(ctx, page)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Fetch All product for guest Successfully", "result": res})
	}
}

func (h *ProductHandler) SearchProductForGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query := ctx.Query("key")
		sort := ctx.Query("sort")
		pageStr := ctx.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)
		res, err := h.service.SearchProductForGuest(ctx, page, query, sort)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Search product for guest Successfully", "result": res})
	}
}

func (h *ProductHandler) FetchAllProductByCategoryForGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		category := ctx.Query("category")
		pageStr := ctx.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)

		res, err := h.service.GetAllByCategoryForGuest(ctx, category, page)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Fetch All product by category for guest Successfully", "result": res})
	}
}

func (h *ProductHandler) FetchProductForGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		productID := ctx.Param("product_id")
		res, err := h.service.GetProductByIdForGuest(ctx, productID)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Fetch product for guest Successfully", "result": res})
	}
}

func (h *ProductHandler) SortProductForGuest() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pageStr := ctx.DefaultQuery("page", "1")
		page, _ := strconv.Atoi(pageStr)
		sortParams := make(map[string]string)
		for _, param := range []string{"price", "stock", "created_at", "total_sales", "average_rating"} {
			if direction := ctx.Query(param); direction != "" {
				sortParams[param] = direction
			}
		}

		res, err := h.service.GetAllProductSortedForCust(ctx, sortParams, page)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Sort product for Cust Successfully", "result": res})
	}
}
