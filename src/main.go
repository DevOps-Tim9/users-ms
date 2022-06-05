package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"user-ms/src/auth0"
	"user-ms/src/handler"
	"user-ms/src/model"
	"user-ms/src/repository"
	"user-ms/src/service"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var db *gorm.DB
var err error

func initDB() *gorm.DB {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB")

	connString := fmt.Sprintf("host=localhost port=%s user=%s dbname=%s sslmode=disable password=%s", port, user, dbName, pass)
	db, err = gorm.Open("postgres", connString)

	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(model.User{})
	db.AutoMigrate(model.FollowingRequest{})
	db.AutoMigrate(model.Follower{})
	return db
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
}

func initFollowerRepository(database *gorm.DB) *repository.FollowerRepository {
	return &repository.FollowerRepository{Database: database}
}

func initFollowingRequestRepository(database *gorm.DB) *repository.FollowingRequestRepository {
	return &repository.FollowingRequestRepository{Database: database}
}

func initFollowingService(followerRepository *repository.FollowerRepository, followingRequestRepository *repository.FollowingRequestRepository) *service.FollowingService {
	return &service.FollowingService{FollowerRepository: followerRepository, FollowingRequestRepository: followingRequestRepository}
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

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database := initDB()

	port := fmt.Sprintf(":%s", os.Getenv("PORT"))

	userRepo := initUserRepo(database)
	auth0Client := initAuth0Client()
	userService := initUserService(userRepo, auth0Client)
	userHandler := initUserHandler(userService)

	followingReqRepo := initFollowingRequestRepository(database)
	followerRepo := initFollowerRepository(database)
	followingService := initFollowingService(followerRepo, followingReqRepo)
	followingHandler := initFollowingHandler(followingService)

	router := gin.Default()

	handleFollowingFunc(followingHandler, router)
	handleUserFunc(userHandler, router)

	http.ListenAndServe(port, cors.AllowAll().Handler(router))
}
