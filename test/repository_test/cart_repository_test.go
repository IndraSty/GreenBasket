package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/repository"
	"github.com/IndraSty/GreenBasket/test"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CartRepositoryTestSuite struct {
	test.MongoTestSuite
	repo domain.CartRepository
}

func (suite *CartRepositoryTestSuite) SetupSuite() {
	suite.MongoTestSuite.SetupSuite()
	suite.repo = repository.NewCartRepository(suite.Client)
}

func (suite *CartRepositoryTestSuite) TearDownSuite() {
	suite.MongoTestSuite.TearDownSuite()
}

func (suite *CartRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	suite.MongoTestSuite.BeforeTest(suiteName, testName)
}

func (suite *CartRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.MongoTestSuite.AfterTest(suiteName, testName)
}

// TestCreateCartSuccess tests the successful creation of a cart
func (suite *CartRepositoryTestSuite) TestCreateCartSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cart := domain.Cart{
		ID:         primitive.NewObjectID(),
		Email:      "test@gmail.com",
		UpdatedAt:  time.Now(),
		TotalPrice: 100.0,
		Items:      []domain.CartItem{},
	}

	err := suite.repo.CreateCart(ctx, &cart)
	suite.Require().NoError(err)

	retrievedCart, err := suite.repo.GetUserCart(ctx, cart.Email)
	suite.Require().NoError(err)
	suite.Require().NotNil(retrievedCart)
	suite.Require().Equal(cart.Email, retrievedCart.Email)
}

func (suite *CartRepositoryTestSuite) TestGetUserCartSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email := "test@gmail.com"

	cart, err := suite.repo.GetUserCart(ctx, email)
	suite.Require().NoError(err)
	suite.Require().NotNil(cart)
	suite.Require().Equal(email, cart.Email)
}

func (suite *CartRepositoryTestSuite) TestGetUserCartFailure() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email := "nonexistent@gmail.com"

	cart, err := suite.repo.GetUserCart(ctx, email)
	suite.Require().Error(err)
	suite.Require().Nil(cart)
}

func (suite *CartRepositoryTestSuite) TestGetAllUserCartSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testTime := time.Now().UTC()

	email := "test@gmail.com"
	items := []domain.CartItem{
		{
			Product_Id:    "productid",
			Product_Name:  "productname",
			Product_Image: []string{"productimage"},
			StoreID:       "storeid",
			Quantity:      12,
			AddedAt:       testTime,
			Selected:      false,
			Price:         10000,
		},
	}
	cart := &domain.Cart{
		ID:         primitive.NewObjectID(),
		Email:      email,
		UpdatedAt:  testTime,
		TotalPrice: 100.0,
		Items:      items,
	}

	err := suite.repo.CreateCart(ctx, cart)
	suite.Require().NoError(err)

	cartItems, err := suite.repo.GetAllCartItem(ctx, email)
	suite.Require().NoError(err)
	suite.Require().NotNil(cartItems)

	// Karena waktu bisa berbeda, kita hanya membandingkan field lainnya
	for i, item := range *cartItems {
		suite.Require().Equal(items[i].Product_Id, item.Product_Id)
		suite.Require().Equal(items[i].Product_Name, item.Product_Name)
		suite.Require().Equal(items[i].Product_Image, item.Product_Image)
		suite.Require().Equal(items[i].StoreID, item.StoreID)
		suite.Require().Equal(items[i].Quantity, item.Quantity)
		suite.Require().Equal(items[i].Selected, item.Selected)
		suite.Require().Equal(items[i].Price, item.Price)
	}
}

func (suite *CartRepositoryTestSuite) TestGetAllUserCartEmpty() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	email := "test@gmail.com"
	cartItems, err := suite.repo.GetAllCartItem(ctx, email)
	suite.Require().NoError(err)
	suite.Require().NotNil(cartItems)
}

func TestCartRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CartRepositoryTestSuite))
}
