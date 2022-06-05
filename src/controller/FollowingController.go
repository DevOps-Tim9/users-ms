package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"user-ms/src/handler"
	"user-ms/src/repository"
	"user-ms/src/service"
)

type FollowingController struct{}

func (controller *FollowingController) InitFollowerRepository(database *gorm.DB) *repository.FollowerRepository {
	return &repository.FollowerRepository{Database: database}
}

func (controller *FollowingController) InitFollowingRequestRepository(database *gorm.DB) *repository.FollowingRequestRepository {
	return &repository.FollowingRequestRepository{Database: database}
}

func (controller *FollowingController) InitFollowingService(followerRepository *repository.FollowerRepository, followingRequestRepository *repository.FollowingRequestRepository) *service.FollowingService {
	return &service.FollowingService{FollowerRepository: followerRepository, FollowingRequestRepository: followingRequestRepository}
}

func (controller *FollowingController) InitFollowingHandler(service *service.FollowingService) *handler.FollowingHandler {
	return &handler.FollowingHandler{Service: service}
}

func (controller *FollowingController) HandleFollowingFunc(handler *handler.FollowingHandler, router *gin.Engine) {
	router.POST("/requests", handler.CreateRequest)
	router.POST("/follower", handler.CreatFollower)
	router.PUT("/requests/:id", handler.UpdateRequest)
	router.GET("/requests", handler.GetRequest)
	router.GET("/requests/:id", handler.GetRequestsByFollowingID)
	router.GET("user/:id/followers", handler.GetFollowers)
	router.GET("user/:id/following", handler.GetFollowing)
	router.DELETE("user/:id/removeFollower/:followingId", handler.RemoveFollowing)
}
