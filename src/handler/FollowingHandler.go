package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"user-ms/src/dto"
	"user-ms/src/service"
	"user-ms/src/utils"

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

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("Following request with id %d updated", id))

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

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("Following request with id %d created", requestId))

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

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("User with id %d started following user with id %d", requestDTO.FollowerId, requestDTO.FollowingId))

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

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("User with id %d is not foolowing user %d anymore", id, followingId))

	ctx.JSON(http.StatusCreated, nil)
}

func AddSystemEvent(time string, message string) error {
	logger := utils.Logger()
	event := dto.EventRequestDTO{
		Timestamp: time,
		Message:   message,
	}

	b, _ := json.Marshal(&event)
	endpoint := os.Getenv("EVENTS_MS")
	logger.Info("Sending system event to events-ms")
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	req.Header.Set("content-type", "application/json")

	_, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Debug("Error happened during sending system event")
		return err
	}

	return nil
}
