package main

import (
	"fmt"
	"net/http"
	"os"
	"user-ms/src/auth0"
	"user-ms/src/handler"
	"user-ms/src/model"
	"user-ms/src/rabbitmq"
	"user-ms/src/repository"
	"user-ms/src/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
)

var db *gorm.DB
var err error

func initDB() (*gorm.DB, error) {
	host := os.Getenv("DATABASE_DOMAIN")
	user := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	name := os.Getenv("DATABASE_SCHEMA")
	port := os.Getenv("DATABASE_PORT")

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host,
		user,
		password,
		name,
		port,
	)
	db, err = gorm.Open("postgres", connectionString)

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.FollowingRequest{})
	db.AutoMigrate(model.Follower{})
	return db, err
}

func initUserRepo(database *gorm.DB) *repository.UserRepository {
	return &repository.UserRepository{Database: database}
}

func initAuth0Client() *auth0.Auth0Client {
	domain := os.Getenv("AUTH0_DOMAIN")
	clientId := os.Getenv("AUTH0_CLIENT_ID")
	clientSecret := os.Getenv("AUTH0_CLIENT_SECRET")
	audience := os.Getenv("AUTH0_AUDIENCE")

	client := auth0.NewAuth0Client(domain, clientId, clientSecret, audience)
	return &client
}

func initUserService(userRepo *repository.UserRepository, auth0Client *auth0.Auth0Client) *service.UserService {
	return &service.UserService{UserRepo: userRepo, Auth0Client: *auth0Client}
}

func initUserHandler(service *service.UserService) *handler.UserHandler {
	return &handler.UserHandler{Service: service}
}

func handleUserFunc(handler *handler.UserHandler, router *gin.Engine) {
	router.POST("/register", handler.Register)
	router.GET("/users", handler.GetByEmail)
	router.PUT("/users", handler.Update)
	router.PUT("/users/block-user", handler.BlockUser)
	router.PUT("/users/unblock-user", handler.UnblockUser)
	router.GET("/users/blocked-users", handler.GetBlockedUsers)
	router.POST("/users/set-notifications", handler.SetNotifications)
	router.GET("/users/get-notifications", handler.GetNotifications)
	router.GET("/users/:id", handler.GetByID)
}

func initFollowerRepository(database *gorm.DB) *repository.FollowerRepository {
	return &repository.FollowerRepository{Database: database}
}

func initFollowingRequestRepository(database *gorm.DB) *repository.FollowingRequestRepository {
	return &repository.FollowingRequestRepository{Database: database}
}

func initFollowingService(followerRepository *repository.FollowerRepository, followingRequestRepository *repository.FollowingRequestRepository, userRepository *repository.UserRepository, channel *amqp.Channel) *service.FollowingService {
	return &service.FollowingService{FollowerRepository: followerRepository, FollowingRequestRepository: followingRequestRepository, UserRepository: userRepository, RabbitMQChannel: channel}
}

func initFollowingHandler(service *service.FollowingService) *handler.FollowingHandler {
	return &handler.FollowingHandler{Service: service}
}

func handleFollowingFunc(handler *handler.FollowingHandler, router *gin.Engine) {
	router.POST("/requests", handler.CreateRequest)
	router.POST("/follower", handler.CreatFollower)
	router.PUT("/requests/:id", handler.UpdateRequest)
	router.GET("/requests", handler.GetRequest)
	router.GET("/requests/:id", handler.GetRequestsByFollowingID)
	router.GET("user/:id/followers", handler.GetFollowers)
	router.GET("user/:id/following", handler.GetFollowing)
	router.DELETE("user/:id/removeFollower/:followingId+", handler.RemoveFollowing)
}

func main() {
	database, _ := initDB()

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	rabbit := rabbitmq.RMQProducer{
		ConnectionString: amqpServerURL,
	}

	channel, _ := rabbit.StartRabbitMQ()

	defer channel.Close()

	port := fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))

	userRepo := initUserRepo(database)
	auth0Client := initAuth0Client()
	userService := initUserService(userRepo, auth0Client)
	userHandler := initUserHandler(userService)

	followingReqRepo := initFollowingRequestRepository(database)
	followerRepo := initFollowerRepository(database)
	followingService := initFollowingService(followerRepo, followingReqRepo, userRepo, channel)
	followingHandler := initFollowingHandler(followingService)

	router := gin.Default()

	handleFollowingFunc(followingHandler, router)
	handleUserFunc(userHandler, router)

	http.ListenAndServe(port, cors.AllowAll().Handler(router))
}
