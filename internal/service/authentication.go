package service

import (
	"context"
	"errors"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type authService struct {
	userRepo  domain.UserRepository
	cacheRepo domain.CacheRepository
	tokenSvc  domain.TokenService
	emailSvc  domain.EmailService
}

func NewAuthService(userRepo domain.UserRepository, cacheRepo domain.CacheRepository,
	tokenSvc domain.TokenService, emailSvc domain.EmailService) domain.AuthService {
	return &authService{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
		tokenSvc:  tokenSvc,
		emailSvc:  emailSvc,
	}
}

func (s *authService) ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error {
	var userReq dto.UserUpdateReq
	val, err := s.cacheRepo.Get("otp:" + req.UserID)
	if err != nil {
		return errors.New("otp code is invalid")
	}

	otp := string(val)
	if otp != req.OTP {
		return errors.New("otp code is invalid")
	}

	val, err = s.cacheRepo.Get("user-id:" + req.UserID)
	if err != nil {
		return errors.New("otp code is invalid")
	}

	_, err = s.userRepo.FindUserByEmail(ctx, string(val))
	if err != nil {
		return errors.New("failed to get user by email: " + err.Error())
	}

	var updateUser primitive.D
	userReq.EmailVerified = true
	updateUser = append(updateUser, bson.E{Key: "email_verified", Value: userReq.EmailVerified})

	_, err = s.userRepo.UpdateUser(ctx, string(val), updateUser)
	if err != nil {
		return errors.New("failed to update email verified: " + err.Error())
	}

	return nil
}

func (s *authService) AuthenticateUser(ctx context.Context, req *dto.UserAuthReq) (*dto.UserAuthRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	passwordErr := util.ValidatePassword(req.Password)
	if passwordErr != "" {
		return nil, errors.New(passwordErr)
	}

	user, err := s.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("failed to find user by email: " + err.Error())
	}

	if !user.EmailVerified {
		return nil, errors.New("your email was not verified")
	}

	passwordIsInvalid, _ := util.VerifyPassword(req.Password, user.Password)
	if !passwordIsInvalid {
		return nil, errors.New("email or password incorrect")
	}

	acc_token, refreshToken, _ := s.tokenSvc.GenerateAllTokens(user.Email, user.First_Name, user.Last_Name, user.User_Id)

	return &dto.UserAuthRes{
		Access_Token:  acc_token,
		Refresh_Token: refreshToken,
	}, nil
}

func (s *authService) RequestEmail(ctx context.Context, req dto.UserReqEmail, action string) error {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return errors.New("Invalid request body: " + err.Error())
	}

	user, err := s.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return errors.New("failed to get user data: " + err.Error())
	}

	switch action {
	case "verify":
		if user.EmailVerified {
			return errors.New("email is already verified")
		}

		otpCode := util.GenarateRandomNumber(4)

		err = s.emailSvc.SendMail(req.Email, "OTP Code", "otp anda "+otpCode)
		if err != nil {
			return errors.New("failed to send email: " + err.Error())
		}

	case "recovery":
		if !user.EmailVerified {
			return errors.New("your email is not verified")
		}

	default:
		return errors.New("invalid action")
	}

	return nil
}
