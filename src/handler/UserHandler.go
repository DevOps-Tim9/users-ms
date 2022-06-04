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
