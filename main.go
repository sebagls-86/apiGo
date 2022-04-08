package main

import (
	"api-go/controllers"
	"api-go/services"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server         *gin.Engine
	userservice    services.UserService
	usercontroller controllers.UserController
	ctx            context.Context
	usercollection *mongo.Collection
	mongoclient    *mongo.Client
	err            error
)

func init() {
	ctx = context.TODO()

	errorLoad := godotenv.Load()
	if errorLoad != nil {
		panic(err)
	}

	mongodb := os.Getenv("MONGODB_URI")

	mongoconn := options.Client().ApplyURI(mongodb)
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongo connection established")

	usercollection = mongoclient.Database("userdb").Collection("users")
	userservice = services.NewUserService(usercollection, ctx)
	usercontroller = controllers.New(userservice)
	server = gin.Default()

}

func main() {
	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("v1")
	usercontroller.RegisterUserRoutes(basepath)

	log.Fatal(server.Run(":9090"))
}
