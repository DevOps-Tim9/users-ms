package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-ms/src/dto"
	"user-ms/src/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
	Logger  *logrus.Entry
}

func (handler *UserHandler) Register(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "POST /register")
	defer span.Finish()

	var userToRegister dto.RegistrationRequestDTO
	if err := ctx.ShouldBindJSON(&userToRegister); err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	userID, err := handler.Service.Register(&userToRegister)

	if err != nil {
		AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("User with username %s failed to register", userToRegister.Username))
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("New user registered with id %d", userID))

	ctx.JSON(http.StatusCreated, userID)
}

func (handler *UserHandler) GetByEmail(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /users")
	defer span.Finish()

	email := ctx.Query("email")
	user, err := handler.Service.GetByEmail(email)
	fmt.Println(err)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (handler *UserHandler) GetByUsername(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /users/username")
	defer span.Finish()

	username := ctx.Query("username")

	users := handler.Service.GetByUsername(username)

	ctx.JSON(http.StatusOK, users)
}

func (handler *UserHandler) GetByID(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /users/:id")
	defer span.Finish()

	idStr := ctx.Param("id")
	id, _ := getId(idStr)
	user, err := handler.Service.GetByID(id)
	fmt.Println(err)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (handler *UserHandler) Update(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "PUT /users")
	defer span.Finish()

	var userToUpdate dto.UserUpdateDTO
	if err := ctx.ShouldBindJSON(&userToUpdate); err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	userDTO, err := handler.Service.Update(&userToUpdate)

	if err != nil {
		AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("User with id %d failed to update his profile info", userDTO.ID))
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("User with id %d updated his profile info", userDTO.ID))
	ctx.JSON(http.StatusOK, userDTO)
}

func extractClaims(tokenStr string) (jwt.MapClaims, bool) {
	token, _ := jwt.Parse(strings.Split(tokenStr, " ")[1], nil)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, true
	} else {
		fmt.Println("Invalid JWT Token")
		return nil, false
	}
}

func getId(idParam string) (int, error) {
	id, err := strconv.ParseInt(idParam, 10, 32)
	if err != nil {
		return 0, errors.New("ID should be a number")
	}
	return int(id), nil
}

func (handler *UserHandler) BlockUser(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "PUT /users/block-user")
	defer span.Finish()

	blockingID := ctx.Query("id")

	claims, _ := extractClaims(ctx.Request.Header.Get("Authorization"))

	id, _ := getId(blockingID)

	err := handler.Service.BlockUser(id, fmt.Sprint(claims["sub"]))

	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("User with auth0 id %s blocked user with id %d", fmt.Sprint(claims["sub"]), id))

	ctx.JSON(http.StatusOK, "User successfully blocked")
}

func (handler *UserHandler) UnblockUser(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "PUT /users/unblock-user")
	defer span.Finish()

	blockingID := ctx.Query("id")

	claims, _ := extractClaims(ctx.Request.Header.Get("Authorization"))

	id, _ := getId(blockingID)

	err := handler.Service.UnblockUser(id, fmt.Sprint(claims["sub"]))
	fmt.Println(err)
	if err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("User with auth0 id %s unblocked user with id %d", fmt.Sprint(claims["sub"]), id))

	ctx.JSON(http.StatusOK, "User successfully unblocked")
}

func (handler *UserHandler) GetBlockedUsers(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /users/blocked-users")
	defer span.Finish()

	claims, _ := extractClaims(ctx.Request.Header.Get("Authorization"))

	blockedUsers := handler.Service.GetBlockedUsers(fmt.Sprint(claims["sub"]))

	ctx.JSON(http.StatusOK, blockedUsers)
}

func (handler *UserHandler) SetNotifications(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "POST /users/set-notifications")
	defer span.Finish()

	var notificationSettings dto.NotificationsUpdateDTO
	if err := ctx.ShouldBindJSON(&notificationSettings); err != nil {
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	claims, _ := extractClaims(ctx.Request.Header.Get("Authorization"))

	err := handler.Service.SetNotifications(&notificationSettings, fmt.Sprint(claims["sub"]))
	if err != nil {
		AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("User with auth0 id %s failed to update notifications settings", fmt.Sprint(claims["sub"])))
		handler.Logger.Debug(err.Error())
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	AddSystemEvent(time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf("User with auth0 id %s updated notifications settings", fmt.Sprint(claims["sub"])))

	ctx.JSON(http.StatusOK, true)
}

func (handler *UserHandler) GetNotifications(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /users/get-notifications")
	defer span.Finish()

	claims, _ := extractClaims(ctx.Request.Header.Get("Authorization"))

	notificationSettings := handler.Service.GetNotifications(fmt.Sprint(claims["sub"]))

	ctx.JSON(http.StatusOK, notificationSettings)
}

func (handler *UserHandler) GetAll(ctx *gin.Context) {
	users, err := handler.Service.GetAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (handler *UserHandler) GetByParam(ctx *gin.Context) {
	searchParam := ctx.Query("param")
	users, err := handler.Service.GetByParam(searchParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}
