package service

import (
	"context"
	"errors"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type passwordService struct {
	userRepo domain.UserRepository
}

func NewPasswordService(userRepo domain.UserRepository) domain.PasswordService {
	return &passwordService{
		userRepo: userRepo,
	}
}

// ChangePassword implements domain.PasswordService.
func (s *passwordService) ChangePassword(ctx context.Context, email string, req dto.PasswordReq) error {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return errors.New("Invalid request body" + err.Error())
	}
	if req.NewPassword != req.RewritePassword {
		return errors.New("password doesn't match")
	}
	passwordErr := util.ValidatePassword(req.NewPassword)
	if passwordErr != "" {
		return errors.New(passwordErr)
	}
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return errors.New("failed to get user data: " + err.Error())
	}
	passwordIsInvalid, _ := util.VerifyPassword(req.NewPassword, user.Password)
	if !passwordIsInvalid {
		return errors.New("password doesn't match")
	}
	var update primitive.D
	updateAt := time.Now()
	password := util.HashPassword(req.NewPassword)
	update = append(update, bson.E{Key: "password", Value: password})
	update = append(update, bson.E{Key: "updated_at", Value: updateAt})
	res, err := s.userRepo.UpdateUser(ctx, email, update)
	if err != nil {
		return errors.New("failed update password" + err.Error())
	}

	if res.ModifiedCount == 0 {
		return errors.New("no password was updated")
	}

	return nil
}

// RecoveryPassword implements domain.PasswordService.
func (s *passwordService) RecoveryPassword(ctx context.Context, email string, req dto.PasswordReq) error {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return errors.New("Invalid request body" + err.Error())
	}
	if req.NewPassword != req.RewritePassword {
		return errors.New("password doesn't match")
	}
	var update primitive.D
	updateAt := time.Now()
	password := util.HashPassword(req.NewPassword)
	update = append(update, bson.E{Key: "password", Value: password})
	update = append(update, bson.E{Key: "updated_at", Value: updateAt})
	res, err := s.userRepo.UpdateUser(ctx, email, update)
	if err != nil {
		return errors.New("failed update password" + err.Error())
	}

	if res.ModifiedCount == 0 {
		return errors.New("no password was updated")
	}

	return nil
}
