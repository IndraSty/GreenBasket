package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/repository"
	"github.com/IndraSty/GreenBasket/test"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoreRepositoryTestSuite struct {
	test.MongoTestSuite
	repo domain.StoreRepository
}

func (suite *StoreRepositoryTestSuite) SetupSuite() {
	suite.MongoTestSuite.SetupSuite()
	suite.repo = repository.NewStoreRepository(suite.Client)
}

func (suite *StoreRepositoryTestSuite) TearDownSuite() {
	suite.MongoTestSuite.TearDownSuite()
}

func (suite *StoreRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	suite.MongoTestSuite.BeforeTest(suiteName, testName)
}

func (suite *StoreRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.MongoTestSuite.AfterTest(suiteName, testName)
}

func (suite *StoreRepositoryTestSuite) TestCreateStoreSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	newStore := domain.Store{
		ID:          primitive.NewObjectID(),
		Name:        "storename",
		Description: "storedescription",
		Logo:        "storelogo",
		Banner:      "storebanner",
		Created_At:  time.Now(),
		Updated_At:  time.Now(),
		Email:       "testemail@gmail.com",
		Store_Id:    "storeid",
	}

	res, err := suite.repo.CreateStore(ctx, newStore)

	suite.Require().NoError(err)
	suite.Require().NotEqual(primitive.NilObjectID, res, "The resulting ID cannot be nil")
}

func (suite *StoreRepositoryTestSuite) TestGetStoreSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	newStore := domain.Store{
		ID:          primitive.NewObjectID(),
		Name:        "storename",
		Description: "storedescription",
		Logo:        "storelogo",
		Banner:      "storebanner",
		Created_At:  time.Now(),
		Updated_At:  time.Now(),
		Email:       "testemail@gmail.com",
		Store_Id:    "storeid",
	}

	_, err := suite.repo.CreateStore(ctx, newStore)
	suite.Require().NoError(err)

	store, err := suite.repo.GetStore(ctx, newStore.Store_Id)

	suite.Require().NoError(err, "There should be no errors when searching for a store with a valid store ID")
	suite.Require().NotNil(store, "The store object cannot be nil when the store is discovered")
	suite.Require().Equal(newStore.Store_Id, store.Store_Id, "The store ID found must be the same as the one searched")

}

func (suite *StoreRepositoryTestSuite) TestGetStoreNotFound() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	storeid := "storeidnoneexist@gmail.com"

	store, err := suite.repo.GetStore(ctx, storeid)

	suite.Require().Error(err, "An error should occur because the storeID was not found")
	suite.Require().Nil(store, "The user object must be nil because the storeID was not found")

}

func (suite *StoreRepositoryTestSuite) TestUpdateStoreSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	email := "testemail@gmail.com"
	storeID := "storeid"
	newStore := domain.Store{
		ID:          primitive.NewObjectID(),
		Name:        "storename",
		Description: "storedescription",
		Logo:        "storelogo",
		Banner:      "storebanner",
		Created_At:  time.Now(),
		Updated_At:  time.Now(),
		Email:       email,
		Store_Id:    storeID,
	}

	_, err := suite.repo.CreateStore(ctx, newStore)
	suite.Require().NoError(err)

	update := bson.D{{Key: "name", Value: "updatedName"}}

	result, err := suite.repo.UpdateStore(ctx, email, storeID, update)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(1, result.ModifiedCount, "Should have modified 1 document")
}

func (suite *StoreRepositoryTestSuite) TestUpdateStoreFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	email := "testemail@gmail.com"
	storeID := "storeid"

	update := bson.D{{Key: "name", Value: "updatedName"}}

	result, err := suite.repo.UpdateStore(ctx, email, storeID, update)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(0, result.ModifiedCount, "Should not have modified any documents")
}

func (suite *StoreRepositoryTestSuite) TestRemoveStoreSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	email := "testemail@gmail.com"
	storeID := "storeid"
	newStore := domain.Store{
		ID:          primitive.NewObjectID(),
		Name:        "storename",
		Description: "storedescription",
		Logo:        "storelogo",
		Banner:      "storebanner",
		Created_At:  time.Now(),
		Updated_At:  time.Now(),
		Email:       email,
		Store_Id:    storeID,
	}

	_, err := suite.repo.CreateStore(ctx, newStore)
	suite.Require().NoError(err)

	result, err := suite.repo.RemoveStore(ctx, email, storeID)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(1, result.DeletedCount, "Should have delete 1 document")
}

func (suite *StoreRepositoryTestSuite) TestRemoveStoreFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	email := "testemail@gmail.com"
	storeID := "storeid"

	result, err := suite.repo.RemoveStore(ctx, email, storeID)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(0, result.DeletedCount, "Should not have deleted any documents")
}

func TestStoreRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(StoreRepositoryTestSuite))
}
