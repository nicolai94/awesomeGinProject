package main

import (
	"awesomeProject/app/router"
	"awesomeProject/config"
	"github.com/joho/godotenv"
	"os"
)

func init() {
	godotenv.Load()
	config.InitLog()
}

func main() {
	port := os.Getenv("PORT")

	config.ConnectToDB()
	config.ConnectRedis()
	//config.ConnectToDB().AutoMigrate()

	init := config.InitDependencies()

	app := router.Init(init)

	app.Run(":" + port)
}
