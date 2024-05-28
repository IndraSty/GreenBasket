package service

import (
	"context"
	"errors"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/asaskevich/govalidator"
)

type contactService struct {
	repo      domain.ContactRepository
	storeRepo domain.StoreRepository
}

func NewContactService(repo domain.ContactRepository, storeRepo domain.StoreRepository) domain.ContactService {
	return &contactService{
		repo:      repo,
		storeRepo: storeRepo,
	}
}

// AddStoreContact implements domain.ContactService.
func (s *contactService) AddStoreContact(ctx context.Context, email string, storeID string, req domain.Contact) (*dto.ContactRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	if store.Contact_Details != nil {
		return nil, errors.New("store already has an contact")
	}

	updateAT := time.Now()

	result, err := s.repo.AddStoreContact(ctx, email, storeID, req, updateAT)
	if err != nil {
		return nil, errors.New("Failed to update store: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no contact was added")
	}

	return &dto.ContactRes{
		UpdateResult: result,
	}, nil
}

// GetStoreContact implements domain.ContactService.
func (s *contactService) GetStoreContact(ctx context.Context, email string, storeID string) (*domain.Contact, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	result, err := s.repo.GetStoreContact(ctx, email, storeID)
	if err != nil {
		return nil, errors.New("store haven't address")
	}

	return result, nil
}

// RemoveStoreContact implements domain.ContactService.
func (s *contactService) RemoveStoreContact(ctx context.Context, email string, storeID string) (*dto.ContactDelRes, error) {
	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	contact, err := s.repo.GetStoreContact(ctx, email, storeID)
	if err != nil || contact == nil {
		return nil, errors.New("store haven't address" + err.Error())
	}

	updateAT := time.Now()
	result, err := s.repo.RemoveStoreContact(ctx, email, storeID, updateAT)
	if err != nil {
		return nil, errors.New("Failed to delete contact: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("contact was not removed")
	}

	return &dto.ContactDelRes{
		UpdateResult: result,
	}, nil
}

// UpdateStoreContact implements domain.ContactService.
func (s *contactService) UpdateStoreContact(ctx context.Context, email string, storeID string, req domain.Contact) (*dto.ContactRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, email)
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	contact, err := s.repo.GetStoreContact(ctx, email, storeID)
	if err != nil || contact == nil {
		return nil, errors.New("store haven't address" + err.Error())
	}

	updateAT := time.Now()
	result, err := s.repo.UpdateStoreContact(ctx, email, storeID, req, updateAT)
	if err != nil {
		return nil, errors.New("Failed to update contact: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("contact was not updated")
	}

	return &dto.ContactRes{
		UpdateResult: result,
	}, nil
}
