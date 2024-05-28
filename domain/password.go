package domain

import (
	"context"

	"github.com/IndraSty/GreenBasket/dto"
)

type PasswordService interface {
	ChangePassword(ctx context.Context, email string, req dto.PasswordReq) error
	RecoveryPassword(ctx context.Context, email string, req dto.PasswordReq) error
}
