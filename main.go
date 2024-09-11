package main

import (
	"awesomeProject/app/router"
	"awesomeProject/app/utils"
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
	utils.Init()
	init := config.InitDependencies()

	app := router.Init(init)

	app.Run(":" + port)
}
