package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
)

type sellerOrderService struct {
	repo        domain.SellerOrderRepository
	sellerRepo  domain.SellerRepository
	orderRepo   domain.OrderRepository
	productRepo domain.ProductRepository
	notifSvc    domain.NotificationService
	cacheRepo   domain.CacheRepository
}

func NewSellerOrderService(repo domain.SellerOrderRepository, sellerRepo domain.SellerRepository,
	orderRepo domain.OrderRepository, productRepo domain.ProductRepository,
	notifSvc domain.NotificationService, cacheRepo domain.CacheRepository) domain.SellerOrderService {
	return &sellerOrderService{
		repo:        repo,
		sellerRepo:  sellerRepo,
		orderRepo:   orderRepo,
		productRepo: productRepo,
		notifSvc:    notifSvc,
		cacheRepo:   cacheRepo,
	}
}

func (s *sellerOrderService) setRedisSO(item any, name, email string) error {
	itemData, err := json.Marshal(item)
	if err != nil {
		return errors.New("failed to marshal seller order data: " + err.Error())
	}

	err = s.cacheRepo.Set(name+email, itemData, time.Hour*1)
	if err != nil {
		return errors.New("failed to set user seller order data in cache: " + err.Error())
	}

	return nil
}

func (s *sellerOrderService) delRedisSO(email string, name ...string) error {
	err := s.cacheRepo.Del(name[0] + email)
	if err != nil {
		return errors.New("failed to delete seller order data in cache: " + err.Error())
	}

	if len(name[1]) != 0 {
		err = s.cacheRepo.Del(name[1] + email)
		if err != nil {
			return errors.New("failed to delete all seller order data in cache: " + err.Error())
		}
	}

	return nil
}

func (s *sellerOrderService) updateRedisSO(ctx context.Context, email, orderID string, name ...string) error {
	data1, err := s.getSellerOrdersWithNoAct(ctx, email, orderID)
	if err != nil {
		return errors.New("failed to get seller order data: " + err.Error())
	}

	err = s.setRedisSO(*data1, name[0], email)
	if err != nil {
		return errors.New("failed to set user seller order data1 in cache: " + err.Error())
	}

	if len(name[1]) != 0 {
		data2, err := s.getAllSellerOrdersWithNoAct(ctx, email)
		if err != nil {
			return errors.New("failed to get seller order data: " + err.Error())
		}
		err = s.setRedisSO(*data2, name[1], email)
		if err != nil {
			return errors.New("failed to set seller order data2 in cache: " + err.Error())
		}
	}

	return nil
}

func (s *sellerOrderService) getSellerOrdersWithNoAct(ctx context.Context, email string, orderID string) (*domain.SellerOrder, error) {
	_, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to find user: " + err.Error())
	}

	result, err := s.repo.GetSellerOrderByEmailAndId(ctx, email, orderID)
	if err != nil {
		return nil, errors.New("failed to get seller order: " + err.Error())
	}

	return result, nil
}

func (s *sellerOrderService) getAllSellerOrdersWithNoAct(ctx context.Context, email string) (*[]domain.SellerOrder, error) {
	_, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to find user: " + err.Error())
	}

	result, err := s.repo.GetAllSellerOrders(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get all orders: " + err.Error())
	}

	return result, nil
}

// GetAllSellerOrders implements domain.SellerOrderService.
func (s *sellerOrderService) GetAllSellerOrders(ctx context.Context, email string) (*[]domain.SellerOrder, error) {
	val, err := s.cacheRepo.Get("all_seller-order:" + email)
	if err == nil {
		var data []domain.SellerOrder
		err := json.Unmarshal(val, &data)
		if err != nil {
			return nil, errors.New("failed to unmarshal seller order data: " + err.Error())
		}
	}

	_, err = s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to find seller: " + err.Error())
	}

	result, err := s.repo.GetAllSellerOrders(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get all orders: " + err.Error())
	}

	err = s.setRedisSO(*result, "all_seller-order:", email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetSellerOrderByEmailAndId implements domain.SellerOrderService.
func (s *sellerOrderService) GetSellerOrderByEmailAndId(ctx context.Context, email string, orderID string) (*domain.SellerOrder, error) {
	val, err := s.cacheRepo.Get("seller-order:" + email)
	if err == nil {
		var data domain.SellerOrder
		err := json.Unmarshal(val, &data)
		if err != nil {
			return nil, errors.New("failed to unmarshal seller order data: " + err.Error())
		}
	}

	_, err = s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to find seller: " + err.Error())
	}

	result, err := s.repo.GetSellerOrderByEmailAndId(ctx, email, orderID)
	if err != nil {
		return nil, errors.New("failed to get the order: " + err.Error())
	}

	if result == nil {
		return nil, errors.New("no order found")
	}

	err = s.setRedisSO(*result, "seller-order:", email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateSellerAndUserOrder implements domain.SellerOrderService.
func (s *sellerOrderService) UpdateSellerAndUserOrderStatus(ctx context.Context, email, orderID, productID string, req *dto.OrderStatusUpdateReq) error {
	err := s.delRedisSO(email, "seller-order:", "all_seller-order:")
	if err != nil {
		return err
	}

	_, err = s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil {
		return errors.New("failed to find seller: " + err.Error())
	}

	sellerOrder, err := s.repo.GetSellerOrderByEmailAndId(ctx, email, orderID)
	if err != nil {
		return errors.New("failed to get the order: " + err.Error())
	}

	order, err := s.orderRepo.GetOrder(ctx, orderID)
	if err != nil {
		return errors.New("failed to get the user order: " + err.Error())
	}

	if sellerOrder.Payment_Status != "SUCCESS" {
		return errors.New("this order has not yet made payment")
	}

	_, err = s.repo.UpdateStatusOrderSeller(ctx, orderID, productID, req)
	if err != nil {
		return errors.New("failed to update the status seller order: " + err.Error())
	}

	_, err = s.orderRepo.UpdateStatusOrder(ctx, orderID, productID, req)
	if err != nil {
		return errors.New("failed to update the status user order: " + err.Error())
	}

	if req.Status == "SHIPPED" {
		for _, item := range order.Items {
			if item.Product_Id == productID {
				updateAT := time.Now()
				res, err := s.productRepo.UpdateStockProduct(ctx, item.StoreID, item.Product_Id, -item.Quantity, updateAT)
				if err != nil {
					return errors.New("failed to update stock product: " + err.Error())
				}

				go s.notificationProductShipped(order.Email, productID, item.StoreID)

				if res.ModifiedCount == 1 {
					product, err := s.productRepo.GetProductById(ctx, productID)
					if err != nil {
						return errors.New("failed to get product by id: " + err.Error())
					}
					if product.Stock <= 2 {
						go s.notificationStockProduct(email, productID)
					}
				}
			}
		}
	}

	defer func() {
		if err := s.updateRedisSO(ctx, email, orderID, "seller-order:", "all_seller-order:"); err != nil {
			log.Println("failed to update seller order in cache: ", err)
		}
	}()

	return nil
}

// CancelOrder implements domain.SellerOrderService.
func (s *sellerOrderService) CancelOrder(ctx context.Context, email string, orderID string, productID string) error {
	err := s.delRedisSO(email, "seller-order:", "all_seller-order:")
	if err != nil {
		return err
	}

	order, err := s.repo.GetSellerOrderByEmailAndId(ctx, email, orderID)
	if err != nil {
		return errors.New("failed to get the order: " + err.Error())
	}

	for _, item := range order.Items {
		if item.Product_Id == productID {
			if item.Status != "PENDING" {
				return errors.New("order status is not 'PENDING'")
			}
		}
	}

	res, err := s.repo.DeleteItem(ctx, email, orderID, productID)
	if err != nil {
		return errors.New("failed to delete item: " + err.Error())
	}

	if res.ModifiedCount == 0 {
		return errors.New("no item deleted")
	}

	for _, item := range order.Items {
		if item.Product_Id == productID {
			res, err := s.orderRepo.DeleteItem(ctx, orderID, productID)
			if err != nil {
				return errors.New("failed to delete item: " + err.Error())
			}

			if res.ModifiedCount == 0 {
				return errors.New("no item deleted")
			}
		}
	}

	defer func() {
		if err := s.updateRedisSO(ctx, email, orderID, "seller-order:", "all_seller-order"); err != nil {
			log.Println("failed to update seller order in cache: ", err)
		}
	}()

	return nil
}

func (s *sellerOrderService) notificationProductShipped(userID, productID, storeID string) error {
	data := map[string]string{
		"product_id": productID,
		"store_id":   storeID,
	}
	err := s.notifSvc.Insert(context.Background(), userID, "USER_PRODUCT_SHIPPED", data)
	if err != nil {
		return errors.New("failed to insert user notification :" + err.Error())
	}

	return nil
}

func (s *sellerOrderService) notificationStockProduct(sellerID, productID string) error {
	data := map[string]string{
		"product_id": productID,
	}
	err := s.notifSvc.Insert(context.Background(), sellerID, "SELLER_LESS_STOCK", data)
	if err != nil {
		return errors.New("failed to insert seller notification :" + err.Error())
	}

	return nil
}
