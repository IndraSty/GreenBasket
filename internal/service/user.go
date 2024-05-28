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

type userService struct {
	repo      domain.UserRepository
	cacheRepo domain.CacheRepository
	emailSvc  domain.EmailService
	cartSvc   domain.CartService
}

func NewUserService(repo domain.UserRepository,
	emailSvc domain.EmailService,
	cacheRepo domain.CacheRepository, cartSvc domain.CartService) domain.UserService {
	return &userService{
		repo:      repo,
		emailSvc:  emailSvc,
		cacheRepo: cacheRepo,
		cartSvc:   cartSvc,
	}
}

func (us *userService) setRedisUser(user domain.User, email string) error {
	userData, err := json.Marshal(user)
	if err != nil {
		return errors.New("failed to marshal user data: " + err.Error())
	}

	err = us.cacheRepo.Set("user:"+email, userData, time.Hour*24)
	if err != nil {
		return errors.New("failed to set user data in cache: " + err.Error())
	}
	return nil
}

func (us *userService) RegisterUser(ctx context.Context, req *dto.UserRegisterReq) (*dto.UserRegisterRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body:" + err.Error())
	}

	passwordErr := util.ValidatePassword(req.Password)
	if passwordErr != "" {
		return nil, errors.New(passwordErr)
	}

	emailExist, err := us.repo.CheckEmailExists(ctx, req.Email)
	if err != nil {
		return nil, errors.New("failed to check email: " + err.Error())
	}

	if emailExist {
		return nil, errors.New("email already registered")
	}

	phoneExists, err := us.repo.CheckPhoneExists(ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if phoneExists {
		return nil, errors.New("phone number already registered")
	}

	password := util.HashPassword(req.Password)

	id := primitive.NewObjectID()
	userId := id.Hex()

	user := domain.User{
		ID:            id,
		First_Name:    req.First_Name,
		Last_Name:     req.Last_Name,
		Email:         req.Email,
		Password:      password,
		Phone:         req.Phone,
		Role:          "User",
		Created_At:    time.Now(),
		Updated_At:    time.Now(),
		EmailVerified: false,
		PhoneVerified: false,
		User_Id:       userId,
	}

	otpCode := util.GenarateRandomNumber(4)

	err = us.emailSvc.SendMail(req.Email, "OTP Code", "otp anda "+otpCode)
	if err != nil {
		return nil, errors.New("failed to send email: " + err.Error())
	}

	insertResult, err := us.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, errors.New("failed to create user: " + err.Error())
	}

	if err = us.cartSvc.CreateCart(ctx, user.Email); err != nil {
		return nil, errors.New("failed to create user cart: " + err.Error())
	}

	exp := 15 * time.Minute
	err = us.cacheRepo.Set("otp:"+userId, []byte(otpCode), exp)
	if err != nil {
		return nil, errors.New("failed to add otp to redis :" + err.Error())
	}
	err = us.cacheRepo.Set("user-id:"+userId, []byte(user.Email), exp)
	if err != nil {
		return nil, errors.New("failed to add user email to redis :" + err.Error())
	}

	return &dto.UserRegisterRes{
		InsertId: insertResult,
	}, nil
}

func (us *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	val, err := us.cacheRepo.Get("user:" + email)
	if err == nil {
		var user domain.User
		err = json.Unmarshal(val, &user)
		if err != nil {
			return nil, errors.New("failed to unmarshal user data: " + err.Error())
		}
		return &user, nil
	}

	// Data not found in cache, fetch from database
	user, err := us.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get user by email: " + err.Error())
	}

	err = us.setRedisUser(*user, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *userService) UpdateUser(ctx context.Context, email string, req *dto.UserUpdateReq) (*dto.UserUpdateRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	user, err := us.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	var updateUser primitive.D
	if req.First_Name != "" {
		updateUser = append(updateUser, bson.E{Key: "first_name", Value: req.First_Name})
	}
	if req.Last_Name != "" {
		updateUser = append(updateUser, bson.E{Key: "last_name", Value: req.Last_Name})
	}

	if req.Image_Url != "" {
		updateUser = append(updateUser, bson.E{Key: "image_url", Value: req.Image_Url})
	}

	user.Updated_At = time.Now()

	result, err := us.repo.UpdateUser(ctx, email, updateUser)
	if err != nil {
		return nil, errors.New("Failed to update user: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no data was updated")
	}

	user, err = us.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	err = us.setRedisUser(*user, email)
	if err != nil {
		return nil, err
	}

	return &dto.UserUpdateRes{
		UpdateResult: result,
	}, nil
}

// AddPhoneNumber implements domain.UserService.
func (us *userService) AddPhoneNumber(ctx context.Context, email string, req *dto.AddPhone) error {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return errors.New("Invalid request body" + err.Error())
	}

	user, err := us.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	var updateUser primitive.D
	if req.PhoneNumber != "" {
		updateUser = append(updateUser, bson.E{Key: "phone", Value: req.PhoneNumber})
	}

	user.Updated_At = time.Now()

	result, err := us.repo.UpdateUser(ctx, email, updateUser)
	if err != nil {
		return errors.New("Failed to update user: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return errors.New("no data was updated")
	}

	user, err = us.repo.FindUserByEmail(ctx, email)
	if err != nil {
		return errors.New("user not found")
	}

	err = us.setRedisUser(*user, email)
	if err != nil {
		return err
	}

	return nil
}
