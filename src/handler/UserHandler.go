package handler

import (
	"fmt"
	"net/http"
	"user-ms/src/dto"
	"user-ms/src/service"

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

