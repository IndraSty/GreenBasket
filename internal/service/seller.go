package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	dto "github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type sellerService struct {
	repo      domain.SellerRepository
	tokenSvc  domain.TokenService
	cacheRepo domain.CacheRepository
	emailSvc  domain.EmailService
}

func NewSellerService(repo domain.SellerRepository, tokenSvc domain.TokenService,
	cacheRepo domain.CacheRepository,
	emailSvc domain.EmailService) domain.SellerService {
	return &sellerService{
		repo:      repo,
		tokenSvc:  tokenSvc,
		cacheRepo: cacheRepo,
		emailSvc:  emailSvc,
	}
}

func (s *sellerService) RegisterSeller(ctx context.Context, req *dto.SellerRegisterReq) (*dto.SellerRegisterRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	passwordErr := util.ValidatePassword(req.Password)
	if passwordErr != "" {
		return nil, errors.New(passwordErr)
	}

	emailExist, err := s.repo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if emailExist {
		return nil, errors.New("email already registered")
	}

	phoneExists, err := s.repo.CheckPhoneExists(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if phoneExists {
		return nil, errors.New("phone number already registered")
	}

	password := util.HashPassword(req.Password)

	id := primitive.NewObjectID()
	sellerId := id.Hex()

	seller := domain.Seller{
		ID:            id,
		First_Name:    req.First_Name,
		Last_Name:     req.Last_Name,
		Email:         req.Email,
		Password:      password,
		Phone:         req.Phone,
		Role:          "Seller",
		Created_At:    time.Now(),
		Updated_At:    time.Now(),
		EmailVerified: false,
		PhoneVerified: false,
		Seller_Id:     sellerId,
	}

	var msg []string

	otpCode := util.GenarateRandomNumber(4)
	err = s.emailSvc.SendMail(req.Email, "OTP Code", "otp anda "+otpCode)
	if err != nil {
		return nil, errors.New("failed to send email: " + err.Error())
	}

	exp := 15 * time.Minute
	err = s.cacheRepo.Set("otp:"+sellerId, []byte(otpCode), exp)
	if err != nil {
		return nil, errors.New("failed to add otp to redis :" + err.Error())
	}
	err = s.cacheRepo.Set("seller-ref:"+sellerId, []byte(seller.Email), exp)
	if err != nil {
		return nil, errors.New("failed to add seller email to redis :" + err.Error())
	}

	insertResult, err := s.repo.CreateSeller(ctx, seller)
	if err != nil {
		return nil, err
	}

	msg = append(msg, "OTP Has been send on Your email")
	msg = append(msg, "Seller created successfully!")

	return &dto.SellerRegisterRes{
		InsertId: insertResult,
		Message:  msg,
	}, nil
}

func (s *sellerService) AuthenticateSeller(ctx context.Context, req *dto.SellerAuthReq) (*dto.SellerAuthRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	passwordErr := util.ValidatePassword(req.Password)
	if passwordErr != "" {
		return nil, errors.New(passwordErr)
	}

	seller, err := s.repo.FindSellerByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("email doesn't exist")
	}

	passwordIsInvalid, _ := util.VerifyPassword(req.Password, seller.Password)
	if !passwordIsInvalid {
		return nil, errors.New("email or password incorrect")
	}

	acc_token, refreshToken, _ := s.tokenSvc.GenerateAllTokens(seller.Email, seller.First_Name, seller.Last_Name, seller.Seller_Id)

	return &dto.SellerAuthRes{
		Access_Token:  acc_token,
		Refresh_Token: refreshToken,
	}, nil
}

// ValidateOTP implements domain.UserService.
func (s *sellerService) ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error {
	var userReq dto.UserUpdateReq
	val, err := s.cacheRepo.Get("otp:" + req.UserID)
	if err != nil {
		return errors.New("otp code is invalid")
	}

	otp := string(val)
	if otp != req.OTP {
		return errors.New("otp code is invalid")
	}

	val, err = s.cacheRepo.Get("seller-ref:" + req.UserID)
	if err != nil {
		return errors.New("otp code is invalid")
	}

	_, err = s.repo.FindSellerByEmail(ctx, string(val))
	if err != nil {
		return errors.New("failed to get seller by email: " + err.Error())
	}

	var updateUser primitive.D
	userReq.EmailVerified = true
	updateUser = append(updateUser, bson.E{Key: "email_verified", Value: userReq.EmailVerified})

	_, err = s.repo.UpdateSeller(ctx, string(val), updateUser)
	if err != nil {
		return errors.New("failed to update email verified: " + err.Error())
	}

	return nil
}

func (s *sellerService) GetSellerByEmail(ctx context.Context, email string) (*domain.Seller, error) {
	val, err := s.cacheRepo.Get("seller:" + email)
	if err == nil {
		var seller domain.Seller
		err = json.Unmarshal(val, &seller)
		if err != nil {
			return nil, errors.New("failed to unmarshal seller data: " + err.Error())
		}
		return &seller, nil
	}

	// Data not found in cache, fetch from database
	seller, err := s.repo.FindSellerByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get seller by email: " + err.Error())
	}

	sellerData, err := json.Marshal(seller)
	if err != nil {
		return nil, errors.New("failed to marshal seller data: " + err.Error())
	}

	err = s.cacheRepo.Set("seller:"+email, sellerData, time.Hour*24)
	if err != nil {
		return nil, errors.New("failed to set seller data in cache: " + err.Error())
	}

	return seller, nil
}

func (s *sellerService) UpdateSeller(ctx context.Context, email string, req *dto.SellerUpdateReq) (*dto.SellerUpdateRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	seller, err := s.repo.FindSellerByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("seller not found!" + err.Error())
	}

	var update primitive.D
	if req.First_Name != "" {
		update = append(update, bson.E{Key: "first_name", Value: req.First_Name})
	}
	if req.Last_Name != "" {
		update = append(update, bson.E{Key: "last_name", Value: req.Last_Name})
	}

	if req.Image_Url != "" {
		update = append(update, bson.E{Key: "image_url", Value: req.Image_Url})
	}

	seller.Updated_At = time.Now()

	result, err := s.repo.UpdateSeller(ctx, email, update)
	if err != nil {
		return nil, errors.New("Failed to update user: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no data was updated")
	}

	_, err = s.GetSellerByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("seller not found! and failed to set redis" + err.Error())
	}

	return &dto.SellerUpdateRes{
		UpdateResult: result,
	}, nil

}
