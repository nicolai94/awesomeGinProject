package config

import (
	controller "awesomeProject/app/controllers"
	"awesomeProject/app/repository"
	"awesomeProject/app/service"
)

type Initialization struct {
	userRepo repository.UserRepository
	userAuth repository.AuthRepository
	userSvc  service.UserService
	authSvc  service.AuthService
	UserCtrl controller.UserController
	AuthCtrl controller.AuthController
}

func NewInitialization(userRepo repository.UserRepository,
	userAuth repository.AuthRepository,
	userService service.UserService,
	authService service.AuthService,
	userCtrl controller.UserController,
	authCtrl controller.AuthController,
) *Initialization {
	return &Initialization{
		userRepo: userRepo,
		userSvc:  userService,
		UserCtrl: userCtrl,
		userAuth: userAuth,
		authSvc:  authService,
		AuthCtrl: authCtrl,
	}
}
