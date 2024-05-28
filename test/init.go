package test

import (
	"context"
	"time"

	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTestSuite struct {
	suite.Suite
	Client *mongo.Client
	db     *mongo.Database
}

func (suite *MongoTestSuite) SetupSuite() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("error when load envi %s", err.Error())
	// }
	// mongoURI := os.Getenv("MONGO_URI")
	mongoURI := "mongodb://localhost:27017"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	suite.Require().NoError(err, "Koneksi ke MongoDB gagal")

	suite.Client = client
	suite.db = client.Database("greenbasket_test")
}

func (suite *MongoTestSuite) TearDownSuite() {
	err := suite.Client.Disconnect(context.Background())
	suite.Require().NoError(err, "Gagal memutus koneksi dari MongoDB")
}

func (suite *MongoTestSuite) BeforeTest(suiteName, testName string) {
	err := suite.db.Drop(context.Background())
	suite.Require().NoError(err, "Gagal menghapus database sebelum test")
}

func (suite *MongoTestSuite) AfterTest(suiteName, testName string) {
	err := suite.db.Drop(context.Background())
	suite.Require().NoError(err, "Gagal menghapus database setelah test")
}
