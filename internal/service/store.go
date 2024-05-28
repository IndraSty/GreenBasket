package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/asaskevich/govalidator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type storeService struct {
	storeRepo   domain.StoreRepository
	sellerRepo  domain.SellerRepository
	salesReport domain.SalesReportRepository
	cacheRepo   domain.CacheRepository
}

func NewStoreService(storeRepo domain.StoreRepository, sellerRepo domain.SellerRepository,
	salesReport domain.SalesReportRepository, cacheRepo domain.CacheRepository) domain.StoreService {
	return &storeService{
		storeRepo:   storeRepo,
		sellerRepo:  sellerRepo,
		salesReport: salesReport,
		cacheRepo:   cacheRepo,
	}
}

func (s *storeService) setRedisStore(store dto.GetStoreRes, email string) error {
	storeData, err := json.Marshal(store)
	if err != nil {
		return errors.New("failed to marshal store data: " + err.Error())
	}

	err = s.cacheRepo.Set("seller_store:"+email, storeData, time.Hour*24)
	if err != nil {
		return errors.New("failed to set store data in cache: " + err.Error())
	}

	return nil
}

func (s *storeService) delRedisStore(email string) error {
	err := s.cacheRepo.Del("seller_store:" + email)
	if err != nil {
		return errors.New("failed to delete store data in cache: " + err.Error())
	}

	return nil
}

func (s *storeService) updateRedisStore(ctx context.Context, email, storeID string) error {
	store, err := s.getStoreWithNoAct(ctx, email, storeID)
	if err != nil {
		return errors.New("failed to get store data: " + err.Error())
	}

	return s.setRedisStore(*store, email)
}

func (s *storeService) getStoreWithNoAct(ctx context.Context, email string, storeID string) (*dto.GetStoreRes, error) {
	seller, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil || seller == nil {
		return nil, errors.New("seller not found")
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, seller.Email)
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	storeRes := dto.GetStoreRes{
		Name:        store.Name,
		Description: store.Description,
		Logo:        store.Logo,
		Banner:      store.Banner,
		Email:       store.Email,
		Store_Id:    store.Store_Id,
	}

	return &storeRes, nil
}

// CreateStore implements domain.StoreService.
func (s *storeService) CreateStore(ctx context.Context, email string, req *dto.StoreReq) (*dto.AddStoreRes, error) {
	_, err := govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	nameExist, _ := s.storeRepo.CheckNameExists(ctx, req.Name)
	if nameExist {
		return nil, errors.New("store name already added! try to another name")
	}

	seller, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("seller not found" + err.Error())
	}

	id := primitive.NewObjectID()
	storeID := id.Hex()
	seller.Store_Id = storeID

	store := domain.Store{
		ID:          id,
		Store_Id:    storeID,
		Name:        req.Name,
		Description: req.Description,
		Logo:        req.Logo,
		Banner:      req.Banner,
		Created_At:  time.Now(),
		Updated_At:  time.Now(),
		Email:       email,
	}

	salesReport := domain.Sales_Report{
		ID:            primitive.NewObjectID(),
		Store_Id:      storeID,
		Email:         email,
		Total_Sales:   0,
		Total_Incomes: 0.0,
		Products:      make([]domain.Product_Sales, 0),
	}

	result, err := s.storeRepo.CreateStore(ctx, store)
	if err != nil {
		return nil, errors.New("failed to created store" + err.Error())
	}

	_, err = s.salesReport.Insert(ctx, salesReport)
	if err != nil {
		return nil, errors.New("failed to insert sales report: " + err.Error())
	}

	return &dto.AddStoreRes{
		InsertId: &result,
	}, nil
}

// DeleteStore implements domain.StoreService.
func (s *storeService) DeleteStore(ctx context.Context, email, storeID string) (*dto.DeleteStoreRes, error) {

	seller, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil || seller == nil {
		return nil, errors.New("seller not found")
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, seller.Email)
	if err != nil || store == nil {
		return nil, errors.New("store not found" + err.Error())
	}

	result, err := s.storeRepo.RemoveStore(ctx, seller.Email, store.Store_Id)
	if err != nil {
		return nil, errors.New("failed to delete store: " + err.Error())
	}

	if result.DeletedCount == 0 {
		return nil, errors.New("no item was deleted")
	}

	err = s.delRedisStore(email)
	if err != nil {
		return nil, err
	}

	return &dto.DeleteStoreRes{
		DeleteResult: result,
	}, nil
}

// GetStoreByIdAndEmail implements domain.StoreService.
func (s *storeService) GetStoreByIdAndEmail(ctx context.Context, email string, storeID string) (*dto.GetStoreRes, error) {
	val, err := s.cacheRepo.Get("seller_store:" + email)
	if err == nil {
		var store dto.GetStoreRes
		err = json.Unmarshal(val, &store)
		if err != nil {
			return nil, errors.New("failed to unmarshal seller store data: " + err.Error())
		}

		return &store, nil
	}

	seller, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil || seller == nil {
		return nil, errors.New("seller not found")
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, seller.Email)
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	storeRes := dto.GetStoreRes{
		Name:        store.Name,
		Description: store.Description,
		Logo:        store.Logo,
		Banner:      store.Banner,
		Email:       store.Email,
		Store_Id:    store.Store_Id,
	}

	storeData, err := json.Marshal(storeRes)
	if err != nil {
		return nil, errors.New("failed to marshal store data: " + err.Error())
	}

	err = s.cacheRepo.Set("seller_store:"+email, storeData, time.Hour*24)
	if err != nil {
		return nil, errors.New("failed to set store data in cache: " + err.Error())
	}

	return &storeRes, nil
}

// UpdateStore implements domain.StoreService.
func (s *storeService) UpdateStore(ctx context.Context, email, storeID string, req *dto.StoreReq) (*dto.UpdateStoreRes, error) {
	err := s.delRedisStore(email)
	if err != nil {
		return nil, err
	}

	_, err = govalidator.ValidateStruct(req)
	if err != nil {
		return nil, errors.New("Invalid request body" + err.Error())
	}

	seller, err := s.sellerRepo.FindSellerByEmail(ctx, email)
	if err != nil || seller == nil {
		return nil, errors.New("seller not found")
	}

	store, err := s.storeRepo.GetStore(ctx, storeID, seller.Email)
	if err != nil || store == nil {
		return nil, errors.New("store not found")
	}

	var update primitive.D
	if req.Name != "" {
		update = append(update, bson.E{Key: "name", Value: req.Name})
	}
	if req.Banner != "" {
		update = append(update, bson.E{Key: "banner", Value: req.Banner})
	}
	if req.Description != "" {
		update = append(update, bson.E{Key: "description", Value: req.Description})
	}
	if req.Logo != "" {
		update = append(update, bson.E{Key: "logo", Value: req.Logo})
	}

	result, err := s.storeRepo.UpdateStore(ctx, seller.Email, store.Store_Id, update)
	if err != nil {
		return nil, errors.New("Failed to update store: " + err.Error())
	}

	if result.ModifiedCount == 0 {
		return nil, errors.New("no store was updated")
	}

	defer func() {
		if err := s.updateRedisStore(ctx, email, storeID); err != nil {
			log.Println("failed to update store data in cache: ", err)
		}
	}()

	return &dto.UpdateStoreRes{
		UpdateResult: result,
	}, nil
}

// SearchStore implements domain.StoreService.
func (s *storeService) SearchStore(ctx context.Context, query string) ([]domain.Store, error) {
	stores, err := s.storeRepo.GetStoreByQuery(ctx, query)
	if err != nil {
		return nil, errors.New("failed to get all store by query: " + err.Error())
	}

	return stores, nil
}
