package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"user-ms/src/dto"
	"user-ms/src/service"

	"github.com/gin-gonic/gin"
)

type FollowingHandler struct {
	Service *service.FollowingService
}

func (handler *FollowingHandler) UpdateRequest(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	fmt.Println(id)
	var requestDTO dto.FollowingRequestDTO
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	requestId, err := handler.Service.UpdateRequest(id, &requestDTO)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requestId)
}

func (handler *FollowingHandler) CreateRequest(ctx *gin.Context) {
	var requestDTO dto.FollowingRequestDTO
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	requestId, err := handler.Service.CreateRequest(&requestDTO)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requestId)
}

func (handler *FollowingHandler) GetRequest(ctx *gin.Context) {
	requests, err := handler.Service.GetRequests()
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requests)
}

func (handler *FollowingHandler) GetRequestsByFollowingID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	requests, err := handler.Service.GetRequestsByFollowingID(id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requests)
}

func (handler *FollowingHandler) CreatFollower(ctx *gin.Context) {
	var requestDTO dto.FollowingRequestDTO
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	requestId, err := handler.Service.CreateFollower(&requestDTO)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requestId)
}

func (handler *FollowingHandler) GetFollowers(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	requests, err := handler.Service.GetFollowers(id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requests)
}

func (handler *FollowingHandler) GetFollowing(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	requests, err := handler.Service.GetFollowing(id)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requests)
}

func (handler *FollowingHandler) RemoveFollowing(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	followingId, err := strconv.Atoi(ctx.Param("followingId"))
	err = handler.Service.RemoveFollowing(id, followingId)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}
