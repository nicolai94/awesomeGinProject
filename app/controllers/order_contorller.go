package controllers

import (
	"awesomeProject/app/service"
	"github.com/gin-gonic/gin"
)

type OrderController interface {
	CreateOrder(c *gin.Context)
}

type OrderControllerImpl struct {
	svc service.OrderService
}

func (u OrderControllerImpl) CreateOrder(c *gin.Context) {
	u.svc.CreateOrder(c)
}

func OrderControllerInit(orderService service.OrderService) *OrderControllerImpl {
	return &OrderControllerImpl{
		svc: orderService,
	}
}
