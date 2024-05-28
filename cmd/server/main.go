package main

import (
	"github.com/IndraSty/GreenBasket/db"
	"github.com/IndraSty/GreenBasket/internal/bootstrap"
	"github.com/IndraSty/GreenBasket/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	cnf := config.Get()

	db.DBInstance(cnf)
	var client *mongo.Client = db.DBInstance(cnf)

	config.NewAuthSetup(cnf).NewAuth()

	app := gin.New()

	config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"http://greenbasket.com"}
	config.AllowAllOrigins = true

	app.Use(cors.New(config))

	bootstrap.Application(&bootstrap.ApplicationConfig{
		Config: cnf,
		Client: client,
		App:    app,
	})

	_ = app.Run(cnf.Server.Host + ":" + cnf.Server.Port)

}
