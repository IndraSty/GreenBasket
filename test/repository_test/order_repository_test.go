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

type OrderRepositoryTestSuite struct {
	test.MongoTestSuite
	repo domain.OrderRepository
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	suite.MongoTestSuite.SetupSuite()
	suite.repo = repository.NewOrderRepository(suite.Client)
}

func (suite *OrderRepositoryTestSuite) TearDownSuite() {
	suite.MongoTestSuite.TearDownSuite()
}

func (suite *OrderRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	suite.MongoTestSuite.BeforeTest(suiteName, testName)
}

func (suite *OrderRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.MongoTestSuite.AfterTest(suiteName, testName)
}

func (suite *OrderRepositoryTestSuite) TestCreateOrderSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	newOrder := domain.Orders{
		ID:               primitive.NewObjectID(),
		Order_id:         "orderid",
		Email:            "testemail@gmail.com",
		Order_Date:       time.Now(),
		Updated_At:       time.Now(),
		Total_Price:      12000,
		Address_Shipping: domain.Address{},
		Payment:          &domain.PaymentOrder{},
		Items:            []domain.OrderItem{},
	}

	res, err := suite.repo.CreateOrder(ctx, newOrder)

	suite.Require().NoError(err)
	suite.Require().NotEqual(primitive.NilObjectID, res, "The resulting ID cannot be nil")
}

func (suite *OrderRepositoryTestSuite) TestGetOrderSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderID := "orderid"
	newOrder := domain.Orders{
		ID:               primitive.NewObjectID(),
		Order_id:         orderID,
		Email:            "testemail@gmail.com",
		Order_Date:       time.Now(),
		Updated_At:       time.Now(),
		Total_Price:      12000,
		Address_Shipping: domain.Address{},
		Payment:          &domain.PaymentOrder{},
		Items:            []domain.OrderItem{},
	}

	_, err := suite.repo.CreateOrder(ctx, newOrder)
	suite.Require().NoError(err)

	order, err := suite.repo.GetOrder(ctx, orderID)

	suite.Require().NoError(err, "There should be no errors when searching for a order with a valid order ID")
	suite.Require().NotNil(order, "The order object cannot be nil when the order is discovered")
	suite.Require().Equal(orderID, order.Order_id, "The order ID found must be the same as the one searched")
}

func (suite *OrderRepositoryTestSuite) TestGetOrderFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderID := "orderid"

	order, err := suite.repo.GetOrder(ctx, orderID)

	suite.Require().Error(err, "An error should occur because the OrderID was not found")
	suite.Require().Nil(order, "The user object must be nil because the OrderID was not found")
}

func (suite *OrderRepositoryTestSuite) TestGetAllOrderSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	email := "testemail@gmail.com"
	newOrder := domain.Orders{
		ID:               primitive.NewObjectID(),
		Order_id:         "orderid",
		Email:            email,
		Order_Date:       time.Now(),
		Updated_At:       time.Now(),
		Total_Price:      12000,
		Address_Shipping: domain.Address{},
		Payment:          &domain.PaymentOrder{},
		Items:            []domain.OrderItem{},
	}

	_, err := suite.repo.CreateOrder(ctx, newOrder)
	suite.Require().NoError(err)

	orders, err := suite.repo.GetAllOrders(ctx, email)

	suite.Require().NoError(err, "There should be no errors when searching for a orders with a valid email")
	suite.Require().NotNil(orders, "The order object cannot be nil when the orders is discovered")
}

func (suite *OrderRepositoryTestSuite) TestGetAllOrderFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	email := "testsemail@gmail.com"

	orders, err := suite.repo.GetAllOrders(ctx, email)

	suite.Require().Error(err, "An error should occur because the Email was not found")
	suite.Require().Nil(orders, "The orders object must be nil because the Email was not found")
}

func (suite *OrderRepositoryTestSuite) TestUpdateStatusOrderSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var req dto.OrderStatusUpdateReq
	req.Status = "Pending"
	orderID := "uniqueorderid"
	productID := "uniqueproductid"
	newOrder := domain.Orders{
		ID:               primitive.NewObjectID(),
		Order_id:         orderID,
		Email:            "testemail@gmail.com",
		Order_Date:       time.Now(),
		Updated_At:       time.Now(),
		Total_Price:      12000,
		Address_Shipping: domain.Address{},
		Payment:          &domain.PaymentOrder{},
		Items: []domain.OrderItem{
			{
				Product_Id:    productID,
				Product_Name:  "productname",
				Product_Image: []string{"productimage"},
				StoreID:       "storeid",
				Order_Status:  "status",
				Quantity:      12,
				Price:         100,
			},
		},
	}

	_, err := suite.repo.CreateOrder(ctx, newOrder)

	suite.Require().NoError(err)

	result, err := suite.repo.UpdateStatusOrder(ctx, orderID, productID, &req)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(1, result.ModifiedCount, "Should have modified 1 document")
}

func (suite *OrderRepositoryTestSuite) TestUpdateStatusOrderFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var req dto.OrderStatusUpdateReq

	req.Status = "Pending"

	orderID := "uniqueorderid"
	productID := "uniqueproductid"

	result, err := suite.repo.UpdateStatusOrder(ctx, orderID, productID, &req)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(0, result.ModifiedCount, "Should not have modified any document")
}

func (suite *OrderRepositoryTestSuite) TestDeleteItemOrderSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderID := "uniqueorderid"
	productID := "uniqueproductid"
	newOrder := domain.Orders{
		ID:               primitive.NewObjectID(),
		Order_id:         orderID,
		Email:            "testemail@gmail.com",
		Order_Date:       time.Now(),
		Updated_At:       time.Now(),
		Total_Price:      12000,
		Address_Shipping: domain.Address{},
		Payment:          &domain.PaymentOrder{},
		Items: []domain.OrderItem{
			{
				Product_Id:    productID,
				Product_Name:  "productname",
				Product_Image: []string{"productimage"},
				StoreID:       "storeid",
				Order_Status:  "status",
				Quantity:      12,
				Price:         100,
			},
		},
	}

	_, err := suite.repo.CreateOrder(ctx, newOrder)

	suite.Require().NoError(err)

	result, err := suite.repo.DeleteItem(ctx, orderID, productID)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(1, result.ModifiedCount, "Should have delete 1 document")
}

func (suite *OrderRepositoryTestSuite) TestDeleteItemOrderFailed() {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderID := "uniqueorderid"
	productID := "uniqueproductid"

	result, err := suite.repo.DeleteItem(ctx, orderID, productID)

	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(0, result.ModifiedCount, "Should not have deleted any documents")
}

func TestOrderRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}
