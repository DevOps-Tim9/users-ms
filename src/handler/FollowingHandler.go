package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"user-ms/src/dto"
	"user-ms/src/service"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type FollowingHandler struct {
	Service *service.FollowingService
	Logger  *logrus.Entry
}

func (handler *FollowingHandler) UpdateRequest(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "PUT /requests/:id")
	defer span.Finish()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	fmt.Println(id)
	var requestDTO dto.FollowingRequestDTO
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	requestId, err := handler.Service.UpdateRequest(id, &requestDTO)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requestId)
}

func (handler *FollowingHandler) CreateRequest(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "POST /requests")
	defer span.Finish()

	var requestDTO dto.FollowingRequestDTO
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	requestId, err := handler.Service.CreateRequest(&requestDTO)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requestId)
}

func (handler *FollowingHandler) GetRequest(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /requests")
	defer span.Finish()

	requests, err := handler.Service.GetRequests()
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requests)
}

func (handler *FollowingHandler) GetRequestsByFollowingID(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /requests/:id")
	defer span.Finish()

	id, err := strconv.Atoi(ctx.Param("id"))
	requests, err := handler.Service.GetRequestsByFollowingID(id)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requests)
}

func (handler *FollowingHandler) CreatFollower(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "POST /follower")
	defer span.Finish()

	var requestDTO dto.FollowingRequestDTO
	if err := ctx.ShouldBindJSON(&requestDTO); err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	requestId, err := handler.Service.CreateFollower(&requestDTO)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requestId)
}

func (handler *FollowingHandler) GetFollowers(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /user/:id/followers")
	defer span.Finish()

	id, err := strconv.Atoi(ctx.Param("id"))
	requests, err := handler.Service.GetFollowers(id)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requests)
}

func (handler *FollowingHandler) GetFollowing(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET user/:id/following")
	defer span.Finish()

	id, err := strconv.Atoi(ctx.Param("id"))
	requests, err := handler.Service.GetFollowing(id)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, requests)
}

func (handler *FollowingHandler) RemoveFollowing(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "DELETE /user/:id/removeFollower/:followingId+")
	defer span.Finish()

	id, err := strconv.Atoi(ctx.Param("id"))
	followingId, err := strconv.Atoi(ctx.Param("followingId"))
	err = handler.Service.RemoveFollowing(id, followingId)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, nil)
}
