package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type orderService struct {
	repo            domain.OrderRepository
	userRepo        domain.UserRepository
	cartRepo        domain.CartRepository
	sellerRepo      domain.SellerRepository
	storeRepo       domain.StoreRepository
	notifSvc        domain.NotificationService
	sellerOrderRepo domain.SellerOrderRepository
	salesReportSvc  domain.SalesReportService
	cacheRepo       domain.CacheRepository
}

func NewOrderService(repo domain.OrderRepository, userRepo domain.UserRepository, cartRepo domain.CartRepository,
	sellerRepo domain.SellerRepository, storeRepo domain.StoreRepository, notifSvc domain.NotificationService,
	sellerOrderRepo domain.SellerOrderRepository, salesReportSvc domain.SalesReportService,
	cacheRepo domain.CacheRepository) domain.OrderService {
	return &orderService{
		repo:            repo,
		userRepo:        userRepo,
		cartRepo:        cartRepo,
		sellerRepo:      sellerRepo,
		storeRepo:       storeRepo,
		notifSvc:        notifSvc,
		sellerOrderRepo: sellerOrderRepo,
		salesReportSvc:  salesReportSvc,
		cacheRepo:       cacheRepo,
	}
}

func (s *orderService) setRedisOrder(item any, name, email string) error {
	itemData, err := json.Marshal(item)
	if err != nil {
		return errors.New("failed to marshal order data: " + err.Error())
	}

	err = s.cacheRepo.Set(name+email, itemData, time.Hour*1)
	if err != nil {
		return errors.New("failed to set user order data in cache: " + err.Error())
	}

	return nil
}

func (s *orderService) delRedisOrder(email string, name ...string) error {
	err := s.cacheRepo.Del(name[0] + email)
	if err != nil {
		return errors.New("failed to delete order data in cache: " + err.Error())
	}

	if len(name[1]) != 0 {
		err = s.cacheRepo.Del(name[1] + email)
		if err != nil {
			return errors.New("failed to delete all order data in cache: " + err.Error())
		}
	}

	return nil
}

func (s *orderService) updateRedisOrder(ctx context.Context, email, orderID string, name ...string) error {
	data1, err := s.getOrderByEmailAndId(ctx, email, orderID)
	if err != nil {
		return errors.New("failed to get order data: " + err.Error())
	}

	err = s.setRedisOrder(*data1, name[0], email)
	if err != nil {
		return errors.New("failed to set user order data1 in cache: " + err.Error())
	}

	if len(name[1]) != 0 {
		data2, err := s.getAllOrders(ctx, email)
		if err != nil {
			return errors.New("failed to get order data: " + err.Error())
		}
		err = s.setRedisOrder(*data2, name[1], email)
		if err != nil {
			return errors.New("failed to set user order data2 in cache: " + err.Error())
		}
	}

	return nil
}

func (s *orderService) getAllOrders(ctx context.Context, email string) (*[]domain.Orders, error) {
	_, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to find user: " + err.Error())
	}

	result, err := s.repo.GetAllOrders(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get all orders: " + err.Error())
	}

	return result, nil
}

func (s *orderService) getOrderByEmailAndId(ctx context.Context, email string, orderID string) (*domain.Orders, error) {
	_, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to find user: " + err.Error())
	}

	result, err := s.repo.GetOrder(ctx, orderID, email)
	if err != nil {
		return nil, errors.New("failed to get the order: " + err.Error())
	}

	return result, nil
}

// CreateOrder implements domain.OrderService.
func (s *orderService) CreateOrder(ctx context.Context, email string) (*dto.InsertOrderRes, error) {
	err := s.delRedisOrder(email, "user-order:", "all_user-order:")
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to find user: " + err.Error())
	}

	if user.Address_Details == nil {
		return nil, errors.New("user doesn't have an addresses")
	}

	cart, err := s.cartRepo.GetUserCart(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get user cart: " + err.Error())
	}

	if len(cart.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	var items []domain.OrderItem
	var totalPrice float64

	for _, item := range cart.Items {
		if item.Selected {
			orderItem := domain.OrderItem{
				Product_Id:    item.Product_Id,
				Product_Name:  item.Product_Name,
				Product_Image: item.Product_Image,
				StoreID:       item.StoreID,
				Order_Status:  "PENDING",
				Quantity:      item.Quantity,
				Price:         item.Price,
			}
			items = append(items, orderItem)
			totalPrice += item.Price * float64(item.Quantity)
		}
	}

	if len(items) == 0 {
		return nil, errors.New("no items selected in the cart")
	}

	id := primitive.NewObjectID()
	orderID := id.Hex()
	order := domain.Orders{
		ID:               id,
		Order_id:         orderID,
		Email:            email,
		Order_Date:       time.Now(),
		Updated_At:       time.Now(),
		Total_Price:      totalPrice,
		Address_Shipping: *user.Address_Details,
		Payment:          &domain.PaymentOrder{},
		Items:            items,
	}

	result, err := s.repo.CreateOrder(ctx, order)
	if err != nil {
		return nil, errors.New("failed to create an order: " + err.Error())
	}

	sellerItems := make(map[string][]domain.OrderItem)

	for _, item := range items {
		sellerItems[item.StoreID] = append(sellerItems[item.StoreID], item)
	}

	for _, items := range sellerItems {
		var sellerOrderItems []domain.SellerOrderItem
		var totalPriceSeller float64
		var emailSeller string

		for _, item := range items {
			sellerOrderItem := domain.SellerOrderItem{
				User_Email:       email,
				Product_Id:       item.Product_Id,
				Product_Name:     item.Product_Name,
				Product_Image:    item.Product_Image,
				Quantity:         item.Quantity,
				Price:            item.Price,
				Status:           "PENDING",
				Address_Shipping: *user.Address_Details,
			}

			seller, err := s.sellerRepo.FindSellerByStoreId(ctx, item.StoreID)
			if err != nil {
				return nil, errors.New("failed to find seller for add order: " + err.Error())
			}
			emailSeller = seller.Email
			sellerOrderItems = append(sellerOrderItems, sellerOrderItem)
			totalPriceSeller += item.Price * float64(item.Quantity)
		}

		sellerOrder := domain.SellerOrder{
			ID:             primitive.NewObjectID(),
			Order_id:       orderID,
			Email:          emailSeller,
			Ordered_At:     time.Now(),
			Updated_At:     time.Now(),
			Total_Price:    totalPriceSeller,
			Payment_Status: "UNPAID",
			Items:          sellerOrderItems,
		}

		_, err := s.sellerOrderRepo.CreateOrderSeller(ctx, sellerOrder)
		if err != nil {
			return nil, errors.New("failed to create a seller order: " + err.Error())
		}

		err = s.delRedisOrder(emailSeller, "seller-order:", "all_seller-order:")
		if err != nil {
			log.Println("failed to update seller order in cache: ", err)
		}
	}

	storeIDs := make(map[string]bool)
	for _, item := range items {
		storeIDs[item.StoreID] = true
	}

	// Send a notification for each unique seller ID
	for storeID := range storeIDs {
		store, err := s.storeRepo.GetStore(ctx, storeID)
		if err != nil {
			return nil, errors.New("failed to find store for notification: " + err.Error())
		}

		seller, err := s.sellerRepo.FindSellerByEmail(ctx, store.Email)
		if err != nil {
			return nil, errors.New("failed to find seller for notification: " + err.Error())
		}
		if seller == nil {
			return nil, errors.New("seller not found for store ID: " + storeID)
		}
		go s.notificationAfterOrder(seller.Email, user.Email, orderID)

	}

	defer func() {
		err := s.updateRedisOrder(ctx, email, orderID, "user-order:", "all_seller-order")
		if err != nil {
			log.Println("failed to update order in cache: ", err)
		}
	}()

	return &dto.InsertOrderRes{
		InsertId: result,
	}, nil
}

// GetAllOrders implements domain.OrderService.
func (s *orderService) GetAllOrders(ctx context.Context, email string) (*[]domain.Orders, error) {
	val, err := s.cacheRepo.Get("all_user-order:" + email)
	if err == nil {
		var data []domain.Orders
		err = json.Unmarshal(val, &data)
		if err != nil {
			return nil, errors.New("failed to unmarshal order data: " + err.Error())
		}

		return &data, nil
	}

	_, err = s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to find user: " + err.Error())
	}

	result, err := s.repo.GetAllOrders(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get all orders: " + err.Error())
	}

	err = s.setRedisOrder(result, "all_user-order:", email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetOrderById implements domain.OrderService.
func (s *orderService) GetOrderByEmailAndId(ctx context.Context, email string, orderID string) (*domain.Orders, error) {
	val, err := s.cacheRepo.Get("user-order:" + email)
	if err == nil {
		var data domain.Orders
		err = json.Unmarshal(val, &data)
		if err != nil {
			return nil, errors.New("failed to unmarshal order data: " + err.Error())
		}

		return &data, nil
	}

	_, err = s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("failed to find user: " + err.Error())
	}

	result, err := s.repo.GetOrder(ctx, orderID, email)
	if err != nil {
		return nil, errors.New("failed to get the order: " + err.Error())
	}

	err = s.setRedisOrder(result, "user-order:", email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FinishOrder implements domain.OrderService.
func (s *orderService) FinishOrder(ctx context.Context, email string, orderID string, productID string, req *dto.OrderStatusUpdateReq) error {
	err := s.delRedisOrder(email, "user-order:", "all_user-order:")
	if err != nil {
		return err
	}

	order, err := s.repo.GetOrder(ctx, orderID, email)
	if err != nil {
		return errors.New("failed to get the order: " + err.Error())
	}

	for _, item := range order.Items {
		if item.Product_Id == productID {
			if item.Order_Status != "SHIPPED" {
				return errors.New("order status is not 'SHIPPED'")
			}
		}
	}

	if req.Status != "FINISHED" {
		return errors.New("order status is not 'FINISHED'")
	}

	res, err := s.repo.UpdateStatusOrder(ctx, orderID, productID, req)
	if err != nil {
		return errors.New("failed to update status order: " + err.Error())
	}

	if res.ModifiedCount == 0 {
		return errors.New("failed to update status order")
	}

	res, err = s.sellerOrderRepo.UpdateStatusOrderSeller(ctx, orderID, productID, req)
	if err != nil {
		return errors.New("failed to update status order seller: " + err.Error())
	}

	if res.ModifiedCount == 0 {
		return errors.New("failed to update status order seller")
	}

	newOrder, err := s.repo.GetOrder(ctx, orderID, email)
	if err != nil {
		return errors.New("failed to get seller order: " + err.Error())
	}
	var sellerID string
	for _, item := range newOrder.Items {
		if item.Product_Id == productID {
			if item.Order_Status == "FINISHED" {
				seller, err := s.sellerRepo.FindSellerByStoreId(ctx, item.StoreID)
				if err != nil {
					return errors.New("failed to get seller with email: " + err.Error())
				}

				sellerID = seller.Seller_Id
				err = s.salesReportSvc.UpdateSalesReport(ctx, seller.Store_Id, seller.Email)
				if err != nil {
					return err
				}
			}
		}
	}

	defer func() {
		err := s.updateRedisOrder(ctx, email, orderID, "user-order:", "all_seller-order")
		if err != nil {
			log.Println("failed to update order in cache: ", err)
		}
	}()

	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return errors.New("failed to get user with email: " + err.Error())
	}
	go s.notificationFinishOrder(sellerID, orderID, productID, user.First_Name)

	return nil
}

// CancelOrder implements domain.OrderService.
func (s *orderService) CancelOrder(ctx context.Context, email string, orderID string, productID string) error {
	err := s.delRedisOrder(email, "user-order:", "all_user-order:")
	if err != nil {
		return err
	}

	order, err := s.repo.GetOrder(ctx, orderID, email)
	if err != nil {
		return errors.New("failed to get the order: " + err.Error())
	}

	for _, item := range order.Items {
		if item.Product_Id == productID {
			if item.Order_Status != "PENDING" {
				return errors.New("order status is not 'PENDING'")
			}
		}
	}

	res, err := s.repo.DeleteItem(ctx, orderID, productID)
	if err != nil {
		return errors.New("failed to delete item: " + err.Error())
	}

	if res.ModifiedCount == 0 {
		return errors.New("failed to delete item")
	}

	for _, item := range order.Items {
		if item.Product_Id == productID {
			seller, err := s.sellerRepo.FindSellerByStoreId(ctx, item.StoreID)
			if err != nil {
				return errors.New("failed to find seller: " + err.Error())
			}

			res, err := s.sellerOrderRepo.DeleteItem(ctx, seller.Email, orderID, item.Product_Id)
			if err != nil {
				return errors.New("failed to delete item: " + err.Error())
			}

			if res.ModifiedCount == 0 {
				return errors.New("failed to delete item")
			}
		}
	}

	defer func() {
		err := s.updateRedisOrder(ctx, email, orderID, "user-order:", "all_seller-order")
		if err != nil {
			log.Println("failed to update order in cache: ", err)
		}
	}()

	return nil
}

func (s *orderService) notificationAfterOrder(sellerEmail string, email string, orderID string) error {
	data := map[string]string{
		"order_id": orderID,
	}
	err := s.notifSvc.Insert(context.Background(), email, "USER_ORDER", data)
	if err != nil {
		return errors.New("failed to insert user notification :" + err.Error())
	}

	err = s.notifSvc.Insert(context.Background(), sellerEmail, "SELLER_ORDER", data)
	if err != nil {
		return errors.New("failed to insert seller notification :" + err.Error())
	}

	return nil
}

func (s *orderService) notificationFinishOrder(sellerEmail, orderID, productID, username string) error {
	data := map[string]string{
		"order_id":   orderID,
		"product_id": productID,
		"username":   username,
	}
	err := s.notifSvc.Insert(context.Background(), sellerEmail, "SELLER_FINISH_ORDER", data)
	if err != nil {
		return errors.New("failed to insert user notification :" + err.Error())
	}

	return nil
}
