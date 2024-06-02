package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/IndraSty/GreenBasket/domain"
	"github.com/IndraSty/GreenBasket/internal/repository"
	"github.com/IndraSty/GreenBasket/internal/util"
	"github.com/IndraSty/GreenBasket/test"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SellerRepositoryTestSuite struct {
	test.MongoTestSuite
	repo domain.SellerRepository
}

func (suite *SellerRepositoryTestSuite) SetupSuite() {
	suite.MongoTestSuite.SetupSuite()
	suite.repo = repository.NewSellerRepository(suite.Client)
}

func (suite *SellerRepositoryTestSuite) TearDownSuite() {
	suite.MongoTestSuite.TearDownSuite()
}

func (suite *SellerRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	suite.MongoTestSuite.BeforeTest(suiteName, testName)
}

func (suite *SellerRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.MongoTestSuite.AfterTest(suiteName, testName)
}

func (suite *SellerRepositoryTestSuite) TestCreateSellerSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	password := "Testpassword123_"
	hashedPassword := util.HashPassword(password)

	id := primitive.NewObjectID()
	sellerId := id.Hex()

	newSeller := domain.Seller{
		ID:            id,
		First_Name:    "Test",
		Last_Name:     "Name",
		Email:         "emailtest@gmail.com",
		Password:      hashedPassword,
		Phone:         "0927328237474",
		Role:          "User",
		Created_At:    time.Now(),
		Updated_At:    time.Now(),
		EmailVerified: false,
		Seller_Id:     sellerId,
	}

	res, err := suite.repo.CreateSeller(ctx, newSeller)

	suite.Require().NoError(err)
	suite.Require().NotEqual(primitive.NilObjectID, res, "ID yang dihasilkan tidak boleh nil")
}

func (suite *SellerRepositoryTestSuite) TestFindSellerByEmailSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newSeller := domain.Seller{
		ID:            primitive.NewObjectID(),
		First_Name:    "Test",
		Last_Name:     "User",
		Email:         "testuser@example.com",
		Password:      "hashedpassword",
		Phone:         "1234567890",
		Role:          "User",
		Created_At:    time.Now(),
		Updated_At:    time.Now(),
		EmailVerified: true,
		Seller_Id:     "uniquesellerid",
	}
	_, err := suite.repo.CreateSeller(ctx, newSeller)
	suite.Require().NoError(err)

	user, err := suite.repo.FindSellerByEmail(ctx, newSeller.Email)

	suite.Require().NoError(err, "Tidak seharusnya ada error saat mencari pengguna dengan email yang valid")
	suite.Require().NotNil(user, "Objek user tidak boleh nil saat pengguna ditemukan")
	suite.Require().Equal(newSeller.Email, user.Email, "Email dari pengguna yang ditemukan harus sama dengan yang dicari")
}

func (suite *SellerRepositoryTestSuite) TestFindSellerByEmailNotFound() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email := "nonexistentemail@example.com"
	user, err := suite.repo.FindSellerByEmail(ctx, email)

	suite.Require().Error(err, "Seharusnya terjadi error karena email tidak ditemukan")
	suite.Require().Nil(user, "Objek user harus nil karena email tidak ditemukan")
}

func (suite *SellerRepositoryTestSuite) TestUpdateSellerSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email := "updateuser@example.com"
	seller := domain.Seller{
		ID:            primitive.NewObjectID(),
		First_Name:    "Test",
		Last_Name:     "User",
		Email:         email,
		Password:      "hashedpassword",
		Phone:         "1234567890",
		Role:          "User",
		Created_At:    time.Now(),
		Updated_At:    time.Now(),
		EmailVerified: true,
		Seller_Id:     "uniquesellerid",
	}
	_, err := suite.repo.CreateSeller(ctx, seller)
	suite.Require().NoError(err)

	update := bson.D{{Key: "first_name", Value: "updatedName"}}

	result, err := suite.repo.UpdateSeller(ctx, email, update)
	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(1, result.ModifiedCount, "Should have modified 1 document")
}

func (suite *SellerRepositoryTestSuite) TestUpdateSellerFailNonExisting() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email := "nonexistinguser@example.com"
	update := bson.D{{Key: "first_name", Value: "UpdatedName"}}

	result, err := suite.repo.UpdateSeller(ctx, email, update)
	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(0, result.ModifiedCount, "Should not have modified any documents")
}

func TestSellerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SellerRepositoryTestSuite))
}
