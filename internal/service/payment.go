package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type paymentService struct {
	notifSvc    domain.NotificationService
	midtransSvc domain.MidtransService
	repo        domain.PaymentRepository
	userRepo    domain.UserRepository
}

func NewPaymentService(notifSvc domain.NotificationService, repo domain.PaymentRepository, userRepo domain.UserRepository, midtransSvc domain.MidtransService) domain.PaymentService {
	return &paymentService{
		notifSvc:    notifSvc,
		repo:        repo,
		midtransSvc: midtransSvc,
		userRepo:    userRepo,
	}
}

// InitializePayment implements domain.PaymentService.
func (s *paymentService) InitializePayment(ctx context.Context, req *dto.PaymentReq) (*dto.PaymentRes, error) {
	payment := domain.Payment{
		ID:        primitive.NewObjectID(),
		OrderID:   req.OrderID,
		UserID:    req.UserID,
		CreatedAt: time.Now(),
		UpdateAt:  time.Now(),
		Status:    "PENDING",
		Amount:    req.Amount,
	}

	err := s.midtransSvc.GenerateSnapURL(ctx, &payment)
	if err != nil {
		return &dto.PaymentRes{}, err
	}

	if err := s.repo.Insert(ctx, &payment); err != nil {
		return &dto.PaymentRes{}, err
	}

	return &dto.PaymentRes{
		Snap_Url: payment.Snap_Url,
	}, nil
}

// ConfirmedPayment implements domain.PaymentService.
func (s *paymentService) ConfirmedPayment(ctx context.Context, orderId string) error {
	payment, err := s.repo.FindByOrderId(ctx, orderId)
	if err != nil {
		return err
	}

	if payment == (&domain.Payment{}) {
		return errors.New("payment request not found")
	}

	data := map[string]string{
		"order_id": payment.OrderID,
		"amount":   fmt.Sprintf("%.2f", payment.Amount),
	}
	err = s.notifSvc.Insert(ctx, payment.UserID, "USER_PAYMENT", data)
	if err != nil {
		return errors.New("failed to insert user payment notification: " + err.Error())
	}

	return nil
}

// UpdatePayment implements domain.PaymentService.
func (s *paymentService) UpdatePayment(ctx context.Context, orderID string, req *dto.UpdatePaymentReq) error {
	payment, err := s.repo.FindByOrderId(ctx, orderID)
	if err != nil {
		return nil
	}

	payment.UpdateAt = time.Now()

	_, err = s.repo.Update(ctx, orderID, req)
	if err != nil {
		return errors.New("Failed to update payment: " + err.Error())
	}

	return nil
}
