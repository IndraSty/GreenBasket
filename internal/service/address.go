package service

import (
	"context"
	"errors"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/asaskevich/govalidator"
)

type addressService struct {
	repo       domain.AddressRepository
	storeRepo  domain.StoreRepository
	userRepo   domain.UserRepository
	sellerRepo domain.SellerRepository
}

func NewAddressService(repo domain.AddressRepository, sellerRepo domain.SellerRepository, userRepo domain.UserRepository, storeRepo domain.StoreRepository) domain.AddressService {
	return &addressService{
		repo:       repo,
		storeRepo:  storeRepo,
		userRepo:   userRepo,
		sellerRepo: sellerRepo,
	}
}

// AddSellerAddress implements domain.AddressService.
func (s *addressService) AddSellerAddress(ctx context.Context, email string, req domain.Address) (*dto.AddressRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	seller, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil || seller == nil {
		return nil, errors.New("seller not found" + err.Error())
	}

	if seller.Address_Details != nil {
		return nil, errors.New("seller already has an address")
	}

	updateAT := time.Now()

	result, err := s.repo.AddSellerAddress(ctx, email, req, updateAT)
	if err != nil {
		return nil, errors.New("Failed to update store: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no address was added")
	}

	return &dto.AddressRes{
		UpdateResult: result,
	}, nil
}

// GetSellerAddress implements domain.AddressService.
func (s *addressService) GetSellerAddress(ctx context.Context, email string) (*domain.Address, error) {
	seller, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil || seller == nil {
		return nil, errors.New("seller not found" + err.Error())
	}

	result, err := s.repo.GetSellerAddress(ctx, email)
	if err != nil {
		return nil, errors.New("seller haven't address" + err.Error())
	}

	return result, nil
}

// RemoveSellerAddress implements domain.AddressService.
func (s *addressService) RemoveSellerAddress(ctx context.Context, email string) (*dto.AddressDelRes, error) {
	seller, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil || seller == nil {
		return nil, errors.New("seller not found" + err.Error())
	}

	address, err := s.repo.GetSellerAddress(ctx, email)
	if err != nil || address == nil {
		return nil, errors.New("seller haven't address" + err.Error())
	}

	updateAT := time.Now()
	result, err := s.repo.RemoveSellerAddress(ctx, email, updateAT)
	if err != nil {
		return nil, errors.New("Failed to delete address: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("address was not removed")
	}

	return &dto.AddressDelRes{
		UpdateResult: result,
	}, nil
}

// UpdateSellerAddress implements domain.AddressService.
func (s *addressService) UpdateSellerAddress(ctx context.Context, email string, req domain.Address) (*dto.AddressRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	seller, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil || seller == nil {
		return nil, errors.New("seller not found" + err.Error())
	}

	address, err := s.repo.GetSellerAddress(ctx, email)
	if err != nil || address == nil {
		return nil, errors.New("seller haven't address" + err.Error())
	}

	updateAT := time.Now()
	result, err := s.repo.UpdateSellerAddress(ctx, email, req, updateAT)
	if err != nil {
		return nil, errors.New("Failed to update address: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("address was not updated")
	}

	return &dto.AddressRes{
		UpdateResult: result,
	}, nil
}

// AddUserAddress implements domain.AddressService.
func (s *addressService) AddUserAddress(ctx context.Context, email string, req domain.Address) (*dto.AddressRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("user not found" + err.Error())
	}

	if user.Address_Details != nil {
		return nil, errors.New("user already has an address")
	}

	updateAT := time.Now()

	result, err := s.repo.AddUserAddress(ctx, email, req, updateAT)
	if err != nil {
		return nil, errors.New("Failed to update store: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no address was added")
	}

	return &dto.AddressRes{
		UpdateResult: result,
	}, nil
}

// GetUserAddress implements domain.AddressService.
func (s *addressService) GetUserAddress(ctx context.Context, email string) (*domain.Address, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("user not found" + err.Error())
	}

	result, err := s.repo.GetUserAddress(ctx, email)
	if err != nil {
		return nil, errors.New("user haven't address" + err.Error())
	}

	return result, nil
}

// RemoveUserAddress implements domain.AddressService.
func (s *addressService) RemoveUserAddress(ctx context.Context, email string) (*dto.AddressDelRes, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("user not found" + err.Error())
	}

	address, err := s.repo.GetUserAddress(ctx, email)
	if err != nil || address == nil {
		return nil, errors.New("user haven't address" + err.Error())
	}

	updateAT := time.Now()
	result, err := s.repo.RemoveAddressUser(ctx, email, updateAT)
	if err != nil {
		return nil, errors.New("Failed to delete address: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("address was not removed")
	}

	return &dto.AddressDelRes{
		UpdateResult: result,
	}, nil
}

// UpdateUserAddress implements domain.AddressService.
func (s *addressService) UpdateUserAddress(ctx context.Context, email string, req domain.Address) (*dto.AddressRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, errors.New("user not found" + err.Error())
	}

	address, err := s.repo.GetUserAddress(ctx, email)
	if err != nil || address == nil {
		return nil, errors.New("user haven't address" + err.Error())
	}

	updateAT := time.Now()
	result, err := s.repo.UpdateAddressUser(ctx, email, req, updateAT)
	if err != nil {
		return nil, errors.New("Failed to update address: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("address was not updated")
	}

	return &dto.AddressRes{
		UpdateResult: result,
	}, nil
}

// AddStoreAddress implements domain.AddressService.
func (s *addressService) AddStoreAddress(ctx context.Context, email string, storeID string, req domain.Address) (*dto.AddressRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	if store.Address_Details != nil {
		return nil, errors.New("store already has an address")
	}

	updateAT := time.Now()

	result, err := s.repo.AddStoreAddress(ctx, email, storeID, req, updateAT)
	if err != nil {
		return nil, errors.New("Failed to update address: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no address was added")
	}

	return &dto.AddressRes{
		UpdateResult: result,
	}, nil
}

// GetStoreAddress implements domain.AddressService.
func (s *addressService) GetStoreAddress(ctx context.Context, email string, storeID string) (*domain.Address, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	result, err := s.repo.GetStoreAddress(ctx, email, storeID)
	if err != nil {
		return nil, errors.New("store haven't address" + err.Error())
	}

	return result, nil
}

// RemoveStoreAddress implements domain.AddressService.
func (s *addressService) RemoveStoreAddress(ctx context.Context, email string, storeID string) (*dto.AddressDelRes, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	address, err := s.repo.GetStoreAddress(ctx, email, storeID)
	if err != nil || address == nil {
		return nil, errors.New("store haven't address" + err.Error())
	}

	updateAT := time.Now()
	result, err := s.repo.RemoveStoreAddress(ctx, email, storeID, updateAT)
	if err != nil {
		return nil, errors.New("Failed to delete address: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("address was not removed")
	}

	return &dto.AddressDelRes{
		UpdateResult: result,
	}, nil
}

// UpdateStoreAddress implements domain.AddressService.
func (s *addressService) UpdateStoreAddress(ctx context.Context, email string, storeID string, req domain.Address) (*dto.AddressRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	address, err := s.repo.GetStoreAddress(ctx, email, storeID)
	if err != nil || address == nil {
		return nil, errors.New("store haven't address" + err.Error())
	}

	updateAT := time.Now()
	result, err := s.repo.UpdateStoreAddress(ctx, email, storeID, req, updateAT)
	if err != nil {
		return nil, errors.New("Failed to update address: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("address was not updated")
	}

	return &dto.AddressRes{
		UpdateResult: result,
	}, nil
}
