package delivery

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/config"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	userRepo domain.UserRepository
	auth     config.AuthSetup
	google   config.Google
	authSvc  domain.AuthService
}

func NewAuthHandler(userRepo domain.UserRepository, auth config.AuthSetup, cnf *config.Config, authSvc domain.AuthService) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
		auth:     auth,
		google:   cnf.Google,
		authSvc:  authSvc,
	}
}

func (h *AuthHandler) BeginAuthHandler(c *gin.Context) {
	provider := c.Param("provider")
	r := c.Request
	w := c.Writer

	r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.BeginAuthHandler(w, r)
}

func (h *AuthHandler) GetAuthCallBackFunc(c *gin.Context) {
	provider := c.Param("provider")
	r := c.Request
	w := c.Writer

	r = r.WithContext(context.WithValue(context.Background(), "provider", provider))
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Panic(err)
		return
	}

	emailExist, err := h.userRepo.CheckEmailExists(c, user.Email)
	if err != nil {
		log.Panic(err)
		return
	}

	var updateUser primitive.D
	if emailExist {
		if user.RefreshToken != "" {
			updateUser = append(updateUser, bson.E{Key: "refresh_token", Value: user.RefreshToken})
			_, err := h.userRepo.UpdateUser(c, user.Email, updateUser)
			if err != nil {
				log.Println("failed update refresh token", err.Error())
			}
		}

	}

	if !emailExist {
		id := primitive.NewObjectID()
		userID := id.Hex()

		if user.FirstName == "" && user.LastName == "" && user.Name != "" {
			nameParts := strings.Split(user.Name, " ")
			if len(nameParts) >= 2 {
				user.FirstName = nameParts[0]
				user.LastName = nameParts[1]
			}
		}

		fmt.Println("1", user.FirstName)

		input := domain.User{
			ID:            id,
			First_Name:    user.FirstName,
			Last_Name:     user.LastName,
			Email:         user.Email,
			Image_Url:     user.AvatarURL,
			Refresh_Token: user.RefreshToken,
			Role:          "User",
			Created_At:    time.Now(),
			Updated_At:    time.Now(),
			User_Id:       userID,
			Oauth_Id:      user.UserID,
			EmailVerified: true,
		}
		h.userRepo.CreateUser(c, input)
	}

	http.Redirect(w, r, "http://localhost:8080", http.StatusFound)
}

func (h *AuthHandler) LogoutHandler(c *gin.Context) {
	r := c.Request
	w := c.Writer

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	})
	gothic.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) ValidateOTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.ValidateOtpReq
		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		err := h.authSvc.ValidateOTP(ctx, req)
		if err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Validate OTP successfully!"})
	}
}

func (h *AuthHandler) AuthenticateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UserAuthReq
		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		res, err := h.authSvc.AuthenticateUser(ctx, &req)
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

func (h *AuthHandler) RequestVerifyEmail() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UserReqEmail
		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		if err := h.authSvc.RequestEmail(ctx, req, "verify"); err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "OTP code has been send on your'e email"})
	}
}

func (h *AuthHandler) RequestEmailForRecoveryPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.UserReqEmail
		if err := ctx.BindJSON(&req); err != nil {
			util.HandleError(ctx, err, http.StatusBadRequest, err.Error())
			return
		}

		if err := h.authSvc.RequestEmail(ctx, req, "recovery"); err != nil {
			util.HandleError(ctx, err, http.StatusInternalServerError, err.Error())
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully identity your email, please set new password"})
	}
}
