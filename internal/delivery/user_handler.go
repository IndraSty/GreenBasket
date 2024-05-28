package delivery

import (
	"log"
	"net/http"

	"github.com/IndraSty/GreenBasket/domain"
	dto "github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc domain.UserService
}

func NewUserHandler(us domain.UserService) *UserHandler {
	return &UserHandler{
		userSvc: us,
	}
}

func (uh *UserHandler) RegisterUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var req dto.UserRegisterReq
		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		res, err := uh.userSvc.RegisterUser(ctx, &req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "User created successfully!", "result": res})
	}
}

func (uh *UserHandler) GetUserHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.MustGet("email").(string)
		user, err := uh.userSvc.GetUserByEmail(ctx, userID)
		if err != nil {
			msg := "Something went wrong while fetching User data"
			util.HandleError(ctx, err, http.StatusInternalServerError, msg)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": user})
	}
}

func (uh *UserHandler) UpdateUserHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		var userInput dto.UserUpdateReq
		if err := ctx.BindJSON(&userInput); err != nil {
			log.Println("Error Input req user:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		result, err := uh.userSvc.UpdateUser(ctx, email, &userInput)
		if err != nil {
			msg := "Error updating user" + err.Error()
			util.HandleError(ctx, err, http.StatusInternalServerError, msg)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Update User Successfully", "result": result})
	}
}

func (uh *UserHandler) AddPhoneNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		email := ctx.MustGet("email").(string)
		var userInput dto.AddPhone
		if err := ctx.BindJSON(&userInput); err != nil {
			log.Println("Error Input req user:", err.Error())
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := uh.userSvc.AddPhoneNumber(ctx, email, &userInput)
		if err != nil {
			msg := "Error updating user" + err.Error()
			util.HandleError(ctx, err, http.StatusInternalServerError, msg)
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Add Phone Number Successfully"})
	}
}
