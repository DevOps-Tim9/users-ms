package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"user-ms/src/dto"
	"user-ms/src/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/opentracing/opentracing-go"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func (handler *UserHandler) Register(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "POST /register")
	defer span.Finish()

	var userToRegister dto.RegistrationRequestDTO
	if err := ctx.ShouldBindJSON(&userToRegister); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	userID, err := handler.Service.Register(&userToRegister)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusCreated, userID)
}

func (handler *UserHandler) GetByEmail(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /users")
	defer span.Finish()

	email := ctx.Query("email")
	user, err := handler.Service.GetByEmail(email)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (handler *UserHandler) GetByID(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /users/:id")
	defer span.Finish()

	idStr := ctx.Param("id")
	id, _ := getId(idStr)
	user, err := handler.Service.GetByID(id)
	fmt.Println(err)
	if err != nil {
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	userDTO, err := handler.Service.Update(&userToUpdate)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

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
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}
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
		ctx.JSON(http.StatusNotFound, err.Error())
		return
	}
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
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	claims, _ := extractClaims(ctx.Request.Header.Get("Authorization"))

	err := handler.Service.SetNotifications(&notificationSettings, fmt.Sprint(claims["sub"]))
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	ctx.JSON(http.StatusOK, true)
}

func (handler *UserHandler) GetNotifications(ctx *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(ctx.Request.Context(), "GET /users/get-notifications")
	defer span.Finish()

	claims, _ := extractClaims(ctx.Request.Header.Get("Authorization"))

	notificationSettings := handler.Service.GetNotifications(fmt.Sprint(claims["sub"]))

	ctx.JSON(http.StatusOK, notificationSettings)
}
