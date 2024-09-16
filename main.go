package main

import (
	"awesomeProject/app/router"
	"awesomeProject/app/utils"
	"awesomeProject/config"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Ошибка загрузки .env файла:", err)
	}
	config.InitLog()
	port := os.Getenv("PORT")

	config.ConnectToDB()
	utils.ConnectToRedis()

	init := config.InitDependencies()

	app := router.Init(init)
	app.Run(":" + port)
}
