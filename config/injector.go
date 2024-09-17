//go:build wireinject
// +build wireinject

// go:build wireinject
package config

import (
	controller "awesomeProject/app/controllers"
	"awesomeProject/app/repository"
	"awesomeProject/app/service"
	"github.com/google/wire"
)

var db = wire.NewSet(ConnectToDB)

var userServiceSet = wire.NewSet(service.UserServiceInit,
	wire.Bind(new(service.UserService), new(*service.UserServiceImpl)),
)

var userRepoSet = wire.NewSet(repository.UserRepositoryInit,
	wire.Bind(new(repository.UserRepository), new(*repository.UserRepositoryImpl)),
)

var userCtrlSet = wire.NewSet(controller.UserControllerInit,
	wire.Bind(new(controller.UserController), new(*controller.UserControllerImpl)),
)

var authServiceSet = wire.NewSet(service.AuthServiceInit,
	wire.Bind(new(service.AuthService), new(*service.AuthServiceImpl)),
)

var authRepoSet = wire.NewSet(repository.AuthRepositoryInit,
	wire.Bind(new(repository.AuthRepository), new(*repository.AuthRepositoryImpl)),
)

var authCtrlSet = wire.NewSet(controller.AuthControllerInit,
	wire.Bind(new(controller.AuthController), new(*controller.AuthControllerImpl)),
)

var orderServiceSet = wire.NewSet(service.OrderServiceInit,
	wire.Bind(new(service.OrderService), new(*service.OrderServiceImpl)),
)

var orderRepoSet = wire.NewSet(repository.OrderRepositoryInit,
	wire.Bind(new(repository.OrderRepository), new(*repository.OrderRepositoryImpl)),
)

var orderCtrlSet = wire.NewSet(controller.OrderControllerInit,
	wire.Bind(new(controller.OrderController), new(*controller.OrderControllerImpl)),
)

func InitDependencies() *Initialization {
	wire.Build(NewInitialization, db, userCtrlSet, userServiceSet, userRepoSet, authServiceSet, authRepoSet, authCtrlSet, orderServiceSet, orderRepoSet, orderCtrlSet)
	return nil
}
