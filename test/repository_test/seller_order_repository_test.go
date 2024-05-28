package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/dto"
	"github.com/IndraSty/GreenBasket/internal/repository"
	"github.com/IndraSty/GreenBasket/test"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SellerOrderRepositoryTestSuite struct {
	test.MongoTestSuite
	repo domain.SellerOrderRepository
}

func (suite *SellerOrderRepositoryTestSuite) SetupSuite() {
	suite.MongoTestSuite.SetupSuite()
	suite.repo = repository.NewSellerOrderRepository(suite.Client)
}

func (suite *SellerOrderRepositoryTestSuite) TearDownSuite() {
	suite.MongoTestSuite.TearDownSuite()
}

func (suite *SellerOrderRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	suite.MongoTestSuite.BeforeTest(suiteName, testName)
}

func (suite *SellerOrderRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.MongoTestSuite.AfterTest(suiteName, testName)
}

func (suite *SellerOrderRepositoryTestSuite) TestCreateSellerOrderSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	newSellerOrder := domain.SellerOrder{
		ID:             primitive.NewObjectID(),
		Order_id:       "orderid",
		Email:          "testemail@gmail.com",
		Ordered_At:     time.Now(),
		Updated_At:     time.Now(),
		Total_Price:    100,
		Payment_Status: "status",
		Items:          []domain.SellerOrderItem{},
	}

	res, err := suite.repo.CreateOrderSeller(ctx, newSellerOrder)

	suite.Require().NoError(err)
	suite.Require().NotEqual(primitive.NilObjectID, res, "The resulting ID cannot be nil")

}
func (suite *SellerOrderRepositoryTestSuite) TestGetSellerOrderSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	email := "testemail@gmail.com"
	orderID := "orderid"
	newSellerOrder := domain.SellerOrder{
		ID:             primitive.NewObjectID(),
		Order_id:       orderID,
		Email:          email,
		Ordered_At:     time.Now(),
		Updated_At:     time.Now(),
		Total_Price:    100,
		Payment_Status: "status",
		Items:          []domain.SellerOrderItem{},
	}

	_, err := suite.repo.CreateOrderSeller(ctx, newSellerOrder)

	suite.Require().NoError(err)

	sellerOrder, err := suite.repo.GetSellerOrderByEmailAndId(ctx, email, orderID)

	suite.Require().NoError(err, "There should be no errors when searching for a seller order with a valid order ID")
	suite.Require().NotNil(sellerOrder, "The order object cannot be nil when the seller order is discovered")
	suite.Require().Equal(orderID, sellerOrder.Order_id, "The order ID found must be the same as the one searched")
}
func (suite *SellerOrderRepositoryTestSuite) TestGetSellerOrderFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	email := "testsemail@gmail.com"
	orderID := "orderids"

	sellerOrder, err := suite.repo.GetSellerOrderByEmailAndId(ctx, email, orderID)

	suite.Require().Error(err, "An error should occur because the OrderID was not found")
	suite.Require().Nil(sellerOrder, "The user object must be nil because the OrderID was not found")
}

func (suite *SellerOrderRepositoryTestSuite) TestGetAllSellerOrderSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	email := "testemail@gmail.com"
	newSellerOrder := domain.SellerOrder{
		ID:             primitive.NewObjectID(),
		Order_id:       "orderid",
		Email:          email,
		Ordered_At:     time.Now(),
		Updated_At:     time.Now(),
		Total_Price:    100,
		Payment_Status: "status",
		Items:          []domain.SellerOrderItem{},
	}

	_, err := suite.repo.CreateOrderSeller(ctx, newSellerOrder)

	suite.Require().NoError(err)

	sellerOrders, err := suite.repo.GetAllSellerOrders(ctx, email)

	suite.Require().NoError(err, "There should be no errors when searching for sellers orders with a valid email")
	suite.Require().NotNil(sellerOrders, "The seller orders object cannot be nil when the orders is discovered")
}
func (suite *SellerOrderRepositoryTestSuite) TestGetAllSellerOrderFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	email := "testemail@gmail.com"

	sellerOrders, err := suite.repo.GetAllSellerOrders(ctx, email)

	suite.Require().Error(err, "An error should occur because the Email was not found")
	suite.Require().Nil(sellerOrders, "The seller orders object must be nil because the Email was not found")
}

func (suite *SellerOrderRepositoryTestSuite) TestUpdateStatusSellerOrderSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var req dto.OrderStatusUpdateReq
	req.Status = "Pending"
	orderID := "uniqueorderid"
	productID := "uniqueproductid"
	newSellerOrder := domain.SellerOrder{
		ID:             primitive.NewObjectID(),
		Order_id:       orderID,
		Email:          "testemail@gmail.com",
		Ordered_At:     time.Now(),
		Updated_At:     time.Now(),
		Total_Price:    100,
		Payment_Status: "status",
		Items: []domain.SellerOrderItem{
			{
				Product_Id:       productID,
				Product_Name:     "productname",
				Product_Image:    []string{"productimage"},
				User_Email:       "usertestemail@gmail.com",
				Status:           "status",
				Quantity:         12,
				Price:            100,
				Address_Shipping: domain.Address{},
			},
		},
	}

	_, err := suite.repo.CreateOrderSeller(ctx, newSellerOrder)

	suite.Require().NoError(err)

	result, err := suite.repo.UpdateStatusOrderSeller(ctx, orderID, productID, &req)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(1, result.ModifiedCount, "Should have modified 1 document")
}
func (suite *SellerOrderRepositoryTestSuite) TestUpdateStatusSellerOrderFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var req dto.OrderStatusUpdateReq

	req.Status = "Pending"

	orderID := "uniqueorderid"
	productID := "uniqueproductid"

	result, err := suite.repo.UpdateStatusOrderSeller(ctx, orderID, productID, &req)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(0, result.ModifiedCount, "Should not have modified any document")
}

func (suite *SellerOrderRepositoryTestSuite) TestDeleteSellerOrderItemSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var req dto.OrderStatusUpdateReq
	req.Status = "Pending"
	orderID := "uniqueorderid"
	productID := "uniqueproductid"
	email := "testemail@gmail.com"
	newSellerOrder := domain.SellerOrder{
		ID:             primitive.NewObjectID(),
		Order_id:       orderID,
		Email:          email,
		Ordered_At:     time.Now(),
		Updated_At:     time.Now(),
		Total_Price:    100,
		Payment_Status: "status",
		Items: []domain.SellerOrderItem{
			{
				Product_Id:       productID,
				Product_Name:     "productname",
				Product_Image:    []string{"productimage"},
				User_Email:       "usertestemail@gmail.com",
				Status:           "status",
				Quantity:         12,
				Price:            100,
				Address_Shipping: domain.Address{},
			},
		},
	}

	_, err := suite.repo.CreateOrderSeller(ctx, newSellerOrder)

	suite.Require().NoError(err)

	result, err := suite.repo.DeleteItem(ctx, email, orderID, productID)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(1, result.ModifiedCount, "Should have delete 1 document")
}
func (suite *SellerOrderRepositoryTestSuite) TestDeleteSellerOrderItemFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	email := "testemail@gmail.com"
	orderID := "uniqueorderid"
	productID := "uniqueproductid"

	result, err := suite.repo.DeleteItem(ctx, email, orderID, productID)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(0, result.ModifiedCount, "Should not have deleted any documents")
}

func TestSellerOrderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SellerOrderRepositoryTestSuite))
}
