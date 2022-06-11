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

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Service *service.UserService
}

func (handler *UserHandler) Register(ctx *gin.Context) {
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
	email := ctx.Query("email")
	user, err := handler.Service.GetByEmail(email)
	fmt.Println(err)
	if err != nil {
		ctx.JSON(http.StatusNotFound, err)
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (handler *UserHandler) Update(ctx *gin.Context) {
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
	claims, _ := extractClaims(ctx.Request.Header.Get("Authorization"))

	blockedUsers := handler.Service.GetBlockedUsers(fmt.Sprint(claims["sub"]))

	ctx.JSON(http.StatusOK, blockedUsers)
}
