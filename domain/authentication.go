package domain

import (
	"context"

	"github.com/IndraSty/GreenBasket/dto"
)

type AuthService interface {
	ValidateOTP(ctx context.Context, req dto.ValidateOtpReq) error
	AuthenticateUser(ctx context.Context, req *dto.UserAuthReq) (*dto.UserAuthRes, error)
	RequestEmail(ctx context.Context, req dto.UserReqEmail, action string) error
}
