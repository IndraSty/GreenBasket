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

type cartService struct {
	repo        domain.CartRepository
	productRepo domain.ProductRepository
	storeRepo   domain.StoreRepository
	cacheRepo   domain.CacheRepository
}

func NewCartService(repo domain.CartRepository, productRepo domain.ProductRepository,
	storeRepo domain.StoreRepository, cacheRepo domain.CacheRepository) domain.CartService {
	return &cartService{
		repo:        repo,
		productRepo: productRepo,
		storeRepo:   storeRepo,
		cacheRepo:   cacheRepo,
	}
}

func (s *cartService) setRedisCart(item []dto.GetCartItemRes, email string) error {
	itemData, err := json.Marshal(item)
	if err != nil {
		return errors.New("failed to marshal user cart item: " + err.Error())
	}

	err = s.cacheRepo.Set("usercart-item:"+email, itemData, time.Hour*1)
	if err != nil {
		return errors.New("failed to set user cart item in cache: " + err.Error())
	}

	return nil
}

func (s *cartService) delRedisCartItem(email string) error {
	err := s.cacheRepo.Del("usercart-item:" + email)
	if err != nil {
		return errors.New("failed to delete cart item in cache: " + err.Error())
	}

	return nil
}

func (s *cartService) updateRedisCart(ctx context.Context, email string) error {
	cartItems, err := s.getAllCartItemWithNoAct(ctx, email)
	if err != nil {
		return errors.New("failed to get all item in user cart: " + err.Error())
	}

	return s.setRedisCart(*cartItems, email)
}

func (s *cartService) getAllCartItemWithNoAct(ctx context.Context, email string) (*[]dto.GetCartItemRes, error) {
	cart, err := s.repo.CheckUserCart(ctx, email)
	if err != nil {
		return nil, errors.New("failed check user cart: " + err.Error())
	}

	if !cart {
		return nil, errors.New("user doesn't have a cart")
	}

	items, err := s.repo.GetAllCartItem(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get all item in user cart: " + err.Error())
	}

	itemRes := make([]dto.GetCartItemRes, len(*items))
	for i, item := range *items {
		store, err := s.storeRepo.GetStore(ctx, item.StoreID)
		if err != nil {
			return nil, errors.New("failed to get store: " + err.Error())
		}

		itemRes[i] = dto.GetCartItemRes{
			Product_Id:    item.Product_Id,
			Product_Name:  item.Product_Name,
			Product_Image: item.Product_Image,
			Store_Name:    store.Name,
			Quantity:      item.Quantity,
			AddedAt:       item.AddedAt,
			Selected:      item.Selected,
			Price:         item.Price,
		}
	}

	return &itemRes, nil
}

// AddToCart implements domain.CartService.
func (s *cartService) AddToCart(ctx context.Context, email, productID string, req *dto.AddCartReq) error {
	err := s.delRedisCartItem(email)
	if err != nil {
		return err
	}

	cart, err := s.repo.CheckUserCart(ctx, email)
	if err != nil {
		return errors.New("failed check user cart: " + err.Error())
	}

	if !cart {
		return errors.New("user doesn't have a cart")
	}

	product, err := s.productRepo.GetProductById(ctx, productID)
	if err != nil {
		return errors.New("failed to get product: " + err.Error())
	}

	if req.Quantity > product.Stock {
		return errors.New("product stock is less than quantity")
	}

	item := domain.CartItem{
		Product_Id:    productID,
		Product_Name:  product.Name,
		Product_Image: product.Images,
		StoreID:       product.Store_id,
		Quantity:      req.Quantity,
		AddedAt:       time.Now(),
		Selected:      true,
		Price:         product.Price,
	}

	totalPrice := float64(req.Quantity) * product.Price

	result, err := s.repo.AddToCart(ctx, email, &item)
	if err != nil {
		return errors.New("failed add this product to the cart")
	}

	if result.ModifiedCount == 0 {
		return errors.New("no product was added to the cart")
	}

	if err := s.repo.UpdateTotalPrice(ctx, email, totalPrice); err != nil {
		return errors.New("failed update total price in the cart")
	}

	defer func() {
		if err := s.updateRedisCart(ctx, email); err != nil {
			log.Println("failed to update user cart in cache: ", err)
		}
	}()

	return nil

}

// RemoveCartItemById implements domain.CartService.
func (s *cartService) RemoveCartItemById(ctx context.Context, email string, productID string) error {
	err := s.delRedisCartItem(email)
	if err != nil {
		return err
	}

	cartExist, err := s.repo.CheckUserCart(ctx, email)
	if err != nil {
		return errors.New("failed to check user cart: " + err.Error())
	}

	if !cartExist {
		return errors.New("user doesn't have a cart")
	}

	_, err = s.productRepo.GetProductById(ctx, productID)
	if err != nil {
		return errors.New("failed to get product: " + err.Error())
	}

	updateAT := time.Now()

	result, err := s.repo.RemoveCartItemById(ctx, email, productID, updateAT)
	if err != nil {
		return errors.New("failed to remove item from user cart: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return errors.New("no item was removed")
	}

	defer func() {
		if err := s.updateRedisCart(ctx, email); err != nil {
			log.Println("failed to update user cart in cache: ", err)
		}
	}()

	return nil
}

// CreateCart implements domain.CartService.
func (s *cartService) CreateCart(ctx context.Context, email string) error {
	cartExist, err := s.repo.CheckUserCart(ctx, email)
	if err != nil {
		return errors.New("failed check user cart: " + err.Error())
	}

	if cartExist {
		return errors.New("user already have a cart")
	}

	cart := domain.Cart{
		ID:         primitive.NewObjectID(),
		Email:      email,
		UpdatedAt:  time.Now(),
		TotalPrice: 0,
		Items:      make([]domain.CartItem, 0),
	}

	if err := s.repo.CreateCart(ctx, &cart); err != nil {
		return errors.New("failed to create user cart: " + err.Error())
	}

	return nil
}

// GetAllCartItem implements domain.CartService.
func (s *cartService) GetAllCartItem(ctx context.Context, email string) (*[]dto.GetCartItemRes, error) {
	val, err := s.cacheRepo.Get("usercart-item:" + email)
	if err == nil {
		var cartitem []dto.GetCartItemRes
		err = json.Unmarshal(val, &cartitem)
		if err != nil {
			return nil, errors.New("failed to unmarshal user cart item: " + err.Error())
		}
		return &cartitem, nil
	}

	cart, err := s.repo.CheckUserCart(ctx, email)
	if err != nil {
		return nil, errors.New("failed check user cart: " + err.Error())
	}

	if !cart {
		return nil, errors.New("user doesn't have a cart")
	}

	items, err := s.repo.GetAllCartItem(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get all item in user cart: " + err.Error())
	}

	itemRes := make([]dto.GetCartItemRes, len(*items))
	for i, item := range *items {
		store, err := s.storeRepo.GetStore(ctx, item.StoreID)
		if err != nil {
			return nil, errors.New("failed to get store: " + err.Error())
		}

		itemRes[i] = dto.GetCartItemRes{
			Product_Id:    item.Product_Id,
			Product_Name:  item.Product_Name,
			Product_Image: item.Product_Image,
			Store_Name:    store.Name,
			Quantity:      item.Quantity,
			AddedAt:       item.AddedAt,
			Selected:      item.Selected,
			Price:         item.Price,
		}
	}

	err = s.setRedisCart(itemRes, email)
	if err != nil {
		return nil, err
	}

	return &itemRes, nil
}

// GetUserCart implements domain.CartService.
func (s *cartService) GetUserCart(ctx context.Context, email string) (*domain.Cart, error) {
	val, err := s.cacheRepo.Get("usercart:" + email)
	if err == nil {
		var cart domain.Cart
		err = json.Unmarshal(val, &cart)
		if err != nil {
			return nil, errors.New("failed to unmarshal user cart: " + err.Error())
		}
		return &cart, nil
	}

	cartExist, err := s.repo.CheckUserCart(ctx, email)
	if err != nil {
		return nil, errors.New("failed check user cart: " + err.Error())
	}

	if !cartExist {
		return nil, errors.New("user doesn't have a cart")
	}

	cart, err := s.repo.GetUserCart(ctx, email)
	if err != nil {
		return nil, errors.New("failed to get user cart: " + err.Error())
	}

	cartData, err := json.Marshal(cart)
	if err != nil {
		return nil, errors.New("failed to marshal cart data: " + err.Error())
	}

	err = s.cacheRepo.Set("user:"+email, cartData, time.Hour*24)
	if err != nil {
		return nil, errors.New("failed to set user cart in cache: " + err.Error())
	}

	return cart, nil
}

// UpdateCartItemById implements domain.CartService.
func (s *cartService) UpdateCartItemById(ctx context.Context, email string, productID string, input *dto.CartItemEditReq) error {
	err := s.delRedisCartItem(email)
	if err != nil {
		return err
	}

	cartExist, err := s.repo.CheckUserCart(ctx, email)
	if err != nil {
		return errors.New("failed to check user cart: " + err.Error())
	}

	if !cartExist {
		return errors.New("user doesn't have a cart")
	}

	cart, err := s.repo.GetUserCart(ctx, email)
	if err != nil {
		return errors.New("failed to get user cart: " + err.Error())
	}

	product, err := s.productRepo.GetProductById(ctx, productID)
	if err != nil {
		return errors.New("failed to get product: " + err.Error())
	}

	updateAT := time.Now()
	oldQuantity := 0
	for _, item := range cart.Items {
		if item.Product_Id == productID {
			oldQuantity = item.Quantity
			break
		}
	}

	// Ensure that the quantity is not less than 0
	if input.Quantity < 0 {
		return errors.New("quantity cannot be less than 0")
	}

	// Calculate the total price difference
	totalPriceDifference := float64(input.Quantity-oldQuantity) * product.Price

	// If the total price difference is negative, set it to 0
	if totalPriceDifference < 0 {
		totalPriceDifference = 0
	}

	update := dto.CartItemEditRepo{
		Quantity:    input.Quantity,
		Selected:    input.Selected,
		Total_Price: totalPriceDifference,
	}
	result, err := s.repo.UpdateCartItemById(ctx, email, productID, &update, updateAT)
	if err != nil {
		return errors.New("failed to update product in the cart: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return errors.New("no item was updated")
	}

	defer func() {
		if err := s.updateRedisCart(ctx, email); err != nil {
			log.Println("failed to update user cart in cache: ", err)
		}
	}()

	return nil

}
