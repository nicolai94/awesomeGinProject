package config

import (
	controller "awesomeProject/app/controllers"
	"awesomeProject/app/repository"
	"awesomeProject/app/service"
)

type Initialization struct {
	userRepo  repository.UserRepository
	userAuth  repository.AuthRepository
	orderRepo repository.OrderRepository
	userSvc   service.UserService
	authSvc   service.AuthService
	orderSvc  service.OrderService
	UserCtrl  controller.UserController
	AuthCtrl  controller.AuthController
	OrderCtrl controller.OrderController
}

func NewInitialization(
	userRepo repository.UserRepository,
	userAuth repository.AuthRepository,
	orderRepo repository.OrderRepository,
	userService service.UserService,
	authService service.AuthService,
	orderService service.OrderService,
	userCtrl controller.UserController,
	authCtrl controller.AuthController,
	orderCtrl controller.OrderController,
) *Initialization {
	return &Initialization{
		userRepo:  userRepo,
		userSvc:   userService,
		UserCtrl:  userCtrl,
		userAuth:  userAuth,
		authSvc:   authService,
		AuthCtrl:  authCtrl,
		orderRepo: orderRepo,
		orderSvc:  orderService,
		OrderCtrl: orderCtrl,
	}
}
