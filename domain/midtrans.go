package domain

import "context"

type MidtransService interface {
	GenerateSnapURL(ctx context.Context, p *Payment) error
	VerifyPayment(ctx context.Context, orderID string) (bool, error)
}
