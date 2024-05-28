package service

import (
	"context"
	"errors"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type midtransService struct {
	config          config.Midtrans
	envi            midtrans.EnvironmentType
	paymentRepo     domain.PaymentRepository
	orderRepo       domain.OrderRepository
	sellerOrderRepo domain.SellerOrderRepository
}

func NewMidtransService(cnf *config.Config, paymentRepo domain.PaymentRepository,
	orderRepo domain.OrderRepository, sellerOrderRepo domain.SellerOrderRepository) domain.MidtransService {
	envi := midtrans.Sandbox
	if cnf.Midtrans.IsProd {
		envi = midtrans.Production
	}

	return &midtransService{
		config:          cnf.Midtrans,
		envi:            envi,
		paymentRepo:     paymentRepo,
		orderRepo:       orderRepo,
		sellerOrderRepo: sellerOrderRepo,
	}
}

// GenerateSnapURL implements domain.MidtransService.
func (s *midtransService) GenerateSnapURL(ctx context.Context, p *domain.Payment) error {
	// 2. Initiate Snap request
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  p.OrderID,
			GrossAmt: int64(p.Amount),
		},
	}

	var client snap.Client
	client.New(s.config.Key, s.envi)
	// 3. Request create Snap transaction to Midtrans
	snapResp, err := client.CreateTransaction(req)
	if err != nil {
		return errors.New("failed create transaction: " + err.Error())
	}
	p.Snap_Url = snapResp.RedirectURL
	return nil
}

// VerifyPayment implements domain.MidtransService.
func (s *midtransService) VerifyPayment(ctx context.Context, orderID string) (bool, error) {
	var client coreapi.Client
	client.New(s.config.Key, s.envi)
	// 4. Check transaction to Midtrans with param orderId
	transactionStatusResp, e := client.CheckTransaction(orderID)

	if e != nil {
		return false, errors.New("failed check transaction : " + e.Error())
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					// TODO set transaction status on your database to 'success'
					return true, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				// TODO set transaction status on your databaase to 'success'
				var req dto.UpdatePaymentReq
				var reqSO dto.OrderSellerUpdateReq
				var reqStatus dto.OrderStatusUpdateReq
				transactionID := transactionStatusResp.TransactionID

				req.Payment_Method = transactionStatusResp.PaymentType
				req.Status = "SUCCESS"
				req.TransactionID = transactionID
				reqSO.Payment_Status = req.Status
				reqSO.Status = "PROCESSED"
				reqStatus.Status = reqSO.Status

				payment, err := s.paymentRepo.FindByOrderId(ctx, orderID)
				if err != nil {
					return false, err
				}

				order, err := s.orderRepo.GetOrder(ctx, orderID)
				if err != nil {
					return false, err
				}

				payment.UpdateAt = time.Now()
				order.Updated_At = time.Now()

				_, err = s.paymentRepo.Update(ctx, orderID, &req)
				if err != nil {
					return false, errors.New("Failed to update payment: " + err.Error())
				}

				_, err = s.orderRepo.UpdateOrder(ctx, orderID, &req)
				if err != nil {
					return false, errors.New("Failed to update order: " + err.Error())
				}

				_, err = s.sellerOrderRepo.UpdateOrderSeller(ctx, orderID, &reqSO)
				if err != nil {
					return false, errors.New("Failed to update seller order: " + err.Error())
				}

				for _, item := range order.Items {
					_, err = s.orderRepo.UpdateStatusOrder(ctx, orderID, item.Product_Id, &reqStatus)
					if err != nil {
						return false, errors.New("Failed to update status order: " + err.Error())
					}
				}
				return true, nil
			}
		}
	}
	return false, nil
}
