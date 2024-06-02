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

type UserRepositoryTestSuite struct {
	test.MongoTestSuite
	repo domain.UserRepository
}

func (suite *UserRepositoryTestSuite) SetupSuite() {
	suite.MongoTestSuite.SetupSuite()
	suite.repo = repository.NewUserRepository(suite.Client)
}

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	suite.MongoTestSuite.TearDownSuite()
}

func (suite *UserRepositoryTestSuite) BeforeTest(suiteName, testName string) {
	suite.MongoTestSuite.BeforeTest(suiteName, testName)
}

func (suite *UserRepositoryTestSuite) AfterTest(suiteName, testName string) {
	suite.MongoTestSuite.AfterTest(suiteName, testName)
}

func (suite *UserRepositoryTestSuite) TestCreateUserSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()
	password := "Testpassword123_"
	hashedPassword := util.HashPassword(password)

	id := primitive.NewObjectID()
	userId := id.Hex()

	newUser := domain.User{
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
		User_Id:       userId,
	}

	res, err := suite.repo.CreateUser(ctx, newUser)

	suite.Require().NoError(err)
	suite.Require().NotEqual(primitive.NilObjectID, res, "The resulting ID cannot be nil")
}

func (suite *UserRepositoryTestSuite) TestFindUserByEmailSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser := domain.User{
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
		User_Id:       "uniqueuserid",
	}
	_, err := suite.repo.CreateUser(ctx, newUser)
	suite.Require().NoError(err)

	user, err := suite.repo.FindUserByEmail(ctx, newUser.Email)

	suite.Require().NoError(err, "There should be no errors when searching for users with valid emails")
	suite.Require().NotNil(user, "The user object cannot be nil when the user is discovered")
	suite.Require().Equal(newUser.Email, user.Email, "The email of the found user must be the same as the one searched for")
}

func (suite *UserRepositoryTestSuite) TestFindUserByEmailNotFound() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email := "nonexistentemail@gmail.com"
	user, err := suite.repo.FindUserByEmail(ctx, email)

	suite.Require().Error(err, "An error should occur because the email was not found")
	suite.Require().Nil(user, "The user object must be nil because the email was not found")
}

func (suite *UserRepositoryTestSuite) TestUpdateUserSuccess() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email := "updateuser@gmail.com"
	user := domain.User{
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
		User_Id:       "uniqueuserid",
	}
	_, err := suite.repo.CreateUser(ctx, user)
	suite.Require().NoError(err)

	update := bson.D{{Key: "first_name", Value: "updatedName"}}

	result, err := suite.repo.UpdateUser(ctx, email, update)
	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(1, result.ModifiedCount, "Should have modified 1 document")
}

func (suite *UserRepositoryTestSuite) TestUpdateUserFailNonExisting() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	email := "nonexistinguser@gmail.com"
	update := bson.D{{Key: "first_name", Value: "UpdatedName"}}

	result, err := suite.repo.UpdateUser(ctx, email, update)
	suite.Require().NoError(err)
	suite.Require().NotNil(result)
	suite.Require().EqualValues(0, result.ModifiedCount, "Should not have modified any documents")
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
