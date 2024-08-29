package controllers

import (
	"awesomeProject/app/service"
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Login(c *gin.Context)
}

type AuthControllerImpl struct {
	svc service.AuthService
}

func (u AuthControllerImpl) Login(c *gin.Context) {
	u.svc.Login(c)
}

func AuthControllerInit(authService service.AuthService) *AuthControllerImpl {
	return &AuthControllerImpl{
		svc: authService,
	}
}
